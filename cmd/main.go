package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	api "mp3-player/internal/handler"
	"mp3-player/internal/libs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 获取缓存大小
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Fatal("Invalid CACHE_SIZE in .env file")
	}

	r := gin.Default()

	fileStorage := libs.NewFileStorage("music_library.json")
	musicCache := libs.NewCache()
	handlers := api.NewHandlers(fileStorage, musicCache, cacheSize)

	// r.GET("/api/libraries", handlers.GetLibraries)
	// r.POST("/api/libraries", handlers.AddLibrary)
	// r.GET("/api/songs", handlers.GetSongs)
	r.GET("/api/stream/:songPath", handlers.StreamSong)
	r.GET("/api/logincheck", handlers.LoginCheck)

	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Fatal(r.Run(":8080"))
}
