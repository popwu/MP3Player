package libs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	config, err := getQiniuConfig(authToken)
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
func getQiniuConfig(authToken string) (QiniuConfig, error) {
	return loadConfig(authToken)
}
