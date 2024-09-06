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
package store

import (
	"fmt"
	"godo/libs"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
)

func StartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	processesMu.Lock()
	defer processesMu.Unlock()
	err := ExecuteScript(name)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Process %s started.", name)
}
func ExecuteStartAll() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := cmd.Cmd.Start(); err != nil {
			return fmt.Errorf("failed to stop process %s: %v", name, err)
		}
	}

	return nil
}
func StartAll(w http.ResponseWriter, r *http.Request) {
	if err := ExecuteStartAll(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All processes started.")
}

// ExecuteScript 执行指定名称的脚本。
// 参数：
// name - 脚本的名称。
// 返回值：
// 返回可能遇到的错误，如果执行成功，则返回nil。
func ExecuteScript(name string) error {
	storeInfo, err := GetStoreInfo(name)
	if err != nil {
		return err
	}
	err = runStart(storeInfo)
	if err != nil {
		return fmt.Errorf("failed to run script: %v", err)
	}
	return nil
}
func RunOutHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	fileName := filepath.Base(url)
	cacheDir := libs.GetCacheDir()
	filePath := filepath.Join(cacheDir, fileName)
	if !libs.PathExists(filePath) {
		libs.ErrorMsg(w, "file not found")
		return
	}
	cmd := exec.Command(filePath)
	if err := cmd.Start(); err != nil {
		libs.ErrorMsg(w, fmt.Sprintf("start error: %v", err))
		return
	}

	// 如果 cmd.Start() 成功，返回成功消息
	libs.SuccessMsg(w, "", "start success")

}
