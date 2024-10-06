package main

import (
	"fmt"
	"log"
	"mp3-player/internal/handler"
	"mp3-player/internal/libs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 获取缓存大小
	cacheSize, err := libs.ParseCacheSize(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Fatal("无效的 CACHE_SIZE 在 .env 文件中")
	}

	r := gin.Default()

	fileStorage := libs.NewFileStorage("music_library.json")
	musicCache := libs.NewCache()
	handlers := handler.NewHandlers(fileStorage, musicCache, cacheSize)

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

func parseCacheSize(size string) (int64, error) {
	size = strings.TrimSpace(size)
	if len(size) == 0 {
		return 0, fmt.Errorf("缓存大小不能为空")
	}

	unit := size[len(size)-1]
	value, err := strconv.ParseInt(size[:len(size)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的缓存大小格式")
	}

	switch unicode.ToUpper(rune(unit)) {
	case 'K':
		return value * 1024, nil
	case 'M':
		return value * 1024 * 1024, nil
	case 'G':
		return value * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("不支持的单位: %c", unit)
	}
}
