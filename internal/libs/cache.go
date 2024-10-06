package libs

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Cache struct {
	items map[string][]byte
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string][]byte),
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.items[key]
	return value, exists
}

func ParseCacheSize(size string) (int64, error) {
	size = strings.TrimSpace(size)
	if len(size) == 0 {
		return 0, fmt.Errorf("缓存大小不能为空")
	}

	unit := size[len(size)-1]
	value, err := strconv.ParseInt(size[:len(size)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的缓存大小格式")
	}

	switch unicode.ToUpper(rune(unit)) {
	case 'K':
		return value * 1024, nil
	case 'M':
		return value * 1024 * 1024, nil
	case 'G':
		return value * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("不支持的单位: %c", unit)
	}
}
