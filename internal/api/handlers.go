package api

import (
	"fmt"
	"net/http"
	"os"

	"mp3-player/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	storage *storage.FileStorage
}

func NewHandlers(storage *storage.FileStorage) *Handlers {
	return &Handlers{storage: storage}
}

func (h *Handlers) GetLibraries(c *gin.Context) {
	libraries, err := h.storage.GetLibraries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, libraries)
}

func (h *Handlers) AddLibrary(c *gin.Context) {
	path := c.PostForm("libraryPath")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径不能为空"})
		return
	}

	if err := h.storage.AddLibrary(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	libraries, err := h.storage.GetLibraries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的库列表
	c.HTML(http.StatusOK, "libraries.html", gin.H{"libraries": libraries})
}

func (h *Handlers) GetSongs(c *gin.Context) {
	songs, err := h.storage.GetSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, songs)
}

func (h *Handlers) StreamSong(c *gin.Context) {
	songPath := c.Param("songPath")

	// 验证文件路径
	if !h.storage.IsValidSongPath(songPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的歌曲路径"})
		return
	}

	// 打开文件
	file, err := os.Open(songPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取文件信息"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Accept-Ranges", "bytes")

	// 将文件内容写入响应
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
}
