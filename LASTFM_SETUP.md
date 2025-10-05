# 🎵 Last.fm API Setup Guide

This guide will help you get your Last.fm API key in just 2 minutes - it's super simple!

## Why Last.fm?

Last.fm is **much easier** than Spotify:
- ✅ No OAuth required - just an API key
- ✅ Works with any music source (Spotify, Apple Music, YouTube Music, etc.)
- ✅ Free and no restrictions
- ✅ Automatic scrobbling from your music apps

## Step 1: Create Last.fm Account

1. Go to [Last.fm](https://www.last.fm) and create an account (if you don't have one)
2. Or log in to your existing account

## Step 2: Get API Key

1. Go to [Last.fm API Account Creation](https://www.last.fm/api/account/create)
2. Fill in the application details:
   - **Application name**: `My Personal Website` (or any name)
   - **Application description**: `Personal website now playing integration`
   - **Callback URL**: Can leave blank or use `http://localhost:8080`
3. Click **Submit**
4. You'll see your **API Key** and **Shared Secret** (you only need the API Key)

## Step 3: Add to .env File

Create or edit your `.env` file:

```env
PORT=8080
LASTFM_API_KEY=your_api_key_here
LASTFM_USERNAME=your_lastfm_username
```

**Replace:**
- `your_api_key_here` with the API key from step 2
- `your_lastfm_username` with your Last.fm username (visible in your profile URL)

## Step 4: Connect Your Music Apps

Last.fm can scrobble (track) music from various sources:

### For Spotify Users:
1. Go to [Last.fm Settings → Applications](https://www.last.fm/settings/applications)
2. Click "Connect" next to Spotify
3. Authorize the connection
4. Your Spotify listening will now be tracked!

### For Other Music Services:
- **Apple Music**: Download [Last.fm Scrobbler](https://www.last.fm/about/trackmymusic)
- **YouTube Music**: Use browser extension like [Web Scrobbler](https://web-scrobbler.com/)
- **Local Files**: Use desktop app for [Windows](https://www.last.fm/about/trackmymusic) or [Mac](https://www.last.fm/about/trackmymusic)

## Step 5: Start Your Server

```powershell
go run .
```

## Step 6: Test It Out

Play some music on Spotify (or your connected service), then:

```powershell
curl http://localhost:8080/api/now-playing
```

You should see your currently playing track! 🎉

## Response Format

```json
{
  "isPlaying": true,
  "songName": "Song Title",
  "artistName": "Artist Name",
  "albumName": "Album Name",
  "albumArt": "https://lastfm.freetls.fastly.net/...",
  "songUrl": "https://www.last.fm/music/...",
  "lastUpdated": "2025-10-06T12:00:00Z"
}
```

If not currently playing, you'll get the last played track with `isPlaying: false` and a `playedAt` timestamp.

## Troubleshooting

### "No tracks found"
- Make sure you've connected your music service to Last.fm
- Play some music and wait a few seconds for it to scrobble
- Check if tracks appear on your Last.fm profile

### "Invalid API Key"
- Double-check the API key in your `.env` file
- Make sure there are no extra spaces or quotes
- Restart your server after updating `.env`

### "User not found"
- Verify your Last.fm username is correct
- Check your profile URL: `https://www.last.fm/user/YOUR_USERNAME`

## Privacy Settings

By default, your Last.fm profile is public. To change:
1. Go to [Privacy Settings](https://www.last.fm/settings/privacy)
2. Adjust "Recent listening information" as desired
3. Note: The API will only work if your recent tracks are visible

## Advantages Over Spotify API

| Feature | Last.fm | Spotify |
|---------|---------|---------|
| Setup Complexity | ⭐ Simple (just API key) | ⭐⭐⭐ Complex (OAuth flow) |
| Music Source | Any (Spotify, Apple, etc.) | Spotify only |
| Token Refresh | Not needed | Every hour |
| Rate Limits | Very generous | Strict |
| Free Tier | Yes | Yes |

## Need Help?

- [Last.fm API Documentation](https://www.last.fm/api)
- [Last.fm Support](https://support.last.fm/)
