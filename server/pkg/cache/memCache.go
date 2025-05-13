package cache

import (
	"sync"
	"time"
)

// memoryCache 内存缓存实现
type memoryCache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *memoryCache {
	return &memoryCache{
		data: make(map[string]CacheItem),
	}
}

func (c *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	// 设置过期清理
	time.AfterFunc(ttl, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		if item, ok := c.data[key]; ok && item.ExpiresAt.Before(time.Now()) {
			delete(c.data, key)
		}
	})

	return nil
}

func (c *memoryCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	if item.ExpiresAt.Before(time.Now()) {
		return nil, ErrCacheExpired
	}
	return item.Value, nil
}
func (c *memoryCache) GetKey(key string, clientId string) (interface{}, error) {
	return c.Get(key + ":" + clientId)
}
func (c *memoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *memoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheItem)
	return nil
}

func (c *memoryCache) Exists(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return false, nil
	}
	if item.ExpiresAt.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}
