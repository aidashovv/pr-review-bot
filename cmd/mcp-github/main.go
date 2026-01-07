package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aidashovv/pr-review-bot/internal/buildinfo"
	"github.com/aidashovv/pr-review-bot/internal/config"
	"github.com/aidashovv/pr-review-bot/internal/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultMaxBodyBytes = 16 * 1024
	defaultMaxDiffBytes = 200_000
	defaultMaxFiles     = 20
)

func main() {
	logger := log.New(os.Stderr, "mcp-github: ", log.LstdFlags|log.Lshortfile)

	cfg := config.LoadFromEnv()
	if err := cfg.ValidateForMCPGitHub(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	logger.Printf("starting %s version=%s commit=%s built=%s",
		"mcp-github", buildinfo.Version, buildinfo.Commit, buildinfo.Built)

	gh := github.NewClient(cfg.GitHubToken)

	s := server.NewMCPServer(
		"mcp-github",
		buildinfo.Version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	metaTool := mcp.NewTool(
		"github-get-pr-meta",
		mcp.WithDescription("fetch github pr metadata by pr-url"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("example: https://github.com/OWNER/REPO/pull/123"),
		),
		mcp.WithNumber("max_body_bytes",
			mcp.DefaultNumber(float64(defaultMaxBodyBytes)),
			mcp.Description("max size of body to return"),
			mcp.Min(1),
		),
	)

	s.AddTool(metaTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		maxBody := req.GetInt("max_body_bytes", defaultMaxBodyBytes)

		ref, err := github.ParsePullRequestURL(url)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		meta, err := gh.GetPRMeta(ctx, ref, maxBody)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(meta), nil
	})

	diffTool := mcp.NewTool(
		"github-get-pr-diff",
		mcp.WithDescription("fetch github pr unified diff by pr-url"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("example: https://github.com/OWNER/REPO/pull/123"),
		),
		mcp.WithNumber("max_bytes",
			mcp.DefaultNumber(float64(defaultMaxDiffBytes)),
			mcp.Description("max bytes of diff to return"),
			mcp.Min(1),
		),
		mcp.WithNumber("max_files",
			mcp.DefaultNumber(float64(defaultMaxFiles)),
			mcp.Description("max number of files (diff --git sections) to include"),
			mcp.Min(1),
		),
	)

	s.AddTool(diffTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		maxBytes := req.GetInt("max_bytes", defaultMaxDiffBytes)
		maxFiles := req.GetInt("max_files", defaultMaxFiles)

		ref, err := github.ParsePullRequestURL(url)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		diffRes, err := gh.GetPRDiff(ctx, ref, maxBytes, maxFiles)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultStructuredOnly(diffRes), nil
	})

	if err := server.ServeStdio(s); err != nil {
		logger.Printf("server error: %v", err)
		os.Exit(1)
	}
}
