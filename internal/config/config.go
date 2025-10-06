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

	// SMTP Configuration for contact form
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	SMTPTo       string

	// Alternative: Discord webhook for contact form
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

		// SMTP settings
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		SMTPTo:       os.Getenv("SMTP_TO"),

		// Discord webhook (alternative to SMTP)
		DiscordWebhook: os.Getenv("DISCORD_WEBHOOK"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.SMTPPort == "" {
		cfg.SMTPPort = "587"
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

	return cfg
}
