package contact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tringldev-server/internal/config"
)

type Service struct {
	config *config.Config
}

type ContactRequest struct {
	Email   string `json:"email" form:"email"`
	Name    string `json:"name" form:"name"`
	Message string `json:"message" form:"message"`
}

type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DiscordWebhook struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []DiscordEmbedField `json:"fields"`
	Timestamp   string              `json:"timestamp"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// Send sends the contact form via Discord webhook
func (s *Service) Send(req *ContactRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.Message) == "" {
		return fmt.Errorf("message is required")
	}

	if s.config.DiscordWebhook == "" {
		return fmt.Errorf("discord webhook not configured (set DISCORD_WEBHOOK)")
	}

	return s.SendToDiscord(req)
}

// sends the contact form to a Discord webhook
func (s *Service) SendToDiscord(req *ContactRequest) error {
	if s.config.DiscordWebhook == "" {
		return fmt.Errorf("discord webhook not configured")
	}

	// Build Discord embed
	embed := DiscordEmbed{
		Title:       "New Contact Form Submission",
		Description: req.Message,
		Color:       5814783, // Blue
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields: []DiscordEmbedField{
			{
				Name:   "Name",
				Value:  req.Name,
				Inline: true,
			},
		},
	}

	// Add email field if provided
	if req.Email != "" {
		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   "Email",
			Value:  req.Email,
			Inline: true,
		})
	}

	webhook := DiscordWebhook{
		Embeds: []DiscordEmbed{embed},
	}

	payload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal discord webhook: %w", err)
	}

	resp, err := http.Post(s.config.DiscordWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}
