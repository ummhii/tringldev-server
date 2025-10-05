package main

import (
	"log"

	"tringldev-server/internal/config"
	"tringldev-server/internal/github"
	"tringldev-server/internal/lastfm"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
)

func main() {
	cfg := config.Load()

	lastfmService := lastfm.NewService(cfg)
	githubService := github.NewService(cfg)

	app := iris.New()

	crs := cors.New().
		AllowOrigin("*").
		Handler()

	app.UseRouter(crs)

	// Health check endpoint
	app.Get("/", func(ctx iris.Context) {
		err := ctx.JSON(iris.Map{
			"status":  "ok",
			"message": "TringlDev API Server",
			"version": "1.0.0",
		})
		if err != nil {
			log.Printf("Failed to send response: %v\n", err)
			return
		}
	})

	// Last.fm endpoint - Get currently playing song
	app.Get("/api/now-playing", func(ctx iris.Context) {
		nowPlaying, err := lastfmService.GetCurrentlyPlaying()
		if err != nil {
			log.Printf("Error fetching now playing: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			err := ctx.JSON(iris.Map{
				"error": "Failed to fetch currently playing song",
			})
			if err != nil {
				log.Printf("Failed to send error response: %v\n", err)
			}
			return
		}

		err = ctx.JSON(nowPlaying)
		if err != nil {
			log.Printf("Failed to send response: %v\n", err)
		}
	})

	// GitHub endpoint - Get pinned repository
	app.Get("/api/pinned-repo", func(ctx iris.Context) {
		// Optional: Get specific repo name from query parameter
		repoName := ctx.URLParam("repo")

		var pinnedRepo *github.PinnedRepo
		var err error

		if repoName != "" {
			pinnedRepo, err = githubService.GetSpecificRepository(repoName)
		} else {
			pinnedRepo, err = githubService.GetPinnedRepository()
		}

		if err != nil {
			log.Printf("Error fetching pinned repo: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			err := ctx.JSON(iris.Map{
				"error": "Failed to fetch pinned repository",
			})
			if err != nil {
				log.Printf("Failed to send error response: %v\n", err)
			}
			return
		}

		err = ctx.JSON(pinnedRepo)
		if err != nil {
			log.Printf("Failed to send response: %v\n", err)
		}
	})

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s\n", addr)
	err := app.Run(iris.Addr(addr))
	if err != nil {
		log.Printf("Failed to start server: %v\n", err)
		return
	}
}
