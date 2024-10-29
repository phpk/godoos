package files

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"godo/libs"
	"net/http"
	"strings"
)

// 带加密读
func HandleReadFile(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")
	fpwd := r.Header.Get("fpwd")
	haspwd := IsHavePwd(fpwd)

	// 获取salt值
	salt := GetSalt(r)

	// 校验文件路径
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 有密码校验密码
	if haspwd {
		if !CheckFilePwd(fpwd, salt) {
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
	fpwd := r.Header.Get("filepwd")
	salt := r.Header.Get("salt")
	// 密码最长16位
	if fpwd == "" || len(fpwd) > 16 {
		libs.ErrorMsg(w, "密码长度为空或者过长,最长为16位")
		return
	}
	// md5加密
	mhash := md5.New()
	mhash.Write([]byte(fpwd))
	v := mhash.Sum(nil)
	pwdstr := hex.EncodeToString(v)

	// 服务端再hash加密
	hashpwd := libs.HashPassword(pwdstr, salt)

	// 服务端存储
	req := libs.ReqBody{
		Name:  "filepwd",
		Value: hashpwd,
	}
	libs.SetConfig(req)

	// salt值存储
	reqSalt := libs.ReqBody{
		Name:  "salt",
		Value: salt,
	}
	libs.SetConfig(reqSalt)
	res := libs.APIResponse{Message: "success", Data: pwdstr}
	json.NewEncoder(w).Encode(res)
}
