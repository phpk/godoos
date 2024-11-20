/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package store

import (
	"godo/files"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// UploadHandler 处理上传的HTTP请求
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// 解析上传的文件
	err := r.ParseMultipartForm(10000 << 20) // 限制最大上传大小为1000MB
	if err != nil {
		http.Error(w, "上传文件过大"+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("files") // 表单字段名为"files"
	if err != nil {
		http.Error(w, "没有找到文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("读取文件内容出错: %v", err)
		http.Error(w, "读取文件内容出错", http.StatusInternalServerError)
		return
	}
	cachePath := libs.GetCacheDir()
	baseName := filepath.Base(header.Filename)

	savePath := filepath.Join(cachePath, baseName)

	out, err := os.Create(savePath)
	if err != nil {
		log.Printf("创建文件出错: %v", err)
		http.Error(w, "保存文件出错", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将文件内容写入到服务器上的文件
	_, err = out.Write(fileBytes)
	if err != nil {
		log.Printf("写入文件出错: %v", err)
		http.Error(w, "写入文件出错", http.StatusInternalServerError)
		return
	}

	zipPath, err := files.Decompress(savePath, cachePath)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	storeFile := filepath.Join(zipPath, "install.json")
	if !libs.PathExists(storeFile) {
		libs.ErrorMsg(w, "store.json not found")
		return
	}
	installCacheInfo, err := GetInstallInfoByPath(storeFile)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	runDir := libs.GetRunDir()
	targetDir := filepath.Join(runDir, installCacheInfo.Name)
	err = files.CopyResource(zipPath, targetDir)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	err = os.RemoveAll(zipPath)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	installInfo, err := Installation(installCacheInfo.Name)
	if err != nil {
		libs.ErrorData(w, installInfo, "the install.json is error:"+err.Error())
		return
	}
	libs.SuccessMsg(w, installInfo, "install the app success!")
}
