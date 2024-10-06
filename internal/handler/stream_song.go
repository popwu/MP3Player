package handler

import (
	"mp3-player/internal/libs"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
