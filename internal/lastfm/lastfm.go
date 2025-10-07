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

type topArtistsResponse struct {
	TopArtists struct {
		Artist []struct {
			Name      string `json:"name"`
			PlayCount string `json:"playcount"`
			URL       string `json:"url"`
		} `json:"artist"`
	} `json:"topartists"`
}

type TopArtist struct {
	Name      string `json:"name"`
	PlayCount string `json:"playcount"`
}

type TopArtistsInfo struct {
	Artists []TopArtist `json:"artists"`
}

type topTracksResponse struct {
	TopTracks struct {
		Track []struct {
			Name      string `json:"name"`
			PlayCount string `json:"playcount"`
			URL       string `json:"url"`
			Artist    struct {
				Name string `json:"name"`
			} `json:"artist"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
		} `json:"track"`
	} `json:"toptracks"`
}

type TopTrack struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	PlayCount string `json:"playcount"`
	AlbumArt  string `json:"albumArt,omitempty"`
	URL       string `json:"url"`
}

type TopTracksInfo struct {
	Tracks []TopTrack `json:"tracks"`
}

type topAlbumsResponse struct {
	TopAlbums struct {
		Album []struct {
			Name      string `json:"name"`
			PlayCount string `json:"playcount"`
			URL       string `json:"url"`
			Artist    struct {
				Name string `json:"name"`
			} `json:"artist"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
		} `json:"album"`
	} `json:"topalbums"`
}

type TopAlbum struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	PlayCount string `json:"playcount"`
	AlbumArt  string `json:"albumArt,omitempty"`
	URL       string `json:"url"`
}

type TopAlbumsInfo struct {
	Albums []TopAlbum `json:"albums"`
}

type RecentTrack struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	AlbumArt  string `json:"albumArt,omitempty"`
	URL       string `json:"url"`
	PlayedAt  string `json:"playedAt"`
	IsPlaying bool   `json:"isPlaying"`
}

type RecentTracksInfo struct {
	Tracks []RecentTrack `json:"tracks"`
}

type userInfoResponse struct {
	User struct {
		Name       string `json:"name"`
		PlayCount  string `json:"playcount"`
		Registered struct {
			UnixTime string `json:"unixtime"`
		} `json:"registered"`
	} `json:"user"`
}

type ListeningStats struct {
	TotalScrobbles string `json:"totalScrobbles"`
	AccountAge     string `json:"accountAge"`
	Username       string `json:"username"`
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

// limit: number of artists to return (default: 10, max: 50)
func (s *Service) GetTopWeeklyArtists(limit int) (*TopArtistsInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.gettopartists&user=%s&api_key=%s&format=json&period=7day&limit=%d",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
		limit,
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

	var lastfmResp topArtistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	info := &TopArtistsInfo{
		Artists: make([]TopArtist, 0),
	}

	for _, artist := range lastfmResp.TopArtists.Artist {
		info.Artists = append(info.Artists, TopArtist{
			Name:      artist.Name,
			PlayCount: artist.PlayCount,
		})
	}

	return info, nil
}

// Accepts: 7day, 1month, 3month, 6month, 12month, overall (default: 7day)
func validatePeriod(period string) string {
	validPeriods := map[string]string{
		"weekly":  "7day",
		"7day":    "7day",
		"monthly": "1month",
		"1month":  "1month",
		"3month":  "3month",
		"6month":  "6month",
		"yearly":  "12month",
		"12month": "12month",
		"alltime": "overall",
		"overall": "overall",
	}

	if val, ok := validPeriods[period]; ok {
		return val
	}
	return "7day"
}

func (s *Service) GetTopArtists(limit int, period string) (*TopArtistsInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	period = validatePeriod(period)

	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.gettopartists&user=%s&api_key=%s&format=json&period=%s&limit=%d",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
		period,
		limit,
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

	var lastfmResp topArtistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	info := &TopArtistsInfo{
		Artists: make([]TopArtist, 0),
	}

	for _, artist := range lastfmResp.TopArtists.Artist {
		info.Artists = append(info.Artists, TopArtist{
			Name:      artist.Name,
			PlayCount: artist.PlayCount,
		})
	}

	return info, nil
}

func (s *Service) GetTopTracks(limit int, period string) (*TopTracksInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	period = validatePeriod(period)

	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.gettoptracks&user=%s&api_key=%s&format=json&period=%s&limit=%d",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
		period,
		limit,
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

	var lastfmResp topTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	info := &TopTracksInfo{
		Tracks: make([]TopTrack, 0),
	}

	for _, track := range lastfmResp.TopTracks.Track {
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

		info.Tracks = append(info.Tracks, TopTrack{
			Name:      track.Name,
			Artist:    track.Artist.Name,
			PlayCount: track.PlayCount,
			AlbumArt:  albumArt,
			URL:       track.URL,
		})
	}

	return info, nil
}

func (s *Service) GetTopAlbums(limit int, period string) (*TopAlbumsInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	period = validatePeriod(period)

	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.gettopalbums&user=%s&api_key=%s&format=json&period=%s&limit=%d",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
		period,
		limit,
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

	var lastfmResp topAlbumsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	info := &TopAlbumsInfo{
		Albums: make([]TopAlbum, 0),
	}

	for _, album := range lastfmResp.TopAlbums.Album {
		albumArt := ""
		for _, img := range album.Image {
			if img.Size == "extralarge" || img.Size == "large" {
				albumArt = img.Text
				break
			}
		}
		if albumArt == "" && len(album.Image) > 0 {
			albumArt = album.Image[len(album.Image)-1].Text
		}

		info.Albums = append(info.Albums, TopAlbum{
			Name:      album.Name,
			Artist:    album.Artist.Name,
			PlayCount: album.PlayCount,
			AlbumArt:  albumArt,
			URL:       album.URL,
		})
	}

	return info, nil
}

func (s *Service) GetRecentTracks(limit int) (*RecentTracksInfo, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=%d",
		url.QueryEscape(s.config.LastFMUsername),
		url.QueryEscape(s.config.LastFMAPIKey),
		limit,
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

	info := &RecentTracksInfo{
		Tracks: make([]RecentTrack, 0),
	}

	for _, track := range lastfmResp.RecentTracks.Track {
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

		isPlaying := track.Attr.NowPlaying == "true"
		playedAt := ""
		if !isPlaying && track.Date.UTS != "" {
			timestamp, err := strconv.ParseInt(track.Date.UTS, 10, 64)
			if err == nil {
				playedAt = time.Unix(timestamp, 0).Format(time.RFC3339)
			}
		}

		info.Tracks = append(info.Tracks, RecentTrack{
			Name:      track.Name,
			Artist:    track.Artist.Text,
			Album:     track.Album.Text,
			AlbumArt:  albumArt,
			URL:       track.URL,
			PlayedAt:  playedAt,
			IsPlaying: isPlaying,
		})
	}

	return info, nil
}

func (s *Service) GetListeningStats() (*ListeningStats, error) {
	apiURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.getinfo&user=%s&api_key=%s&format=json",
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

	var lastfmResp userInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastfmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	accountAge := ""
	if lastfmResp.User.Registered.UnixTime != "" {
		timestamp, err := strconv.ParseInt(lastfmResp.User.Registered.UnixTime, 10, 64)
		if err == nil {
			regDate := time.Unix(timestamp, 0)
			accountAge = time.Since(regDate).Round(24 * time.Hour).String()
		}
	}

	return &ListeningStats{
		TotalScrobbles: lastfmResp.User.PlayCount,
		AccountAge:     accountAge,
		Username:       lastfmResp.User.Name,
	}, nil
}
