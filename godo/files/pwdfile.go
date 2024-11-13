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
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 非加密文件直接返回base64编码
	isHide := IsHaveHiddenFile(basePath, path)
	//Liuziwang888!@#
	if !isHide {
		fileContent, err := ReadFile(basePath, path)
		if err != nil {
			libs.HTTPError(w, http.StatusNotFound, err.Error())
			return
		}
		data := string(fileContent)
		if !strings.HasPrefix(data, "link::") {
			data = base64.StdEncoding.EncodeToString(fileContent)
		}
		resp := libs.APIResponse{Message: "success", Data: data}
		json.NewEncoder(w).Encode(resp)
		return
	}
	// 有隐藏文件说明这是一个加密过的文件,需要验证密码
	fPwd := r.Header.Get("pwd")
	if fPwd == "" {
		libs.HTTPError(w, http.StatusBadRequest, "密码不能为空")
		return
	}
	salt, err := GetSalt(r)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !CheckFilePwd(fPwd, salt) {
		libs.HTTPError(w, http.StatusBadRequest, "密码错误")
		return
	}
	// 密码正确则读取内容并解密
	fileContent, err := ReadFile(basePath, path)
	if err != nil {
		libs.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}
	// 解密
	fileContent, err = libs.DecryptData(fileContent, libs.EncryptionKey)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := string(fileContent)
	if !strings.HasPrefix(data, "link::") {
		data = base64.StdEncoding.EncodeToString(fileContent)
	}
	resp := libs.APIResponse{Message: "success", Data: data}
	json.NewEncoder(w).Encode(resp)
}

// 设置文件密码
func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	fPwd := r.Header.Get("Pwd")
	salt, err := GetSalt(r) // 获取盐值
	// 处理获取盐值时的错误
	if err != nil || fPwd == "" {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 服务端再hash加密
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
	filePwd := r.Header.Get("Pwd")
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
