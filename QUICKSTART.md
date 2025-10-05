# Quick Start Guide

Get your API running in 3 minutes!

## 1️⃣ Install Dependencies

```powershell
go mod download
```

## 2️⃣ Set Up Last.fm (Super Easy!)

### Get Last.fm API Key
1. Go to https://www.last.fm/api/account/create
2. Fill in:
   - Name: `My Personal Website`
   - Description: `Personal website now playing`
3. Click "Submit" and copy your **API Key**

### Add to .env file
Create a `.env` file:
```env
PORT=8080
LASTFM_API_KEY=paste_your_api_key_here
LASTFM_USERNAME=your_lastfm_username
```

### Connect Your Music Service
1. Go to https://www.last.fm/settings/applications
2. Connect Spotify, Apple Music, or whatever you use
3. Play some music!

## Set Up GitHub (Optional)

1. Go to https://github.com/settings/tokens
2. Generate a new token (classic)
3. Select scope: `public_repo`
4. Copy the token

Add to `.env`:
```env
GITHUB_TOKEN=paste_your_github_token_here
GITHUB_USERNAME=your_github_username
```

## 4️⃣ Start the Server

```powershell
go run .
```

## 5️⃣ Test Your API

### Check if server is running:
```powershell
curl http://localhost:8080/
```

### Get currently playing song:
```powershell
curl http://localhost:8080/api/now-playing
```

### Get pinned repository:
```powershell
curl http://localhost:8080/api/pinned-repo
```