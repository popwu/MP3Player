package handler

import (
	"mp3-player/internal/libs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) LoginCheck(c *gin.Context) {
	accessKey := c.PostForm("qiniu_access_key")
	secretKey := c.PostForm("qiniu_secret_key")
	bucket := c.PostForm("qiniu_bucket")
	zone := c.PostForm("qiniu_zone")

	status, message, fileList := libs.CheckQiniuAccess(accessKey, secretKey, bucket, zone)

	if status == http.StatusOK {
		// 生成 auth_token
		authToken := libs.GenerateAuthToken()

		// 保存配置到文件
		config := libs.QiniuConfig{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Bucket:    bucket,
			Zone:      zone,
		}
		if err := libs.SaveConfig(authToken, config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "无法保存配置"})
			return
		}

		// 保存文件列表到本地
		if err := h.storage.SaveFileList(fileList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "无法保存文件列表"})
			return
		}

		c.JSON(status, gin.H{
			"status":  "success",
			"message": message,
			"data": gin.H{
				"auth_token": authToken,
				"filelist":   fileList,
			},
		})
	} else {
		c.JSON(status, gin.H{"status": "error", "message": message})
	}
}
