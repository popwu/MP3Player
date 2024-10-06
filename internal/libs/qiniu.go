package libs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Zone      string
}

func GetFileFromCloud(authToken, key string) ([]byte, error) {
	config, err := GetQiniuConfig(authToken)
	if err != nil {
		return nil, err
	}

	mac := auth.New(config.AccessKey, config.SecretKey)
	domain := os.Getenv("QINIU_DOMAIN") // 假设域名存储在环境变量中

	deadline := time.Now().Add(time.Second * 3600).Unix() // 1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	resp, err := http.Get(privateAccessURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func loadConfig(token string) (QiniuConfig, error) {
	var config QiniuConfig
	data, err := os.ReadFile(filepath.Join("./configs", token+".json"))
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

// 在其他需要使用七牛云配置的地方，可以这样获取配置
func GetQiniuConfig(authToken string) (QiniuConfig, error) {
	return loadConfig(authToken)
}

func GenerateAuthToken() string {
	return uuid.New().String()
}

func SaveConfig(token string, config QiniuConfig) error {
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

func CheckQiniuAccess(accessKey, secretKey, bucket, zone string) (int, string, *[]storage.ListItem) {
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
