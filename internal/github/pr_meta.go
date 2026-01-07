package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PRMeta struct {
	Repo      string    `json:"repo"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	BodyCut   bool      `json:"body_cut"`
	Author    string    `json:"author"`
	Base      string    `json:"base"`
	Head      string    `json:"head"`
	HeadSHA   string    `json:"head_sha"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type prResponse struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`

	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`

	Head struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Client) GetPRMeta(ctx context.Context, ref PRRef, maxBodyBytes int) (PRMeta, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", apiBaseURL, ref.Owner, ref.Repo, ref.Number)

	req, err := c.newRequest(ctx, "GET", endpoint, acceptJSON)
	if err != nil {
		return PRMeta{}, err
	}

	status, headers, body, _, err := c.doRequest(req, 2*1024*1024)
	if err != nil {
		return PRMeta{}, err
	}

	if status < 200 || status >= 300 {
		return PRMeta{}, classifyAPIError(status, headers, body)
	}

	var pr prResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return PRMeta{}, fmt.Errorf("failed to parse github response")
	}

	bodyTrimmed, cut := TruncateUTF8ByBytes(pr.Body, maxBodyBytes)

	return PRMeta{
		Repo:      ref.Owner + "/" + ref.Repo,
		Number:    pr.Number,
		Title:     pr.Title,
		Body:      bodyTrimmed,
		BodyCut:   cut,
		Author:    pr.User.Login,
		Base:      pr.Base.Ref,
		Head:      pr.Head.Ref,
		HeadSHA:   pr.Head.SHA,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt,
	}, nil
}
