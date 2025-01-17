package user

import (
	"encoding/json"
	"godo/libs"
	"godo/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// 全局变量，用于存储锁屏状态
var userLockedScreen map[uint]bool
var mu sync.Mutex // 用于保护锁屏状态的并发安全

// 注册系统用户
func RegisterSysUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		libs.ErrorMsg(w, "invalid request method")
		return
	}

	var user model.SysUser

	// 获取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("读取请求体错误:", err)
		libs.ErrorMsg(w, "invalid input")
		return
	}

	// 解析请求体
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("解析请求体错误:", err)
		libs.ErrorMsg(w, "invalid input")
		return
	}

	// 检查用户名和密码是否为空
	if user.Username == "" || user.Password == "" {
		libs.ErrorMsg(w, "username or password is empty")
		return
	}

	// 创建用户
	if err := model.Db.Create(&user).Error; err != nil {
		libs.ErrorMsg(w, "failed to create user")
		return
	}

	// 返回成功消息
	libs.SuccessMsg(w, user.ID, "")
}

// 系统用户登录（锁屏）
func LockedScreenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		libs.ErrorMsg(w, "invalid request method")
		return
	}

	var req model.SysUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		libs.ErrorMsg(w, "invalid input")
		return
	}

	var user model.SysUser
	if err := model.Db.Where("username = ? and password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		libs.ErrorMsg(w, "invalid username or password")
		return
	}

	userEncodeStr, err := libs.EncodeFile("godoos", strconv.Itoa(int(user.ID)))
	if err != nil {
		libs.ErrorMsg(w, "failed to encode user")
		return
	}

	// w.Header().Set("sysuser", userEncodeStr)

	// 登录成功，锁屏状态设置为 true
	mu.Lock()
	if userLockedScreen == nil {
		userLockedScreen = make(map[uint]bool)
	}
	userLockedScreen[user.ID] = true
	mu.Unlock()

	libs.SuccessMsg(w, userEncodeStr, "system locked")
}

// 系统用户登录（解锁）
func UnLockScreenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		libs.ErrorMsg(w, "invalid request method")
		return
	}

	var req model.SysUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		libs.ErrorMsg(w, "invalid input")
		return
	}

	var user model.SysUser
	if err := model.Db.Where("username = ? and password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		libs.ErrorMsg(w, "invalid username or password")
		return
	}

	// 登录成功，锁屏状态设置为 false
	mu.Lock()
	if userLockedScreen == nil {
		userLockedScreen = make(map[uint]bool)
	}
	userLockedScreen[user.ID] = false
	mu.Unlock()

	libs.SuccessMsg(w, nil, "system unlock")
}

// 检查锁屏状态
func CheckLockedScreenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		libs.ErrorMsg(w, "invalid request method")
		return
	}

	decoderUid := r.Header.Get("sysuser")
	if decoderUid == "" {
		//获取锁屏状态
		libs.SuccessMsg(w, false, "")
		return
	}

	uidStr, err := libs.DecodeFile("godoos", decoderUid)
	if err != nil {
		libs.ErrorMsg(w, "invalid uid")
		return
	}
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		libs.ErrorMsg(w, "parse uid fail")
		return
	}

	mu.Lock()
	defer mu.Unlock()
	status, ok := userLockedScreen[uint(uid)]
	if !ok {
		libs.ErrorMsg(w, "user not found")
		return
	}

	//获取锁屏状态
	libs.SuccessMsg(w, status, "")
}
