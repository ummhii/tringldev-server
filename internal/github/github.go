package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"tringldev-server/internal/config"
)

type Service struct {
	config *config.Config
}

type repository struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	HTMLURL     string   `json:"html_url"`
	Language    string   `json:"language"`
	Stars       int      `json:"stargazers_count"`
	Forks       int      `json:"forks_count"`
	Topics      []string `json:"topics"`
	UpdatedAt   string   `json:"updated_at"`
	Homepage    string   `json:"homepage"`
}

type PinnedRepo struct {
	Name        string   `json:"name"`
	FullName    string   `json:"fullName"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Language    string   `json:"language"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
	Topics      []string `json:"topics"`
	UpdatedAt   string   `json:"updatedAt"`
	Homepage    string   `json:"homepage,omitempty"`
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) GetPinnedRepository() (*PinnedRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=10", s.config.GithubUsername)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if s.config.GithubToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.GithubToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error %d: %s", resp.StatusCode, string(body))
	}

	var repos []repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	// Return the most recently updated repo
	repo := repos[0]

	return toPinnedRepo(&repo), nil
}

func (s *Service) GetAllPublicRepositories() ([]*PinnedRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=public&sort=updated&per_page=6", s.config.GithubUsername)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if s.config.GithubToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.GithubToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error %d: %s", resp.StatusCode, string(body))
	}

	var repos []repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	// Convert all repos to PinnedRepo format
	pinnedRepos := make([]*PinnedRepo, 0, len(repos))
	for i := range repos {
		pinnedRepos = append(pinnedRepos, toPinnedRepo(&repos[i]))
	}

	return pinnedRepos, nil
}

func (s *Service) GetSpecificRepository(repoName string) (*PinnedRepo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", s.config.GithubUsername, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if s.config.GithubToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.GithubToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error %d: %s", resp.StatusCode, string(body))
	}

	var repo repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return toPinnedRepo(&repo), nil
}

func toPinnedRepo(repo *repository) *PinnedRepo {
	return &PinnedRepo{
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		URL:         repo.HTMLURL,
		Language:    repo.Language,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		Topics:      repo.Topics,
		UpdatedAt:   repo.UpdatedAt,
		Homepage:    repo.Homepage,
	}
}
