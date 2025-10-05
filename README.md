# TringlDev Server API

A personal backend API for displaying currently playing Spotify songs and pinned GitHub projects.

## Features

- **Now Playing**: Displays the current song you're listening to via Last.fm (works with Spotify, Apple Music, YouTube Music, etc.)
- **Pinned Repository**: Shows one of your GitHub repositories
- **CORS Enabled**: Ready to be consumed by your frontend
- **Simple Setup**: No OAuth complexity - just an API key!

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

#### Last.fm Setup:
1. Go to [Last.fm API Account Creation](https://www.last.fm/api/account/create)
2. Create an app and get your API key
3. Add API key and your Last.fm username to your `.env` file
4. Connect your music service (Spotify, Apple Music, etc.) to Last.fm
5. Start listening to music!

See [LASTFM_SETUP.md](LASTFM_SETUP.md) for detailed instructions.

#### GitHub Setup:
1. Go to [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)
2. Generate a new token (classic)
3. Select scopes: `public_repo` (for public repositories)
4. Copy the token and your username

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
├── main.go          # Main application entry point
├── config.go        # Configuration loading
├── spotify.go       # Spotify API service
├── github.go        # GitHub API service
├── .env.example     # Example environment variables
└── README.md        # This file
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
```

## CORS Configuration

By default, CORS is enabled for all origins (`*`). For production, update the CORS configuration in `main.go` to only allow your frontend domain.

## Troubleshooting

- **Spotify Token Issues**: Make sure your refresh token is valid and hasn't expired
- **GitHub Rate Limiting**: Use a GitHub token to increase API rate limits
- **Missing Environment Variables**: Check the console logs for warnings about missing configuration

## License

MIT
