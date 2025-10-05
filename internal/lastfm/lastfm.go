package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"tringldev-server/internal/config"
)

type Service struct {
	config *config.Config
}

type recentTracksResponse struct {
	RecentTracks struct {
		Track []struct {
			Name   string `json:"name"`
			Artist struct {
				Text string `json:"#text"`
			} `json:"artist"`
			Album struct {
				Text string `json:"#text"`
			} `json:"album"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			URL  string `json:"url"`
			Attr struct {
				NowPlaying string `json:"nowplaying"`
			} `json:"@attr"`
			Date struct {
				UTS string `json:"uts"`
			} `json:"date"`
		} `json:"track"`
	} `json:"recenttracks"`
}

type NowPlayingInfo struct {
	IsPlaying   bool   `json:"isPlaying"`
	SongName    string `json:"songName,omitempty"`
	ArtistName  string `json:"artistName,omitempty"`
	AlbumName   string `json:"albumName,omitempty"`
	AlbumArt    string `json:"albumArt,omitempty"`
	SongURL     string `json:"songUrl,omitempty"`
	LastUpdated string `json:"lastUpdated"`
	PlayedAt    string `json:"playedAt,omitempty"`
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) GetCurrentlyPlaying() (*NowPlayingInfo, error) {
	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=1",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Last.fm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("last.fm api error %d: %s", resp.StatusCode, string(body))
	}

	var lastfmResp recentTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	info := &NowPlayingInfo{
		IsPlaying:   false,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	if len(lastfmResp.RecentTracks.Track) == 0 {
		return info, nil
	}

	track := lastfmResp.RecentTracks.Track[0]

	isNowPlaying := track.Attr.NowPlaying == "true"
	info.IsPlaying = isNowPlaying

	albumArt := ""
	for _, img := range track.Image {
		if img.Size == "extralarge" || img.Size == "large" {
			albumArt = img.Text
			break
		}
	}

	if albumArt == "" && len(track.Image) > 0 {
		albumArt = track.Image[len(track.Image)-1].Text
	}

	info.SongName = track.Name
	info.ArtistName = track.Artist.Text
	info.AlbumName = track.Album.Text
	info.AlbumArt = albumArt
	info.SongURL = track.URL

	if !isNowPlaying && track.Date.UTS != "" {
		timestamp, err := strconv.ParseInt(track.Date.UTS, 10, 64)
		if err == nil {
			info.PlayedAt = time.Unix(timestamp, 0).Format(time.RFC3339)
		}
	}

	return info, nil
}
