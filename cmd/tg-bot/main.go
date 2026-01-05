package tgbot

import (
	"fmt"
	"os"

	"github.com/aidashovv/pr-review-bot/internal/buildinfo"
	"github.com/aidashovv/pr-review-bot/internal/config"
)

func main() {
	cfg := config.LoadFromEnv()

	if err := cfg.ValidateForBot(); err != nil {
		fmt.Fprintln(os.Stderr, "config error", err)
		os.Exit(1)
	}

	fmt.Printf(
		"tg-bot starting (version=%s commit=%s built=%s)\n",
		buildinfo.Version, buildinfo.Commit, buildinfo.Built,
	)

	// TODO: tg-bot start
	select {}
}
