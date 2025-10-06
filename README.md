# tringl.dev backend

A personal backend API for displaying currently playing songs and pinned GitHub projects and other stuff.

## Features

- **Now Playing**: Displays the current song you're listening to via Last.fm (works with Spotify, Apple Music, YouTube Music, etc.)
- **Pinned Repository**: Shows one of your GitHub repositories
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
Sends a contact form message via Discord or email

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
.
├── cmd/
│   └── server/
│       └── main.go      # Main application entry point & routes
└── internal/
    ├── config/
    │   └── config.go    # Configuration management
    ├── lastfm/
    │   └── service.go   # Last.fm API service
    ├── github/
    │   └── service.go   # GitHub API service
    └── contact/
        └── service.go   # Contact form service (Discord webhook)
```

### Testing the API

Test the endpoints using curl:

```bash
# Health check
curl http://localhost:8080/

# Now playing
curl http://localhost:8080/api/now-playing

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