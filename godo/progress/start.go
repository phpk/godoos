package progress

import (
	"fmt"
	"net/http"

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
	cmd, exists := processes[name]
	if !exists {
		return fmt.Errorf("process information for '%s' not found", name)
	}
	if cmd.Running {
		return fmt.Errorf("process %s is already running", name)
	}
	if err := cmd.Cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process %s: %v", name, err)
	}
	cmd.Running = true
	return nil
	// 确保name存在于ProcessInfoMap中
	// _, ok := ProcessInfoMap[name]
	// if !ok {
	// 	return fmt.Errorf("process information for '%s' not found", name)
	// }
	// // 根据操作系统添加.exe后缀
	// binaryExt := ""
	// if runtime.GOOS == "windows" {
	// 	binaryExt = ".exe"
	// }
	// scriptName := name + binaryExt
	// exeDir := libs.GetRunDir()                            // 获取执行程序的目录
	// scriptPath := filepath.Join(exeDir, name, scriptName) // 拼接脚本的完整路径
	// // 检查脚本文件的存在性
	// if !libs.PathExists(scriptPath) {
	// 	return fmt.Errorf("script file %s does not exist", scriptPath)
	// }

	// // 设置并启动脚本执行命令
	// var cmd *exec.Cmd
	// cmd = exec.Command(scriptPath)
	// if runtime.GOOS == "windows" {
	// 	// 在Windows上，通过设置CreationFlags来隐藏窗口
	// 	cmd = SetHideConsoleCursor(cmd)
	// }
	// go func() {
	// 	// 启动脚本命令并返回可能的错误
	// 	if err := cmd.Start(); err != nil {
	// 		log.Printf("failed to start process %s: %v", name, err)
	// 		return
	// 	}
	// 	RegisterProcess(name, cmd)
	// }()
	// return nil
}
