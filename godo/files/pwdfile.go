package files

import (
	"encoding/base64"
	"encoding/json"
	"godo/libs"
	"net/http"
	"strconv"
	"strings"
)

func HandleReadFile(w http.ResponseWriter, r *http.Request) {

	// 初始值
	path := r.URL.Query().Get("path")
	hasPwd, err := GetPwdFlag()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 检查是否需要密码，如果需要，则从请求头中获取文件密码和盐值
	var fPwd string
	var salt string
	if hasPwd {
		fPwd = r.Header.Get("filePwd") // 获取文件密码
		salt, err = GetSalt(r)         // 获取盐值
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, err.Error()) // 处理获取盐值时的错误
			return
		}
	}

	// 校验文件路径
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 获取文件路径
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 读取内容
	fileContent, err := ReadFile(basePath, path)
	if err != nil {
		libs.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}
	fileData := string(fileContent)
	// 无加密情况
	if !hasPwd {
		// 判断文件开头是否以link:开头
		if !strings.HasPrefix(fileData, "link::") {
			fileData = base64.StdEncoding.EncodeToString(fileContent)
		}
		resp := libs.APIResponse{Message: "success", Data: fileData}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 加密情况
	// 1. 验证文件密码
	if !CheckFilePwd(fPwd, salt) {
		libs.HTTPError(w, http.StatusBadRequest, "密码错误")
		return
	}

	// 2. 解密文件内容
	fileContent, err = libs.DecryptData(fileContent, libs.EncryptionKey)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. base64编码后返回
	// 判断文件开头是否以link:开头
	fileData = string(fileContent)
	if !strings.HasPrefix(fileData, "link::") {
		fileData = base64.StdEncoding.EncodeToString(fileContent)
	}

	// 初始响应
	res := libs.APIResponse{Code: 0, Message: "success", Data: fileData}
	json.NewEncoder(w).Encode(res)
}

// 设置文件密码
func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	fPwd := r.Header.Get("filePwd")
	salt, err := GetSalt(r) // 获取盐值
	// 处理获取盐值时的错误
	if err != nil || fPwd == "" {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 服务端再hash加密
	// "Lqda0ez6DeBhKOHDUklSO1SDJ7QAwHLgUqFYFfN6kU4="
	hashPwd := libs.HashPassword(fPwd, salt)

	// 服务端存储
	err = libs.SetConfigByName("filePwd", hashPwd)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// salt值存储
	err = libs.SetConfigByName("salt", salt)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := libs.APIResponse{Message: "密码设置成功"}
	json.NewEncoder(w).Encode(res)
}

// 更改文件密码
func HandleChangeFilePwd(w http.ResponseWriter, r *http.Request) {
	filePwd := r.Header.Get("filePwd")
	salt, err := GetSalt(r)          // 获取盐值
	if err != nil || filePwd == "" { // 检查错误和密码是否为空
		libs.ErrorMsg(w, "参数错误")
		return
	}
	newPwd := libs.HashPassword(filePwd, salt)
	err = libs.SetConfigByName("filePwd", newPwd)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	libs.SuccessMsg(w, "success", "The file password change success!")
}

// 更改加密状态
func HandleSetIsPwd(w http.ResponseWriter, r *http.Request) {
	isPwd := r.URL.Query().Get("ispwd")
	// 0非加密机器 1加密机器
	isPwdValue, err := strconv.Atoi(isPwd)
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	var isPwdBool bool
	if isPwdValue == 0 {
		isPwdBool = false
	} else {
		isPwdBool = true
	}
	err = libs.SetConfigByName("isPwd", isPwdBool)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	libs.SuccessMsg(w, "success", "")
}
