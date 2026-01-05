package mcpgithub

import (
	"fmt"
	"os"

	"github.com/aidashovv/pr-review-bot/internal/buildinfo"
	"github.com/aidashovv/pr-review-bot/internal/config"
)

func main() {
	cfg := config.LoadFromEnv()

	if err := cfg.ValidateForMCPGitHub(); err != nil {
		fmt.Fprintln(os.Stderr, "config error", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr,
		"mcp-github starting (version=%s commit=%s built=%s)\n",
		buildinfo.Version, buildinfo.Commit, buildinfo.Built,
	)

	// TODO: mcp stdio + tools
	select {}
}
