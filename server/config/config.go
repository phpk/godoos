package config

import (
	"encoding/json"
	"godocms/pkg/cache"
	"os"
)

var Config Info
var Cache cache.Cache

type Info struct {
	System System      `json:"system"`
	Log    Log         `json:"log"`
	Email  Email       `json:"email"`
	Cache  CacheConfig `json:"cache"`
}
type System struct {
	Debug         bool   `json:"debug"`
	Port          int    `json:"port"`
	Host          string `json:"host"`
	Scheme        string `json:"scheme"`
	SessionType   string `json:"sessionType"`
	SessionSecret string `json:"sessionSecret"`

	IpAccess []string `json:"ipAccess"`
}
type CacheConfig struct {
	Type    string `json:"type"`    // 缓存类型: memory/file
	FileDir string `json:"fileDir"` // 文件缓存目录(仅文件缓存需要)
}
type Log struct {
	WriteFile  bool   `json:"writeFile"`
	Path       string `json:"path"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
}
type Email struct {
	From         string `json:"from"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	Port         int    `json:"port"`
	IsSsl        bool   `json:"isSsl"`
	UserLogin    bool   `json:"userLogin"`
	AutoRegister bool   `json:"autoRegister"`
	Enable       bool   `json:"enable"`
	AdminLogin   bool   `json:"adminLogin"`
}

func LoadConfig() error {
	data, err := os.ReadFile("config/system.json")
	if err != nil {
		return err
	}
	json.Unmarshal(data, &Config)
	Cache = cache.NewCache(Config.Cache.Type, Config.Cache.FileDir)
	return nil
}
