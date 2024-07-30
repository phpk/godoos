package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type PortRangeResponse struct {
	Start        int   `json:"start"`
	End          int   `json:"end"`
	EnabledPorts []int `json:"enabled_ports"`
}

func getProcessIdsOnPort(port int) ([]string, error) {
	osType := runtime.GOOS

	var cmd *exec.Cmd
	var output []byte
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-ti", fmt.Sprintf("tcp:%d", port))
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Get-Process | Where-Object {$_.Id -eq "+strconv.Itoa(port)+"} | Select-Object -ExpandProperty Id")
	default:
		return nil, fmt.Errorf("unsupported operating system")
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// 如果lsof或powershell命令找不到任何进程，它会返回非零退出代码，这是正常情况
			if exitError.ExitCode() != 1 {
				return nil, fmt.Errorf("failed to list processes on port %d: %v", port, err)
			}
		} else {
			return nil, fmt.Errorf("failed to list processes on port %d: %v", port, err)
		}
	}

	pids := strings.Fields(strings.TrimSpace(string(output)))
	return pids, nil
}
func listEnabledPorts(portRangeStart, portRangeEnd int) ([]int, error) {
	var usedPorts []int
	var wg sync.WaitGroup

	for i := portRangeStart; i <= portRangeEnd; i++ {
		currentPort := i // 创建一个新的变量来绑定当前的i值
		wg.Add(1)
		go func() { // 注意这里不再直接传入port，而是使用currentPort
			defer wg.Done()

			pids, err := getProcessIdsOnPort(currentPort)
			if err != nil {
				log.Printf("Error checking port %d: %v", currentPort, err)
			}

			if len(pids) > 0 {
				usedPorts = append(usedPorts, currentPort)
			}
		}()
	}

	wg.Wait()

	return usedPorts, nil
}

func killProcess(pid int) error {
	osType := runtime.GOOS

	var cmd *exec.Cmd
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid)) // /F 表示强制结束
	default:
		return fmt.Errorf("unsupported operating system")
	}

	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to kill process with PID %d: %v", pid, err)
	}

	return err
}

func killAllProcessesOnPort(port int, w http.ResponseWriter) {
	pids, err := getProcessIdsOnPort(port)
	if err != nil {
		http.Error(w, "Failed to list processes", http.StatusInternalServerError)
		return
	}

	for _, pidStr := range pids {
		if pidStr == "" {
			continue
		}

		pidInt, err := strconv.Atoi(pidStr)
		if err != nil {
			log.Printf("Failed to convert PID to integer: %v", err)
			continue
		}

		if err := killProcess(pidInt); err != nil {
			log.Printf("Failed to kill process with PID %d: %v", pidInt, err)
			continue
		}
	}

	fmt.Fprintf(w, "All processes on port %d have been killed", port)
}
func KillPortHandler(w http.ResponseWriter, r *http.Request) {
	portStr := r.URL.Query().Get("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, "Invalid port number", http.StatusBadRequest)
		return
	}
	killAllProcessesOnPort(port, w)
}
func ListPortsHandler(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	// 设置默认值
	start := 56711
	end := 56730

	// 如果参数存在，则尝试转换为整数，否则使用默认值
	if startStr != "" {
		start, _ = strconv.Atoi(startStr)
	}

	if endStr != "" {
		end, _ = strconv.Atoi(endStr)
	}

	ports, err := listEnabledPorts(start, end)
	if err != nil {
		http.Error(w, "Failed to list ports", http.StatusInternalServerError)
		return
	}

	// 构造JSON响应结构体
	response := PortRangeResponse{
		Start:        start,
		End:          end,
		EnabledPorts: ports,
	}

	// 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 编码并写入响应体
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		return
	}
}
