package config

import (
	"fmt"
	"os"
)

type Config struct {
	TgBotToken  string
	GitHubToken string
	LogLevel    string
}

func LoadFromEnv() Config {
	return Config{
		TgBotToken:  os.Getenv("TG_BOT_TOKEN"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}
}

func (c Config) ValidateForBot() error {
	if c.TgBotToken == "" {
		return fmt.Errorf("TG_BOT_TOKEN is required")
	}

	return nil
}

func (c Config) ValidateForMCPGitHub() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}

	return nil
}
