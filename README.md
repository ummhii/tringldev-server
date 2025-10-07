# tringl.dev backend

A personal backend API for displaying last.fm stats and pinned GitHub projects and other stuff.

## Features

- **Last.fm Stats**: Displays Last.fm stats.
- **Github Stats**: Shows one of your GitHub repositories
- **Contact Form**: Receive messages via Discord webhook

## API Endpoints

### `GET /`
Health check endpoint
```json
{
  "status": "ok",
  "message": "TringlDev API Server",
  "version": "1.0.0"
}
```

### `GET /api/now-playing`
Returns the currently playing song from Last.fm (or most recently played)

**Response (when playing):**
```json
{
  "isPlaying": true,
  "songName": "Song Name",
  "artistName": "Artist Name",
  "albumName": "Album Name",
  "albumArt": "https://...",
  "songUrl": "https://www.last.fm/music/...",
  "lastUpdated": "2025-10-06T12:00:00Z"
}
```

**Response (when not playing):**
```json
{
  "isPlaying": false,
  "lastUpdated": "2025-10-06T12:00:00Z"
}
```

### `GET /api/top-artists`
Returns your top artists from Last.fm

**Query Parameters:**
- `limit` (optional): Number of artists to return (default: 10, max: 50)
- `period` (optional): Time period - `weekly`/`7day`, `monthly`/`1month`, `3month`, `6month`, `yearly`/`12month`, `alltime`/`overall` (default: `7day`)

**Example:** `/api/top-artists?limit=20&period=monthly`

**Response:**
```json
{
  "artists": [
    {
      "name": "Artist Name",
      "playcount": "123"
    },
    {
      "name": "Another Artist",
      "playcount": "45"
    }
  ]
}
```

### `GET /api/top-tracks`
Returns your top tracks from Last.fm

**Query Parameters:**
- `limit` (optional): Number of tracks to return (default: 10, max: 50)
- `period` (optional): Time period - `weekly`/`7day`, `monthly`/`1month`, `3month`, `6month`, `yearly`/`12month`, `alltime`/`overall` (default: `7day`)

**Example:** `/api/top-tracks?limit=15&period=alltime`

**Response:**
```json
{
  "tracks": [
    {
      "name": "Track Name",
      "artist": "Artist Name",
      "playcount": "89",
      "albumArt": "https://...",
      "url": "https://www.last.fm/music/..."
    }
  ]
}
```

### `GET /api/top-albums`
Returns your top albums from Last.fm

**Query Parameters:**
- `limit` (optional): Number of albums to return (default: 10, max: 50)
- `period` (optional): Time period - `weekly`/`7day`, `monthly`/`1month`, `3month`, `6month`, `yearly`/`12month`, `alltime`/`overall` (default: `7day`)

**Example:** `/api/top-albums?limit=15&period=yearly`

**Response:**
```json
{
  "albums": [
    {
      "name": "Album Name",
      "artist": "Artist Name",
      "playcount": "67",
      "albumArt": "https://...",
      "url": "https://www.last.fm/music/..."
    }
  ]
}
```

### `GET /api/recent-tracks`
Returns recently played tracks from Last.fm

**Query Parameters:**
- `limit` (optional): Number of tracks to return (default: 10, max: 50)

**Example:** `/api/recent-tracks?limit=20`

**Response:**
```json
{
  "tracks": [
    {
      "name": "Track Name",
      "artist": "Artist Name",
      "album": "Album Name",
      "albumArt": "https://...",
      "url": "https://www.last.fm/music/...",
      "playedAt": "2025-10-06T12:00:00Z",
      "isPlaying": false
    }
  ]
}
```

### `GET /api/stats`
Returns your Last.fm listening statistics

**Response:**
```json
{
  "totalScrobbles": "15423",
  "accountAge": "8760h0m0s",
  "username": "your_username"
}
```

### `GET /api/pinned-repo`
Returns a pinned GitHub repository

**Query Parameters:**
- `repo` (optional): Specific repository name to fetch

**Example:** `/api/pinned-repo?repo=tringldev-server`

**Response:**
```json
{
  "name": "repo-name",
  "fullName": "username/repo-name",
  "description": "Repository description",
  "url": "https://github.com/username/repo-name",
  "language": "Go",
  "stars": 42,
  "forks": 5,
  "topics": ["api", "golang"],
  "updatedAt": "2025-10-06T12:00:00Z",
  "homepage": "https://example.com"
}
```

### `POST /api/contact`
Sends a contact form message via Discord webhook

**Form Data:**
- `name` (required): Sender's name
- `email` (optional): Sender's email
- `message` (required): Message content

**Example:**
```bash
curl -X POST http://localhost:8080/api/contact \
  -d "name=John Doe&email=john@example.com&message=Hello!"
```

**Response (HTML for HTMX):**
```html
<div class="success-message">Message sent successfully! I'll get back to you soon.</div>
```

Messages are sent to your Discord channel via webhook.

**Rate Limit:** 5 requests per minute per IP address

## Rate Limiting

All API endpoints are protected with rate limiting:

- **General Endpoints** (`/api/now-playing`, `/api/pinned-repo`): 60 requests per minute (burst of 10)
- **Contact Form** (`/api/contact`): 5 requests per minute (burst of 5)

When rate limit is exceeded, you'll receive a `429 Too Many Requests` response:
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

- `message` (required): Message content

**Example:**
```bash
curl -X POST http://localhost:8080/api/contact \
  -d "name=John Doe&email=john@example.com&message=Hello!"
```

**Response (HTML for HTMX):**
```html
<div class="success-message">Message sent successfully! I'll get back to you soon.</div>
```

Messages are sent to your Discord channel via webhook.

## Setup

### 1. Clone and Install Dependencies

```bash
go mod download
```

### 2. Configure Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` and fill in your credentials:

### 3. Run the Server

```bash
go run .
```

Or build and run:
```bash
go build -o tringldev-server
./tringldev-server
```

The server will start on `http://localhost:8080` (or the port specified in `.env`)

## Development

### Project Structure
```
tringldev-server/
├── cmd/
│   └── server/
│       └── main.go              # Main application entry point & routes
└── internal/
    ├── config/
    │   └── config.go            # Configuration management
    ├── middleware/
    │   └── ratelimit.go         # Rate limiting middleware
    ├── lastfm/
    │   └── service.go           # Last.fm API service
    ├── github/
    │   └── service.go           # GitHub API service
    └── contact/
        └── service.go           # Contact form service (Discord webhook)
```

### Testing the API

Test the endpoints using curl:

```bash
# Health check
curl http://localhost:8080/

# Now playing
curl http://localhost:8080/api/now-playing

# Top weekly artists
curl http://localhost:8080/api/top-artists

# Pinned repo (most recent)
curl http://localhost:8080/api/pinned-repo

# Specific repo
curl http://localhost:8080/api/pinned-repo?repo=your-repo-name

# Contact form
curl -X POST http://localhost:8080/api/contact \
  -d "name=Test&email=test@example.com&message=Hello!"
```

## CORS Configuration

By default, CORS is enabled for all origins (`*`). For production, update the CORS configuration in `main.go` to only allow your frontend domain.

## Deployment
Quick deploy to Fly.io:

```bash
fly launch
fly secrets set LASTFM_API_KEY=your_key LASTFM_USERNAME=your_username
fly secrets set GITHUB_TOKEN=your_token GITHUB_USERNAME=your_username
fly secrets set DISCORD_WEBHOOK=your_webhook_url
fly deploy
```