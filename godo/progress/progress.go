package progress

import (
	"fmt"
	"godo/libs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Name     string
	Running  bool
	ExitCode int
	Pid      int
	Cmd      *exec.Cmd
}

var (
	processesMu sync.RWMutex
	processes   = make(map[string]*Process)
)

func RegisterProcess(name string, cmdstr *exec.Cmd) {
	processesMu.Lock()
	defer processesMu.Unlock()
	processes[name] = &Process{
		Name:    name,
		Running: true,
		Pid:     cmdstr.Process.Pid,
		Cmd:     cmdstr,
	}
}
func GetCmd(name string) *Process {
	// processesMu.Lock()
	// defer processesMu.Unlock()
	info, ok := processes[name]
	if !ok {
		return &Process{
			Name:    name,
			Running: false,
		}
	}
	return info
}

// StartCmd 执行指定名称的脚本。
// 参数：
// name - 脚本的名称。
// 返回值：
// 返回可能遇到的错误，如果执行成功，则返回nil。
func StartCmd(name string) error {
	info, ok := processes[name]
	if ok && info.Running {
		return fmt.Errorf("process information for '%s' is runing", name)
	}
	//appName := name
	scriptPath, err := libs.GetCmdPath(name)
	if err != nil {
		return fmt.Errorf("failed to extract zip file: %v", err)
	}

	// 设置并启动脚本执行命令
	var cmd *exec.Cmd
	switch name {
	case "frpc":
		// 检查配置文件
		configPath := filepath.Join(filepath.Dir(scriptPath), "frpc.ini")
		log.Printf("Config file not found at %s, creating new file", configPath)
		if !libs.PathExists(configPath) {
			return fmt.Errorf("frpc config file not found")
		}

		params := []string{
			"-c",
			configPath,
		}
		cmd = exec.Command(scriptPath, params...)
	default:
		cmd = exec.Command(scriptPath)
	}
	log.Printf("Starting %s", scriptPath)
	//cmd = exec.Command(scriptPath)
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	go func() {
		// 启动脚本命令并返回可能的错误
		if err := cmd.Start(); err != nil {
			//log.Printf("failed to start process %s: %v", name, err)
			return
		}
		RegisterProcess(name, cmd)
		// 等待命令完成
		if err := cmd.Wait(); err != nil {
			log.Printf("command failed for %s: %v", name, err)
			//return
		} else {
			log.Printf("%s command completed successfully", name)
		}

		// 命令完成后，更新进程信息
		processesMu.Lock()
		defer processesMu.Unlock()
		if p, ok := processes[name]; ok {
			p.Running = false
			p.ExitCode = cmd.ProcessState.ExitCode()
		}
	}()
	return nil
}
func StopCmd(name string) error {
	cmd, ok := processes[name]
	if !ok {
		return fmt.Errorf("process information for '%s' not found", name)
	}
	if !cmd.Running {
		return nil
	}

	// 停止进程并更新status
	// TODO: 如果有多个pid 的情况，需要处理
	if err := KillByPid(cmd.Pid); err != nil {
		return fmt.Errorf("failed to kill process %s: %v", name, err)
	}

	// if err := cmd.Cmd.Process.Kill(); err != nil {
	// 	return fmt.Errorf("failed to kill process %s: %v", name, err)
	// }
	//delete(processes, name) // 更新status，表示进程已停止
	cmd.Running = false
	return nil
}

func RestartCmd(name string) error {
	if err := StopCmd(name); err != nil {
		log.Printf("stopping the app encountered an error: %s", err)
		return err
	}
	return StartCmd(name)
}
func StopAllCmd() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := cmd.Cmd.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to stop process %s: %v", name, err)
		}
	}
	return nil
}

func KillByPid(pid int) error {
	log.Println("Killing process with PID:", pid)
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("failed to create process object: %w", err)
	}

	// 杀死进程及其子进程
	if err := p.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}
	return nil
}

// // findPidsWindows 在 Windows 系统下查找具有指定名称的进程的 PID
// func findPidsWindows(name string) ([]int, error) {
// 	cmd := exec.Command("tasklist")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute tasklist: %w", err)
// 	}
// 	lines := strings.Split(string(output), "\r\n")
// 	var pids []int
// 	for _, line := range lines {
// 		fields := strings.Fields(line)
// 		if len(fields) >= 3 {
// 			if strings.Contains(strings.ToLower(fields[0]), strings.ToLower(name)) {
// 				pid, err := strconv.Atoi(fields[1])
// 				if err != nil {
// 					return nil, fmt.Errorf("failed to convert PID to integer: %w", err)
// 				}
// 				pids = append(pids, pid)
// 			}
// 		}
// 	}
// 	return pids, nil
// }
