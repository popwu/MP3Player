package api

import (
	"encoding/json"
	"fmt"
	"mp3-player/internal/libs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func (h *Handlers) LoginCheck(c *gin.Context) {
	accessKey := c.PostForm("qiniu_access_key")
	secretKey := c.PostForm("qiniu_secret_key")
	bucket := c.PostForm("qiniu_bucket")
	zone := c.PostForm("qiniu_zone")

	status, message, fileList := checkQiniuAccess(accessKey, secretKey, bucket, zone)

	if status == http.StatusOK {
		// 生成 auth_token
		authToken := generateAuthToken()

		// 保存配置到文件
		config := libs.QiniuConfig{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Bucket:    bucket,
			Zone:      zone,
		}
		if err := saveConfig(authToken, config); err != nil {
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

func generateAuthToken() string {
	return uuid.New().String()
}

func saveConfig(token string, config libs.QiniuConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	configDir := "./configs"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, token+".json"), data, 0644)
}

func checkQiniuAccess(accessKey, secretKey, bucket, zone string) (int, string, *[]storage.ListItem) {
	// 检查参数是否为空
	if accessKey == "" || secretKey == "" || bucket == "" {
		return http.StatusBadRequest, "请提供 accessKey, secretKey 和 bucket", nil
	}

	// 创建七牛云认证对象
	mac := auth.New(accessKey, secretKey)

	// 创建存储空间管理器
	cfg := storage.Config{
		Zone: getZone(zone),
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	// 列出存储空间中的所有文件
	var fileList []storage.ListItem
	marker := ""
	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, "", "", marker, 1000)
		if err != nil {
			return http.StatusUnauthorized, "认证失败或无法访问存储空间", nil
		}
		fileList = append(fileList, entries...)
		if !hasNext {
			break
		}
		marker = nextMarker
	}

	if len(fileList) > 0 {
		return http.StatusOK, "成功访问存储空间", &fileList
	} else {
		return http.StatusOK, "存储空间为空", &fileList
	}
}

func getZone(zone string) *storage.Zone {
	zoneMap := map[string]*storage.Zone{
		"huadong":  &storage.ZoneHuadong,
		"huabei":   &storage.ZoneHuabei,
		"huanan":   &storage.ZoneHuanan,
		"beimei":   &storage.ZoneBeimei,
		"xinjiapo": &storage.ZoneXinjiapo,
	}

	if z, ok := zoneMap[zone]; ok {
		return z
	}
	fmt.Printf("警告：未知的区域 '%s'，使用默认区域（华东）\n", zone)
	return &storage.ZoneHuadong
}
