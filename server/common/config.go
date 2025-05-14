package common

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

func LoadSystemInfo() error {
	data, err := os.ReadFile("config/system.json")
	if err != nil {
		return err
	}
	json.Unmarshal(data, &Config)
	loginData, err := os.ReadFile("config/login.json")
	if err != nil {
		return err
	}
	json.Unmarshal(loginData, &LoginConf)
	Cache = cache.NewCache(Config.Cache.Type, Config.Cache.FileDir)
	return nil
}
