package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func (h *Handlers) LoginCheck(c *gin.Context) {
	accessKey := c.PostForm("qiniu_access_key")
	secretKey := c.PostForm("qiniu_secret_key")
	bucket := c.PostForm("qiniu_bucket")
	zone := c.PostForm("qiniu_zone")

	status, message := checkQiniuAccess(accessKey, secretKey, bucket, zone)
	c.JSON(status, gin.H{"status": "ok", "message": message})
}

func checkQiniuAccess(accessKey, secretKey, bucket, zone string) (int, string) {
	// 检查参数是否为空
	if accessKey == "" || secretKey == "" || bucket == "" {
		return http.StatusBadRequest, "请提供 accessKey, secretKey 和 bucket"
	}

	// 创建七牛云认证对象
	mac := auth.New(accessKey, secretKey)

	// 创建存储空间管理器
	cfg := storage.Config{
		Zone: getZone(zone),
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	// 列出存储空间中的文件
	entries, _, _, _, err := bucketManager.ListFiles(bucket, "", "", "", 1)
	if err != nil {
		return http.StatusUnauthorized, "认证失败或无法访问存储空间"
	}

	if len(entries) > 0 {
		return http.StatusOK, "成功访问存储空间"
	} else {
		return http.StatusOK, "存储空间为空"
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
