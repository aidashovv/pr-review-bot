package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiBaseURL = "https://api.github.com"
	userAgent  = "pr-review-bot-mcp-github"
	acceptJSON = "application/vnd.github+json"
	acceptDiff = "application/vnd.github.diff"
)

type Client struct {
	http  *http.Client
	token string
}

func NewClient(token string) *Client {
	return &Client{
		http:  &http.Client{Timeout: 10 * time.Second},
		token: strings.TrimSpace(token),
	}
}

func (c *Client) newRequest(ctx context.Context, method, url, accept string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", userAgent)

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, maxBytes int) (
	status int,
	headers http.Header,
	body []byte,
	truncated bool, err error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, nil, false, err
	}
	defer resp.Body.Close()

	limit := int64(maxBytes)
	if maxBytes <= 0 {
		limit = 10 * 1024 * 1024 // 10 MB
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return resp.StatusCode, resp.Header, nil, false, err
	}

	truncated = int64(len(data)) > limit
	if truncated {
		data = data[:limit]
	}

	return resp.StatusCode, resp.Header, data, truncated, nil
}

type apiError struct {
	Message string `json:"message"`
}

func classifyAPIError(status int, headers http.Header, body []byte) error {
	msg := ""

	var apiErr apiError
	if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
		msg = strings.ToLower(apiErr.Message)
	}

	remaining := headers.Get("X-RateLimit-Remaining")
	if (status == 403 || status == 429) && (strings.Contains(msg, "rate limit") || remaining == "0") {
		return errors.New("rate limited")
	}

	switch status {
	case 401:
		return errors.New("unauthorized (invalid token)")
	case 403:
		return errors.New("forbidden (repo access or rate limit)")
	case 404:
		return errors.New("not found")
	default:
		if apiErr.Message != "" {
			return fmt.Errorf("github api error (%d): %s", status, apiErr.Message)
		}
		return fmt.Errorf("github api error (%d)", status)
	}
}
