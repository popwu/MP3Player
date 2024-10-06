package api

import (
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
		// 注意：以下测试需要有效的七牛云凭证才能通过
		// {
		// 	name:           "有效凭证_空存储空间",
		// 	accessKey:      "valid_access_key",
		// 	secretKey:      "valid_secret_key",
		// 	bucket:         "empty_bucket",
		// 	zone:           "huadong",
		// 	expectedStatus: http.StatusOK,
		// 	expectedMsg:    "存储空间为空",
		// },
		// {
		// 	name:           "有效凭证_非空存储空间",
		// 	accessKey:      "valid_access_key",
		// 	secretKey:      "valid_secret_key",
		// 	bucket:         "non_empty_bucket",
		// 	zone:           "huadong",
		// 	expectedStatus: http.StatusOK,
		// 	expectedMsg:    "成功访问存储空间",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := checkQiniuAccess(tt.accessKey, tt.secretKey, tt.bucket, tt.zone)
			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedMsg, msg)
		})
	}
}
