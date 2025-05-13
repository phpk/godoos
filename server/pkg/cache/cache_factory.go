package cache

import (
	"errors"
	"time"
)

// 初始化缓存
// cache, err := cache.NewCache(cache.CacheConfig{
// 	Type: cache.FileCacheType,
// 	FileDir: "data/cache",
// 	})
// 	if err != nil {
// 	log.Fatal(err)
// 	}

// 	// 设置缓存(有效期10分钟)
// 	err = cache.Set("user:1", userData, 10*time.Minute)

// 	// 获取缓存
// 	val, err := cache.Get("user:1")

//	// 删除缓存
//	err = cache.Delete("user:1")

// CacheItem 缓存项结构
type CacheItem struct {
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expiresAt"`
}

// Cache 定义缓存接口
type Cache interface {
	// Set 设置缓存
	Set(key string, value interface{}, ttl time.Duration) error
	// Get 获取缓存
	Get(key string) (interface{}, error)
	GetKey(key string, val string) (interface{}, error)
	// Delete 删除缓存
	Delete(key string) error
	// Clear 清空缓存
	Clear() error
	// Exists 检查缓存是否存在
	Exists(key string) (bool, error)
}

var (
	ErrCacheMiss    = errors.New("cache: key not found")
	ErrCacheExpired = errors.New("cache: key expired")
)

// NewCache 创建缓存实例
func NewCache(cacheType string, fileDir string) Cache {
	switch cacheType {
	case "memory":
		return NewMemoryCache()
	case "file":
		if fileDir == "" {
			return NewMemoryCache()
		}
		return NewFileCache(fileDir)
	default:
		return NewMemoryCache()
	}
}
