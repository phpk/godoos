package store

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
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
func GetStoreInfo(name string) (StoreInfo, error) {
	var storeInfo StoreInfo
	exeDir := libs.GetRunDir()
	infoPath := filepath.Join(exeDir, name, "info.json")
	if !libs.PathExists(infoPath) {
		return storeInfo, fmt.Errorf("process information for '%s' not found", name)
	}

	content, err := os.ReadFile(infoPath)
	if err != nil {
		return storeInfo, fmt.Errorf("failed to read info.json: %v", err)
	}
	if err := json.Unmarshal(content, &storeInfo); err != nil {
		return storeInfo, fmt.Errorf("failed to unmarshal info.json: %v", err)
	}
	scriptPath := storeInfo.Setting.BinPath
	if !libs.PathExists(scriptPath) {
		return storeInfo, fmt.Errorf("script file '%s' not found", scriptPath)
	}
	return storeInfo, nil
}

// ExecuteScript 执行指定名称的脚本。
// 参数：
// name - 脚本的名称。
// 返回值：
// 返回可能遇到的错误，如果执行成功，则返回nil。
func ExecuteScript(name string) error {
	cmd, exists := processes[name]
	if !exists {
		log.Printf("process information for '%s' not found", name)
		//return fmt.Errorf("process information for '%s' not found", name)
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
	if cmd.Running {
		return nil
	}

	// if err := cmd.Cmd.Start(); err != nil {
	// 	log.Printf("failed to start process %s: %v", name, err)
	// 	return nil
	// }
	processes[name].Running = true
	go func() {

		if err := cmd.Cmd.Start(); err != nil {
			log.Printf("failed to start process %s: %v", name, err)
			return
		}
		// processesMu.Lock()
		// defer processesMu.Unlock()

		processes[name].Cmd = cmd.Cmd
	}()

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
