package files

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleWriteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fullFilePath := filepath.Join(basePath, path)
	newFile, err := os.Create(fullFilePath)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer newFile.Close()
	//  获取文件内容
	fileContent, _, err := r.FormFile("content")
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer fileContent.Close()
	fileData, err := io.ReadAll(fileContent)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	configPwd, ishas := libs.GetConfig("filePwd")
	//log.Printf("configPwd: %s", configPwd)
	// 如果不是加密文件或者exe文件
	if !ishas || strings.HasPrefix(string(fileData), "link::") {
		// 没开启加密，直接明文写入
		_, err := newFile.Write(fileData)
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, "数据写入失败")
			return
		}
		err = CheckAddDesktop(path)
		if err != nil {
			log.Printf("Error adding file to desktop: %s", err.Error())
		}
		libs.SuccessMsg(w, "", "文件写入成功")
		return
	} else {
		// 开启加密后，写入加密数据
		configPwdStr, ok := configPwd.(string)
		if !ok {
			libs.HTTPError(w, http.StatusInternalServerError, "配置文件密码格式错误")
			return
		}
		_, err = newFile.WriteString(fmt.Sprintf("@%s@", configPwdStr))
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, "密码写入失败")
			return
		}
		entryData, err := libs.EncryptData(fileData, []byte(configPwdStr))
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, "文件加密失败")
			return
		}
		_, err = newFile.Write(entryData)
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, "文件写入失败")
			return
		}
		err = CheckAddDesktop(path)
		if err != nil {
			log.Printf("Error adding file to desktop: %s", err.Error())
		}
		libs.SuccessMsg(w, "", "文件写入成功")
	}
}

func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		libs.ErrorMsg(w, "文件路径不能为空")
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	fullFilePath := filepath.Join(basePath, path)
	file, err := os.Open(fullFilePath)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	strData := string(fileData)
	//log.Printf("isPwd: %v", strData)
	if strings.HasPrefix(strData, "link::") {
		res := libs.APIResponse{Message: "文件读取成功", Data: strData}
		json.NewEncoder(w).Encode(res)
		return
	}
	// 判断是否为加密文件
	isPwd := IsPwdFile(fileData)

	if !isPwd {
		// 未加密文件，直接返回
		content := base64.StdEncoding.EncodeToString(fileData)
		res := libs.APIResponse{Message: "文件读取成功", Data: content}
		json.NewEncoder(w).Encode(res)
		return
	}
	Pwd := r.Header.Get("Pwd")
	filePwd := strData[1:33]
	// Pwd为空，info密码与文件密码做比对
	if Pwd == "" {
		configPwd, ishas := libs.GetConfig("filePwd")
		if !ishas {
			libs.ErrorMsg(w, "未设置密码")
			return
		}
		configPwdStr, ok := configPwd.(string)
		if !ok {
			libs.ErrorMsg(w, "后端配置文件密码格式错误")
			return
		}
		// 校验密码
		if filePwd != configPwdStr {
			res := libs.APIResponse{Message: "密码错误,请输入正确的密码", Code: -1, Error: "needPwd"}
			json.NewEncoder(w).Encode(res)
			return
		}
		decryptData, err := libs.DecryptData(fileData[34:], []byte(filePwd))
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
		content := base64.StdEncoding.EncodeToString(decryptData)
		res := libs.APIResponse{Message: "加密文件读取成功", Data: content}
		json.NewEncoder(w).Encode(res)
		return
	} else {
		// Pwd不为空，Pwd与文件密码做比对
		if Pwd != filePwd {
			res := libs.APIResponse{Message: "密码错误,请输入正确的密码", Code: -1, Error: "needPwd"}
			json.NewEncoder(w).Encode(res)
			return
		}
		decryptData, err := libs.DecryptData(fileData[34:], []byte(filePwd))
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
		content := base64.StdEncoding.EncodeToString(decryptData)
		res := libs.APIResponse{Message: "加密文件读取成功", Data: content}
		json.NewEncoder(w).Encode(res)
	}
}

func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	Pwd := r.Header.Get("Pwd")
	if Pwd == "" {
		err := libs.DeleteConfig("filePwd")
		if err != nil {
			libs.ErrorMsg(w, fmt.Sprintf("取消加密失败:%s", err.Error()))
			return
		}
		libs.SuccessMsg(w, nil, "取消加密成功")
		return
	}
	err := libs.SetConfigByName("filePwd", Pwd)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "设置密码成功")
}
