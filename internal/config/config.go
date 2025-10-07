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
	WhitelistedIPs []string
	AllowedOrigins []string
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

	// Parse whitelisted IPs (comma-separated)
	if whitelistedIPs := os.Getenv("WHITELISTED_IPS"); whitelistedIPs != "" {
		for _, ip := range splitAndTrim(whitelistedIPs, ",") {
			if ip != "" {
				cfg.WhitelistedIPs = append(cfg.WhitelistedIPs, ip)
			}
		}
	}

	// Parse allowed origins (comma-separated)
	if allowedOrigins := os.Getenv("ALLOWED_ORIGINS"); allowedOrigins != "" {
		for _, origin := range splitAndTrim(allowedOrigins, ",") {
			if origin != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
			}
		}
	} else {
		// Default to allow all if not specified
		cfg.AllowedOrigins = []string{"*"}
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

func splitAndTrim(s, sep string) []string {
	var result []string
	for i := 0; i < len(s); {
		// Skip leading separator
		for i < len(s) && string(s[i]) == sep {
			i++
		}
		if i >= len(s) {
			break
		}
		// Find next separator
		start := i
		for i < len(s) && string(s[i]) != sep {
			i++
		}
		// Trim spaces from the token
		token := s[start:i]
		trimmed := trimSpace(token)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for start < end && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
