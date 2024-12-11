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

func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		libs.ErrorMsg(w, "文件路径不能为空")
		return
	}
	stream := r.URL.Query().Get("stream")
	isStream := false
	if stream != "" {
		isStream = true
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
	text := string(fileData)
	//log.Printf("isPwd: %v", strData)
	if strings.HasPrefix(text, "link::") {
		libs.SuccessMsg(w, text, "文件读取成功")
		return
	}
	// 判断是否为加密文件
	isPwd := libs.IsEncryptedFile(text)

	if !isPwd {
		// 未加密文件，直接返回

		if isStream {
			w.Header().Set("Content-Type", "application/octet-stream")
			_, err := w.Write(fileData)
			if err != nil {
				libs.ErrorMsg(w, err.Error())
				return
			}
		} else {
			if len(fileData) > 0 {
				text = base64.StdEncoding.EncodeToString(fileData)
			}
			//text = base64.StdEncoding.EncodeToString(fileData)
			libs.SuccessMsg(w, text, "文件读取成功")
		}

		return
	}
	filePwd := r.Header.Get("pwd")
	// 获取文件密钥
	fileSecret := libs.GetConfigString("filePwd")
	if filePwd == "" && fileSecret == "" {
		libs.Error(w, "加密文件，需要密码", "needPwd")
		return
	}
	decodeFile := ""
	needPwd := false
	if fileSecret != "" {
		decodeFile, err = libs.DecodeFile(fileSecret, text)
		if err != nil {
			libs.Error(w, "请输入文件密码", "needPwd")
			return
		}
		needPwd = true
	}
	if filePwd != "" && !needPwd {
		decodeFile, err = libs.DecodeFile(filePwd, text)
		if err != nil {
			libs.Error(w, "文件密码输入错误", "needPwd")
			return
		}
	}
	if isStream {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err := w.Write([]byte(decodeFile))
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	} else {
		if len(decodeFile) > 0 {
			decodeFile = base64.StdEncoding.EncodeToString([]byte(decodeFile))
		}
		libs.SuccessMsg(w, decodeFile, "加密文件读取成功")
	}

}
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
	//  获取文件内容
	fileContent, _, err := r.FormFile("content")
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer fileContent.Close()
	content, err := io.ReadAll(fileContent)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		libs.ErrorMsg(w, "Failed to create file.")
		return
	}
	defer file.Close()

	oldContent, err := io.ReadAll(file)
	text := string(oldContent)
	//log.Printf("text: %s", text)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	isPwdFile := libs.IsEncryptedFile(text)

	filePwd := r.Header.Get("pwd")
	// 获取文件密钥
	fileSecret := libs.GetConfigString("filePwd")
	haslink := strings.HasPrefix(string(content), "link::")
	//log.Printf("haslink:%v,fileSecret: %s,isPwdFile:%v,filePwd:%s", haslink, fileSecret, isPwdFile, filePwd)
	needPwd := false

	if fileSecret != "" || filePwd != "" {
		needPwd = true
	}
	if isPwdFile {
		needPwd = true
	}
	if haslink {
		needPwd = false
	}
	log.Printf("needPwd:%v", needPwd)
	// 即不是加密用户又不是加密文件
	if !needPwd {
		// 直接写入新内容
		file.Truncate(0)
		file.Seek(0, 0)
		//log.Printf("write file content%v", content)
		_, err = file.Write(content)
		if err != nil {
			libs.ErrorMsg(w, "Failed to write file content.")
			return
		}
		CheckAddDesktop(path)
		libs.SuccessMsg(w, "", "文件写入成功")
		return
	}

	//log.Printf("fileSecret:%s,filePwd:%s,ispwdfile:%v", fileSecret, filePwd, isPwdFile)

	// 是加密文件,写入需继续加密
	pwd := ""
	if isPwdFile {
		//先尝试系统解密
		if fileSecret != "" {
			_, err := libs.DecodeFile(fileSecret, text)
			if err == nil {
				pwd = fileSecret
			}
		}
		//先用户输入解密
		if pwd == "" && filePwd != "" {
			_, err := libs.DecodeFile(filePwd, text)
			if err == nil {
				pwd = filePwd
			}
		}

	} else {
		// 不是加密文件，先判断是否有用户输入
		if filePwd != "" {
			pwd = filePwd
		} else {
			// 没有用户输入，则使用系统默认密码
			if fileSecret != "" {
				pwd = fileSecret
			}
		}
	}
	if pwd == "" {
		libs.Error(w, "文件密码错误", "needPwd")
		//log.Printf("pwd is empty, returning error")
		return
	}
	//log.Printf("pwd: %s", pwd)

	encryData, err := libs.EncodeFile(pwd, string(content))
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	// 清空文件内容
	file.Truncate(0)
	file.Seek(0, 0)

	_, err = file.Write([]byte(encryData))
	if err != nil {
		libs.ErrorMsg(w, fmt.Sprintf("文件内容写入失败: %s", err.Error()))
		return
	}
	CheckAddDesktop(path)
	libs.SuccessMsg(w, "", "文件写入成功")
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

// 定义结构体来匹配 JSON 数据结构
type CallbackData struct {
	Key     string   `json:"key"`
	Status  int      `json:"status"`
	Users   []string `json:"users"`
	Actions []struct {
		Type   int    `json:"type"`
		UserID string `json:"userid"`
	} `json:"actions"`
}

// 定义 OnlyOffice 预期的响应结构
type OnlyOfficeResponse struct {
	Error int `json:"error"`
}

func OnlyOfficeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 打印原始请求体
	fmt.Printf("Received raw data: %s\n", body)

	// 解析 JSON 数据
	// var callbackData CallbackData
	// err = json.Unmarshal(body, &callbackData)
	// if err != nil {
	// 	http.Error(w, "Failed to decode request body", http.StatusBadRequest)
	// 	return
	// }

	// // 打印解析后的回调数据
	// fmt.Printf("Received callback: %+v\n", callbackData)

	// 构造响应
	response := OnlyOfficeResponse{
		Error: 0,
	}

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 返回 JSON 响应
	json.NewEncoder(w).Encode(response)

}
