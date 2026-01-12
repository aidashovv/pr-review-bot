package config

import (
	"fmt"
	"os"
)

type Config struct {
	GitHubToken string
}

func LoadFromEnv() Config {
	return Config{GitHubToken: os.Getenv("GITHUB_TOKEN")}
}

func (c Config) ValidateForMCPGitHub() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}

	return nil
}
