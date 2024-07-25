package progress

import (
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func StartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	processesMu.Lock()
	defer processesMu.Unlock()
	cmd, ok := processes[name]
	if ok && !cmd.Running {
		err := cmd.Cmd.Start()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			fmt.Fprintf(w, "failed to start process %s: %v", name, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	err := ExecuteScript(name)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Process %s started.", name)
}
func ExecuteStartAll() error {
	log.Println("Starting all processes...")

	userExeDir := libs.GetRunDir()

	userFis, err := os.ReadDir(userExeDir)
	if err != nil {
		return fmt.Errorf("failed to read user scripts directory: %w", err)
	}

	var names []string
	for _, userFi := range userFis {
		if userFi.IsDir() && !strings.HasPrefix(userFi.Name(), ".") {
			names = append(names, userFi.Name())
		}
	}

	err = batchExecuteScript(names)
	if err != nil {
		return fmt.Errorf("failed to start processes: %w", err)
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
func batchExecuteScript(names []string) error {
	for _, name := range names {
		// 检查processes中是否存在相同的name
		if _, exists := processes[name]; exists {
			log.Printf("Process with name %s already exists", name)
			continue
		}
		if err := ExecuteScript(name); err != nil {
			return errors.Wrapf(err, "failed to script: %v", err)
		}
	}
	return nil
}

// ExecuteScript 执行指定名称的脚本。
// 参数：
// name - 脚本的名称。
// 返回值：
// 返回可能遇到的错误，如果执行成功，则返回nil。
func ExecuteScript(name string) error {
	// 确保name存在于ProcessInfoMap中
	info, ok := ProcessInfoMap[name]
	if !ok {
		return fmt.Errorf("process information for '%s' not found", name)
	}
	// 根据操作系统添加.exe后缀
	binaryExt := ""
	if runtime.GOOS == "windows" {
		binaryExt = ".exe"
	}
	scriptName := name + binaryExt
	exeDir := libs.GetRunDir()                            // 获取执行程序的目录
	scriptPath := filepath.Join(exeDir, name, scriptName) // 拼接脚本的完整路径
	// 检查脚本文件的存在性
	if !PathExists(scriptPath) {
		return fmt.Errorf("script file %s does not exist", scriptPath)
	}

	// 设置并启动脚本执行命令
	var cmd *exec.Cmd
	cmd = exec.Command(scriptPath)
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = setHideConsoleCursor(cmd)
	}
	go func() {
		// 启动脚本命令并返回可能的错误
		if err := cmd.Start(); err != nil {
			log.Printf("failed to start process %s: %v", name, err)
			return
		}
		pingURL := fmt.Sprintf("http://localhost:%s/%s", info.Port, info.PingPath)
		registerProcess(name, pingURL, cmd)
	}()
	return nil
}
func PathExists(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else if os.IsExist(err) {
		return true
	} else {
		return false
	}
}
