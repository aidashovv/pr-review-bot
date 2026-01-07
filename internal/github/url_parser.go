package github

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type PRRef struct {
	Owner  string
	Repo   string
	Number int
}

// example: https://github.com/OWNER/REPO/pull/123
func ParsePullRequestURL(raw string) (PRRef, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Hostname() == "" || u.Scheme == "" {
		return PRRef{}, fmt.Errorf("unsupported url format")
	}

	host := strings.ToLower(u.Hostname())
	if host != "github.com" && host != "www.github.com" {
		return PRRef{}, fmt.Errorf("unsupported host")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "pull" {
		return PRRef{}, fmt.Errorf("unsupported url format")
	}

	owner, repo := parts[0], parts[1]
	if owner == "" || repo == "" {
		return PRRef{}, fmt.Errorf("unsupported url format")
	}

	n, convErr := strconv.Atoi(parts[3])
	if convErr != nil || n <= 0 {
		return PRRef{}, fmt.Errorf("invalid pull request number")
	}

	return PRRef{
		Owner:  owner,
		Repo:   repo,
		Number: n,
	}, nil
}
