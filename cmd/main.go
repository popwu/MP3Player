package main

import (
	"log"
	"net/http"

	"mp3-player/internal/api"
	"mp3-player/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	fileStorage := storage.NewFileStorage("music_library.json")
	handlers := api.NewHandlers(fileStorage)

	r.GET("/api/libraries", handlers.GetLibraries)
	r.POST("/api/libraries", handlers.AddLibrary)
	r.GET("/api/songs", handlers.GetSongs)
	r.GET("/api/stream/:songPath", handlers.StreamSong)

	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Fatal(r.Run(":8080"))
}
