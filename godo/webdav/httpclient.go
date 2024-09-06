// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package webdav

import (
	"fmt"
	"godo/libs"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var client *Client

func InitWebdav() error {
	clientInfo, ok := libs.GetConfig("webdavClient")
	if !ok {
		return fmt.Errorf("failed to find configuration for webdavClient")
	}

	// 检查 clientInfo 是否为 map 类型
	infoMap, ok := clientInfo.(map[string]interface{})
	if !ok {
		return fmt.Errorf("configuration for webdavClient is not a map")
	}

	url, ok := infoMap["url"].(string)
	if !ok {
		return fmt.Errorf("missing 'url' in webdavClient configuration")
	}

	username, ok := infoMap["username"].(string)
	if !ok {
		return fmt.Errorf("missing 'username' in webdavClient configuration")
	}

	password, ok := infoMap["password"].(string)
	if !ok {
		return fmt.Errorf("missing 'password' in webdavClient configuration")
	}

	client = NewClient(url, username, password)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect to WebDAV server: %v", err)
	}
	return nil
}
func HandlePing(w http.ResponseWriter, r *http.Request) {
	err := client.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleReadDir: 读取目录内容
func HandleReadDir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	files, err := client.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, f := range files {
		w.Write([]byte(f.Name() + "\n"))
	}
}

// HandleStat: 获取文件状态
func HandleStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	fi, err := client.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileInfo, ok := fi.(*File)
	if !ok {
		http.Error(w, "Unexpected file info type", http.StatusInternalServerError)
		return
	}
	// 直接调用 fi.String()，因为 fi 应该是 *File 类型
	w.Write([]byte(fileInfo.String()))
}

// HandleChmod: 改变文件权限 (不支持)
func HandleChmod(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not supported", http.StatusNotImplemented)
}

// HandleExists: 检查文件是否存在
func HandleExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	_, err := client.Stat(path)
	if err != nil {
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}
	w.Write([]byte("File exists"))
}

// HandleReadFile: 读取文件内容
func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	data, err := client.Read(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// HandleUnlink: 删除文件
func HandleUnlink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	err := client.Remove(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File deleted"))
}

// HandleClear: 清空目录 (不支持)
func HandleClear(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not supported", http.StatusNotImplemented)
}

// HandleRename: 重命名文件或目录
func HandleRename(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oldPath := vars["oldPath"]
	newPath := vars["newPath"]
	err := client.Rename(oldPath, newPath, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File renamed"))
}

// HandleMkdir: 创建目录
func HandleMkdir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	err := client.Mkdir(path, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Directory created"))
}

// HandleRmdir: 删除目录
func HandleRmdir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	err := client.RemoveAll(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Directory removed"))
}

// HandleCopyFile: 复制文件
func HandleCopyFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oldPath := vars["oldPath"]
	newPath := vars["newPath"]
	err := client.Copy(oldPath, newPath, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File copied"))
}

// HandleWriteFile: 写入文件
func HandleWriteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = client.Write(path, body, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File written"))
}

// HandleAppendFile: 追加到文件
func HandleAppendFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 读取现有文件内容
	existingContent, err := client.Read(path)
	if err != nil && !IsErrNotFound(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 如果文件不存在，则创建新文件
	if IsErrNotFound(err) {
		err = client.Write(path, body, 0)
	} else {
		// 如果文件存在，则追加内容
		err = client.Write(path, append(existingContent, body...), 0)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File appended"))
}
