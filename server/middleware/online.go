package middleware

import (
	"encoding/json"
	"fmt"
	"godocms/config"
	"godocms/libs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// UserActivity 结构体用于存储用户的单个活动数据
type UserActivity struct {
	Duration  int64     `json:"duration"` // 修改为 int64 类型
	Timestamp time.Time `json:"timestamp"`
}

// UserActivities 结构体用于存储用户的多个活动数据
type UserActivities map[string]UserActivity

// UserStatus 结构体用于存储每个 IP 地址的状态和活动数据
type UserStatus struct {
	Activities UserActivities `json:"activities"`
	LastSeen   time.Time      `json:"last_seen"`
}

// UserStatuses 映射表用于存储每个 IP 地址的状态和活动数据
var (
	UserStatusesMutex sync.Mutex
	UserStatuses      sync.Map
)
var ifWrite = false

// func init() {
// 	InitOnlineFile()
// }

// recordUserActivity 中间件记录用户的 IP 地址、请求的 URL 和开始时间
func RecordUserActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path // 获取请求的 URL
		if strings.HasPrefix(url, "/static/") || strings.HasPrefix(url, "/upload/") {
			c.Next() // 继续后续处理
			return
		}
		ip := libs.GetIpAddress(c) // 获取客户端 IP 地址
		currentTime := time.Now()
		//log.Printf("IP:%s,URL:%s,Time:%s", ip, url, currentTime)
		// 检查 URL 是否以 /static/ 开头，如果是则跳过记录

		// 更新用户的活动数据
		UserStatusesMutex.Lock()
		userStatus, ok := UserStatuses.Load(ip)
		if ok {
			// 如果用户已经存在，则更新其停留时长
			status := userStatus.(UserStatus)
			if activity, exists := status.Activities[url]; exists {
				durationInSeconds := int64(currentTime.Sub(activity.Timestamp).Seconds()) // 转换为秒
				activity.Duration += durationInSeconds
				activity.Timestamp = currentTime
				status.Activities[url] = activity
			} else {
				status.Activities[url] = UserActivity{
					Timestamp: currentTime,
					Duration:  0,
				}
			}
			status.LastSeen = currentTime
			userStatus = status
		} else {
			// 如果用户不存在，则创建新的状态记录
			userStatus = UserStatus{
				Activities: UserActivities{
					url: UserActivity{
						Timestamp: currentTime,
						Duration:  0,
					},
				},
				LastSeen: currentTime,
			}
		}

		// 将用户状态数据保存回映射表
		UserStatuses.Store(ip, userStatus)
		ifWrite = true
		UserStatusesMutex.Unlock()
		c.Next() // 继续后续处理
	}
}

// 读取 JSON 文件
func InitOnlineFile() error {
	logPath := getLogPath()
	//log.Printf("logpath is:%s", logPath)
	data, err := os.ReadFile(logPath)
	if os.IsNotExist(err) {
		return nil // 文件不存在，不需要做任何事情
	} else if err != nil {
		return err
	}

	// 创建临时映射表用于反序列化
	var mapped map[string]UserStatus
	if err := json.Unmarshal(data, &mapped); err != nil {
		return err
	}

	// 清空 UserStatuses 映射表
	UserStatusesMutex.Lock()
	UserStatuses = sync.Map{}
	UserStatusesMutex.Unlock()

	// 将数据加载到 UserStatuses 映射表中
	for k, v := range mapped {
		UserStatuses.Store(k, v)
	}
	StartTicker()
	return nil
}

// 启动定时器
func StartTicker() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	go func() {
		for range ticker.C {
			if ifWrite {
				if err := WriteJSONFile(); err != nil {
					fmt.Println("Error writing to file:", err)
				}
				ifWrite = false
			}
		}
	}()
}

// 获取日志文件目录
func getLogPath() string {
	filename := fmt.Sprintf("online_%s.json", time.Now().Format("2006-01-02"))
	// 获取程序执行文件所在目录的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		return filename
	}
	exeDir := filepath.Dir(exePath)

	// 构建 logs 文件夹的绝对路径
	logDir := filepath.Join(exeDir, "logs")
	if !libs.PathExists(logDir) {
		os.Mkdir(logDir, 0755)
	}

	return filepath.Join(logDir, filename)
}

// 写入 JSON 文件
func WriteJSONFile() error {
	// 将 sync.Map 转换为普通的映射表
	mapped := make(map[string]UserStatus)
	UserStatuses.Range(func(key, value interface{}) bool {
		mapped[key.(string)] = value.(UserStatus)
		return true
	})

	data, err := json.MarshalIndent(mapped, "", "  ")
	if err != nil {
		return err
	}

	// 将文件写入 logs 目录
	logPath := getLogPath()
	return os.WriteFile(logPath, data, 0644)
}

func UserAccessRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt64("userId")
		if uid == 0 {
			c.Next()
			return
		}

		UserStatusesMutex.Lock()
		// 获取缓存的数据
		var origin map[int64]int64
		var ok bool
		data, _ := config.Cache.Get("user_online")
		if data == nil {
			origin = make(map[int64]int64)
		} else {
			if origin, ok = data.(map[int64]int64); !ok {
				origin = make(map[int64]int64)
			}
		}

		origin[uid] = time.Now().Unix()
		config.Cache.Set("user_online", origin, 60*24)
		UserStatusesMutex.Unlock()

		c.Next()
	}
}

// 启动定时心跳检查任务
func StartHeartbeatCheckOnline() {
	ticker := time.NewTicker(3 * time.Minute) // 每 3 分钟执行一次
	defer ticker.Stop()

	for range ticker.C {
		// 每次心跳检查，检查不活跃的用户
		checkInactiveUsers()
	}
}

// 检查不在线的用户
func checkInactiveUsers() {
	UserStatusesMutex.Lock() // 锁定对缓存的访问
	defer UserStatusesMutex.Unlock()

	now := time.Now().Unix()
	// 获取缓存的数据
	var origin map[int64]int64
	var ok bool
	data, _ := config.Cache.Get("user_online")
	if data == nil {
		log.Println("No user activity data found in cache.")
		return
	}

	// 安全地进行类型断言
	if origin, ok = data.(map[int64]int64); !ok {
		log.Println("Cache data is not of the expected type.")
		return
	}

	// 遍历所有用户的活动记录，检查是否有不在线的用户
	for uid, lastActivityTime := range origin {
		if now-lastActivityTime > 30*60 { // 如果用户的最后活动时间超过 30 分钟, 视为不在线
			delete(origin, uid)
		}
	}

	config.Cache.Set("user_online", origin, 60*24)
}
