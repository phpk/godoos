package files

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// 定义结构体来匹配 JSON 数据结构
type CallbackData struct {
	Key        string `json:"key"`
	Status     int    `json:"status"`
	URL        string `json:"url"`
	ChangesURL string `json:"changesurl"`
	History    struct {
		ServerVersion string `json:"serverVersion"`
		Changes       []struct {
			DocumentSHA256 string `json:"documentSha256"`
			Created        string `json:"created"`
			User           struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		} `json:"changes"`
	} `json:"history"`
	Users   []string `json:"users"`
	Actions []struct {
		Type   int    `json:"type"`
		UserID string `json:"userid"`
	} `json:"actions"`
	LastSave    string `json:"lastsave"`
	NotModified bool   `json:"notmodified"`
	FileType    string `json:"filetype"`
}

// 定义 OnlyOffice 预期的响应结构
type OnlyOfficeResponse struct {
	Error int `json:"error"`
}

// OnlyOffice 回调处理函数
func OnlyOfficeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 打印原始请求体
	//fmt.Printf("Received raw data: %s\n", body)

	// 解析 JSON 数据
	var callbackData CallbackData
	err = json.Unmarshal(body, &callbackData)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// 打印解析后的回调数据
	//fmt.Printf("Received callback: %+v\n", callbackData)

	// 使用 key 查找对应的 path
	mapOnlyOfficeMutex.Lock()
	path, exists := OnlyOfficekeyPathMap[callbackData.Key]
	if exists {
		delete(OnlyOfficekeyPathMap, callbackData.Key) // 删除已使用的 key
		if callbackData.Status == 2 {
			err := downloadFile(path, callbackData.URL)
			log.Printf("Download file error: %v", err)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to download file: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	mapOnlyOfficeMutex.Unlock()
	// 构造响应
	response := OnlyOfficeResponse{
		Error: 0,
	}

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 返回 JSON 响应
	json.NewEncoder(w).Encode(response)
}

// 下载文件
func downloadFile(filePath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	basePath, err := libs.GetOsDir()
	if err != nil {
		return err
	}
	fullFilePath := filepath.Join(basePath, filePath)
	// 使用 os.OpenFile 创建或截断文件
	out, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer out.Close()
	// 获取文件密钥
	fileSecret := libs.GetConfigString("filePwd")
	if fileSecret != "" {
		// 读取文件内容到内存中
		fileBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		encryData, err := libs.EncodeFile(fileSecret, string(fileBytes))
		if err != nil {
			return err
		}
		_, err = out.Write([]byte(encryData))
		return err
	} else {
		_, err = io.Copy(out, resp.Body)
		return err
	}

}
