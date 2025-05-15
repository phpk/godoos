package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// fileCache 文件缓存实现
type fileCache struct {
	mu  sync.RWMutex
	dir string
}

// NewFileCache 创建文件缓存实例
func NewFileCache(dir string) *fileCache {
	return &fileCache{
		dir: dir,
	}
}

func (c *fileCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}

	// 强制将 Value 转换为 []byte
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	item := struct {
		Data      []byte    `json:"data"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}

	encoded, err := json.Marshal(item)
	if err != nil {
		return err
	}

	filePath := filepath.Join(c.dir, key)
	return os.WriteFile(filePath, encoded, 0644)
}

type FileCacheItem struct {
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (c *fileCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := filepath.Join(c.dir, key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var item FileCacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	if item.ExpiresAt.Before(time.Now()) {
		_ = os.Remove(filePath)
		return nil, ErrCacheExpired
	}

	return item.Data, nil
}
func (c *fileCache) GetKey(key string, clientId string) (interface{}, error) {
	return c.Get(key + ":" + clientId)
}
func (c *fileCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	filePath := filepath.Join(c.dir, key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func (c *fileCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_ = os.Remove(filepath.Join(c.dir, entry.Name()))
	}
	return nil
}

func (c *fileCache) Exists(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := filepath.Join(c.dir, key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	// 检查是否过期
	val, err := c.Get(key)
	if err == ErrCacheExpired {
		return false, nil
	}
	return val != nil, err
}

// PathExists 检查路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
