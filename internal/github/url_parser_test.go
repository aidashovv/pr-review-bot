package github

import (
	"fmt"
	"testing"
)

func TestParsePullRequestURL_OK(t *testing.T) {
	ref, err := ParsePullRequestURL("https://github.com/cat/Hello-World/pull/1347")
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	if ref.Owner != "cat" || ref.Repo != "Hello-World" || ref.Number != 1347 {
		t.Fatalf("unexpected ref: %+v", ref)
	}
}

func TestParsePullRequestURL_Bad(t *testing.T) {
	cases := []string{
		"",
		"https://gitlab.com/a/b/pull/1",
		"https://github.com/a/b/issues/1",
		"https://github.com/a/b/pull/abc",
		"https://github.com/a/b/pull/0",
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			if _, err := ParsePullRequestURL(c); err == nil {
				t.Fatalf("expected error for %q", c)
			}
		})
	}
}
