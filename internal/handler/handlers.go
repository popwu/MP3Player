package api

import (
	"net/http"

	"mp3-player/internal/libs"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	storage   *libs.FileStorage
	cache     *libs.Cache
	cacheSize int
}

func NewHandlers(storage *libs.FileStorage, cache *libs.Cache, cacheSize int) *Handlers {
	return &Handlers{storage: storage, cache: cache, cacheSize: cacheSize}
}

// func (h *Handlers) GetLibraries(c *gin.Context) {
// 	libraries, err := h.storage.GetLibraries()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, libraries)
// }

// func (h *Handlers) AddLibrary(c *gin.Context) {
// 	path := c.PostForm("libraryPath")
// 	if path == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "路径不能为空"})
// 		return
// 	}

// 	if err := h.storage.AddLibrary(path); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	libraries, err := h.storage.GetLibraries()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 返回更新后的库列表
// 	c.HTML(http.StatusOK, "libraries.html", gin.H{"libraries": libraries})
// }

// func (h *Handlers) GetSongs(c *gin.Context) {
// 	songs, err := h.storage.GetSongs()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, songs)
// }

func (h *Handlers) StreamSong(c *gin.Context) {
	songPath := c.Param("songPath")
	authToken := c.GetHeader("Authorization")

	// 检查缓存
	if cachedData, exists := h.cache.Get(songPath); exists {
		c.Data(http.StatusOK, "audio/mpeg", cachedData)
		return
	}

	// 从云端获取文件
	data, err := libs.GetFileFromCloud(authToken, songPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法从云端获取文件"})
		return
	}

	// 将文件添加到缓存
	h.cache.Set(songPath, data)

	c.Data(http.StatusOK, "audio/mpeg", data)
}
