package tests

import (
	"fmt"
	"mp3-player/internal/libs"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func loadEnv() {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			godotenv.Load(filepath.Join(dir, ".env"))
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func TestCheckQiniuAccess(t *testing.T) {
	loadEnv()

	// 打印环境变量
	fmt.Println("QINIU_ACCESS_KEY:", os.Getenv("QINIU_ACCESS_KEY"))
	fmt.Println("QINIU_SECRET_KEY:", os.Getenv("QINIU_SECRET_KEY"))
	fmt.Println("QINIU_BUCKET:", os.Getenv("QINIU_BUCKET"))
	fmt.Println("QINIU_ZONE:", os.Getenv("QINIU_ZONE"))

	// 保存原始环境变量
	originalAccessKey := os.Getenv("QINIU_ACCESS_KEY")
	originalSecretKey := os.Getenv("QINIU_SECRET_KEY")
	originalBucket := os.Getenv("QINIU_BUCKET")
	originalZone := os.Getenv("QINIU_ZONE")

	// 测试结束后恢复环境变量
	defer func() {
		os.Setenv("QINIU_ACCESS_KEY", originalAccessKey)
		os.Setenv("QINIU_SECRET_KEY", originalSecretKey)
		os.Setenv("QINIU_BUCKET", originalBucket)
		os.Setenv("QINIU_ZONE", originalZone)
	}()

	tests := []struct {
		name           string
		accessKey      string
		secretKey      string
		bucket         string
		zone           string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "缺少参数",
			accessKey:      "",
			secretKey:      "",
			bucket:         "",
			zone:           "",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "请提供 accessKey, secretKey 和 bucket",
		},
		{
			name:           "无效凭证",
			accessKey:      "invalid_access_key",
			secretKey:      "invalid_secret_key",
			bucket:         "test_bucket",
			zone:           "huadong",
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "认证失败或无法访问存储空间",
		},
		{
			name:           "有效凭证",
			accessKey:      os.Getenv("QINIU_ACCESS_KEY"),
			secretKey:      os.Getenv("QINIU_SECRET_KEY"),
			bucket:         os.Getenv("QINIU_BUCKET"),
			zone:           os.Getenv("QINIU_ZONE"),
			expectedStatus: http.StatusOK,
			expectedMsg:    "成功访问存储空间",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("运行测试: %s\n", tt.name)
			fmt.Printf("使用的参数: accessKey=%s, secretKey=%s, bucket=%s, zone=%s\n",
				tt.accessKey, tt.secretKey, tt.bucket, tt.zone)

			status, msg, err := libs.CheckQiniuAccess(tt.accessKey, tt.secretKey, tt.bucket, tt.zone)

			fmt.Printf("得到的结果: status=%d, msg=%s, err=%v\n", status, msg, err)
			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedMsg, msg)
		})
	}
}
