package github

import (
	"context"
	"fmt"
)

type DiffResult struct {
	Diff      string `json:"diff"`
	Truncated bool   `json:"truncated"`
}

func (c *Client) GetPRDiff(ctx context.Context, ref PRRef, maxBytes int, maxFiles int) (DiffResult, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", apiBaseURL, ref.Owner, ref.Repo, ref.Number)

	req, err := c.newRequest(ctx, "GET", endpoint, acceptDiff)
	if err != nil {
		return DiffResult{}, err
	}

	status, headers, body, truncatedBytes, err := c.doRequest(req, maxBytes)
	if err != nil {
		return DiffResult{}, err
	}

	if status < 200 || status >= 300 {
		return DiffResult{}, classifyAPIError(status, headers, body)
	}

	diffText := string(body)

	diffLimited, truncatedFiles := LimitDiffFiles(diffText, maxFiles)

	return DiffResult{
		Diff:      diffLimited,
		Truncated: truncatedBytes || truncatedFiles,
	}, nil
}
