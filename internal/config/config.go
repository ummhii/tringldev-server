package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LastFMAPIKey   string
	LastFMUsername string
	GithubToken    string
	GithubUsername string
	Port           string
	DiscordWebhook string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &Config{
		LastFMAPIKey:   os.Getenv("LASTFM_API_KEY"),
		LastFMUsername: os.Getenv("LASTFM_USERNAME"),
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		GithubUsername: os.Getenv("GITHUB_USERNAME"),
		Port:           os.Getenv("PORT"),
		DiscordWebhook: os.Getenv("DISCORD_WEBHOOK"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Warn if required configs are missing
	if cfg.LastFMAPIKey == "" {
		log.Println("Warning: LASTFM_API_KEY not set")
	}
	if cfg.LastFMUsername == "" {
		log.Println("Warning: LASTFM_USERNAME not set")
	}
	if cfg.GithubToken == "" {
		log.Println("Warning: GITHUB_TOKEN not set")
	}
	if cfg.GithubUsername == "" {
		log.Println("Warning: GITHUB_USERNAME not set")
	}
	if cfg.DiscordWebhook == "" {
		log.Println("Warning: DISCORD_WEBHOOK not set")
	}

	return cfg
}
