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
package files

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func HandleSystemInfo(w http.ResponseWriter, r *http.Request) {
	info, err := libs.GetSystemInfo()
	if err != nil {
		libs.ErrorMsg(w, "Failed to retrieve file information.")
		return
	}
	libs.SuccessMsg(w, info, "File information retrieved successfully.")
}

// HandleReadDir handles reading a directory
func HandleReadDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	files, err := ReadDir(basePath, path)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var osFileInfos []OsFileInfo
	for _, entry := range files {
		// 过滤掉以.开头的文件（隐藏文件）
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		osFileInfo, err := GetFileInfo(entry, basePath, path)
		if err != nil {
			libs.HTTPError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// 如果是文件，读取内容
		if osFileInfo.IsFile {
			file, err := os.Open(filepath.Join(basePath, osFileInfo.Path))
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open file: %v", err))
				return
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read file content: %v", err))
				return
			}
			osFileInfo.Content = string(content)

			// 检查文件内容是否以"link::"开头
			if strings.HasPrefix(osFileInfo.Content, "link::") {
				osFileInfo.IsSymlink = true
			} else {
				osFileInfo.Content = ""
			}
		}

		osFileInfos = append(osFileInfos, *osFileInfo)
	}
	// 按照 ModTime 进行降序排序
	sort.Slice(osFileInfos, func(i, j int) bool {
		return osFileInfos[i].ModTime.Before(osFileInfos[j].ModTime)
	})
	res := libs.APIResponse{
		Message: "Directory read successfully.",
		Data:    osFileInfos,
	}
	json.NewEncoder(w).Encode(res)
}

// HandleStat retrieves file information
func HandleStat(w http.ResponseWriter, r *http.Request) {
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

	// 调用GetFileInfo，传入路径而不是fs.DirEntry
	osFileInfo, err := GetFileInfo(path, basePath, "")
	if err != nil {
		libs.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	res := libs.APIResponse{
		Message: "File information retrieved successfully.",
		Data:    osFileInfo,
	}
	json.NewEncoder(w).Encode(res)
}

// HandleExists checks if a file or directory exists
func HandleExists(w http.ResponseWriter, r *http.Request) {
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
	exists := Exists(basePath, path)
	message := "File does not exist."
	if exists {
		message = "File exists."
	}
	res := libs.APIResponse{
		Message: message,
		Data:    exists,
	}
	json.NewEncoder(w).Encode(res)
}

// HandleReadFile reads a file's content
func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fpwd := r.Header.Get("fpwd")
	haspwd := IsHavePwd(fpwd)
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
	content := string(fileContent)
	// 检查文件内容是否以"link::"开头
	if !strings.HasPrefix(content, "link::") {
		content = base64.StdEncoding.EncodeToString(fileContent)
	}

	// 初始响应
	res := libs.APIResponse{Code: 0, Message: "success"}
	switch haspwd {
	case true:
		// 有密码检验密码
		isreal := CheckFilePwd(fpwd)
		// 密码正确返回原文，否则返回加密文本
		if isreal {
			res.Data = content
		} else {
			data, err := libs.EncryptData(fileContent, libs.EncryptionKey)
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, err.Error())
				return
			}
			res.Data = base64.StdEncoding.EncodeToString(data)
		}
	case false:
		res.Data = content
	}

	json.NewEncoder(w).Encode(res)
}

// HandleUnlink removes a file
func HandleUnlink(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := CheckDeleteDesktop(path)
	if err != nil {
		log.Printf("Error deleting file from desktop: %s", err.Error())
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = Unlink(basePath, path)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	res := libs.APIResponse{Message: fmt.Sprintf("File '%s' successfully removed.", path)}
	json.NewEncoder(w).Encode(res)
}

// HandleClear removes the entire filesystem (Caution: Use with care!)
func HandleClear(w http.ResponseWriter, r *http.Request) {
	// basePath, err := libs.GetOsDir()
	// if err != nil {
	// 	libs.HTTPError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// err = Clear(basePath)
	// if err != nil {
	// 	libs.HTTPError(w, http.StatusConflict, err.Error())
	// 	return
	// }
	err := RecoverOsSystem()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := libs.APIResponse{Message: "FileSystem successfully cleared."}
	json.NewEncoder(w).Encode(res)
}

// HandleRename renames a file
func HandleRename(w http.ResponseWriter, r *http.Request) {
	oldPath := r.URL.Query().Get("oldPath")
	newPath := r.URL.Query().Get("newPath")
	if err := validateFilePathPair(oldPath, newPath); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = CheckDeleteDesktop(oldPath)
	if err != nil {
		log.Printf("Error deleting file from desktop: %s", err.Error())
	}
	err = Rename(basePath, oldPath, newPath)
	if err != nil {
		log.Printf("Error renaming file: %s", err.Error())
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}

	err = CheckAddDesktop(newPath)
	if err != nil {
		log.Printf("Error adding file to desktop: %s", err.Error())
	}
	res := libs.APIResponse{Message: fmt.Sprintf("File '%s' successfully renamed to '%s'.", oldPath, newPath)}
	json.NewEncoder(w).Encode(res)
}

// HandleMkdir creates a directory
func HandleMkdir(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("dirPath")
	if err := validateFilePath(dirPath); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = Mkdir(basePath, dirPath)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	err = CheckAddDesktop(dirPath)
	if err != nil {
		log.Printf("Error adding file to desktop: %s", err.Error())
	}
	res := libs.APIResponse{Message: fmt.Sprintf("Directory '%s' created successfully.", dirPath)}
	json.NewEncoder(w).Encode(res)
}

// HandleRmdir removes a directory
func HandleRmdir(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("dirPath")
	if err := validateFilePath(dirPath); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = CheckDeleteDesktop(dirPath)
	if err != nil {
		log.Printf("Error deleting file from desktop: %s", err.Error())
	}
	err = Rmdir(basePath, dirPath)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	res := libs.APIResponse{Message: fmt.Sprintf("Directory '%s' successfully removed.", dirPath)}
	json.NewEncoder(w).Encode(res)
}

// HandleCopyFile copies a file
func HandleCopyFile(w http.ResponseWriter, r *http.Request) {
	srcPath := r.URL.Query().Get("srcPath")
	dstPath := r.URL.Query().Get("dstPath")
	if err := validateFilePathPair(srcPath, dstPath); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = CopyFile(filepath.Join(basePath, srcPath), filepath.Join(basePath, dstPath))
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	err = CheckAddDesktop(dstPath)
	if err != nil {
		log.Printf("Error adding file to desktop: %s", err.Error())
	}
	res := libs.APIResponse{Message: fmt.Sprintf("File '%s' successfully copied to '%s'.", srcPath, dstPath)}
	json.NewEncoder(w).Encode(res)
}

// HandleWriteFile writes content to a file
func HandleWriteFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("filePath")
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取文件内容
	fileContent, _, err := r.FormFile("content")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer fileContent.Close()
	// 输出到控制台进行调试
	//fmt.Printf("Body content: %v\n", fileContent)
	file, err := os.Create(filepath.Join(basePath, filePath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, fileContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	err = CheckAddDesktop(filePath)
	if err != nil {
		log.Printf("Error adding file to desktop: %s", err.Error())
	}
	res := libs.APIResponse{Message: fmt.Sprintf("File '%s' successfully written.", filePath)}
	json.NewEncoder(w).Encode(res)
}

// HandleAppendFile appends content to a file
func HandleAppendFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("filePath")
	if err := validateFilePath(filePath); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	content, _, err := r.FormFile("content")
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = AppendToFile(filepath.Join(basePath, filePath), content)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	res := libs.APIResponse{Message: fmt.Sprintf("Content appended to file '%s'.", filePath)}
	json.NewEncoder(w).Encode(res)
}

// HandleChmod changes the permissions of a file
func HandleChmod(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
	}

	// 解析请求体
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// 使用解码后的数据
	path := reqData.Path
	modeStr := reqData.Mode

	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	mode, err := parseMode(modeStr)
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, fmt.Sprintf("Invalid mode: %s", err.Error()))
		return
	}

	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Chmod(basePath, path, mode)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}

	res := libs.APIResponse{Message: fmt.Sprintf("Permissions of file '%s' successfully changed to %o.", path, mode)}
	json.NewEncoder(w).Encode(res)
}
func HandleZip(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	ext := r.URL.Query().Get("ext")
	if ext == "" {
		ext = "zip"
	}
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sourcePath := filepath.Join(basePath, path)
	compressedFilePath := filepath.Join(basePath, path+"."+ext)
	err = CompressFileOrFolder(sourcePath, compressedFilePath)
	if err != nil {
		libs.HTTPError(w, http.StatusConflict, err.Error())
		return
	}
	res := libs.APIResponse{Message: fmt.Sprintf("success zip %s.%s.", path, ext)}
	json.NewEncoder(w).Encode(res)
}

func HandleUnZip(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	if err := validateFilePath(path); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	zipFilePath := filepath.Join(basePath, path)
	destPath := filepath.Dir(zipFilePath)
	zipPath, err := Decompress(zipFilePath, destPath)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	res := libs.APIResponse{Message: fmt.Sprintf("success unzip %s.", path), Data: zipPath}
	json.NewEncoder(w).Encode(res)
}

// parseMode converts a string representation of a mode to an os.FileMode
func parseMode(modeStr string) (os.FileMode, error) {
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}
func HandleDesktop(w http.ResponseWriter, r *http.Request) {
	rootInfo, err := GetDesktop()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	libs.SuccessMsg(w, rootInfo, "success")
}

// 设置文件密码
func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	fpwd := r.Header.Get("filepwd")
	// 密码最长16位
	if fpwd == "" || len(fpwd) > 16 {
		libs.ErrorMsg(w, "密码长度为空或者过长,最长为16位")
		return
	}
	// 服务端存储
	req := libs.ReqBody{
		Name:  "filepwd",
		Value: fpwd,
	}
	libs.SetConfig(req)
	// 客户端加密
	mhash := md5.New()
	mhash.Write([]byte(fpwd))
	v := mhash.Sum(nil)
	pwdstr := hex.EncodeToString(v)
	res := libs.APIResponse{Message: "success", Data: pwdstr}
	json.NewEncoder(w).Encode(res)
}
