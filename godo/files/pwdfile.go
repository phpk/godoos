package files

import (
	"encoding/base64"
	"encoding/json"
	"godo/libs"
	"net/http"
	"strings"
)

// 带加密读
func HandleReadFile(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")
	fPwd := r.Header.Get("fPwd")
	hasPwd := IsHavePwd(fPwd)

	// 获取salt值
	salt := GetSalt(r)

	// 校验文件路径
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 有密码校验密码
	if hasPwd {
		if !CheckFilePwd(fPwd, salt) {
			libs.HTTPError(w, http.StatusBadRequest, "密码错误")
			return
		}
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

	// 解密
	data, err := libs.DecryptData(fileContent, libs.EncryptionKey)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	content := string(data)
	// 检查文件内容是否以"link::"开头
	if !strings.HasPrefix(content, "link::") {
		content = base64.StdEncoding.EncodeToString(data)
	}

	// 初始响应
	res := libs.APIResponse{Code: 0, Message: "success", Data: content}

	json.NewEncoder(w).Encode(res)
}

// 设置文件密码
func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	fPwd := r.Header.Get("filepPwd")
	salt := r.Header.Get("salt")
	// 服务端再hash加密
	hashPwd := libs.HashPassword(fPwd, salt)

	// 服务端存储
	req := libs.ReqBody{
		Name:  "filePwd",
		Value: hashPwd,
	}
	libs.SetConfig(req)

	// salt值存储
	reqSalt := libs.ReqBody{
		Name:  "salt",
		Value: salt,
	}
	libs.SetConfig(reqSalt)
	res := libs.APIResponse{Message: "密码设置成功"}
	json.NewEncoder(w).Encode(res)
}

// 更改文件密码
func HandleChangeFilePwd(w http.ResponseWriter, r *http.Request) {
	filePwd := r.Header.Get("filePwd")
	salt := r.Header.Get("salt")
	if filePwd == "" || salt == "" {
		libs.ErrorMsg(w, "密码为空")
		return
	}
	newPwd := libs.HashPassword(filePwd, salt)
	pwdReq := libs.ReqBody{
		Name:  "filePwd",
		Value: newPwd,
	}
	libs.SetConfig(pwdReq)
	libs.SuccessMsg(w, "success", "The file password change success!")
}
