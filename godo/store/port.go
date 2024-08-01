package store

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type ProcessSystemInfo struct {
	PID   int    `json:"pid"`
	Port  int    `json:"port"`
	Proto string `json:"proto"`
	Name  string `json:"name"`
}

type AllProcessesResponse struct {
	Processes []ProcessSystemInfo `json:"processes"`
}

var processInfoRegex = regexp.MustCompile(`(\d+)\s+.*:\s*(\d+)\s+.*LISTEN\s+.*:(\d+)`)

func listAllProcesses() ([]ProcessSystemInfo, error) {
	osType := runtime.GOOS

	var cmd *exec.Cmd
	var output []byte
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-i", "-n", "-P")
	case "windows":
		cmd = exec.Command("netstat", "-ano")
		cmd = SetHideConsoleCursor(cmd)

	default:
		return nil, fmt.Errorf("unsupported operating system")
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list all processes: %v", err)
	}

	processes := make([]ProcessSystemInfo, 0)

	// 解析输出
	switch osType {
	case "darwin", "linux":
		scanner := bufio.NewScanner(bytes.NewBuffer(output)) // 使用bufio.Scanner

		for scanner.Scan() {
			line := scanner.Text()
			matches := processInfoRegex.FindStringSubmatch(line)
			if matches != nil {
				pid, _ := strconv.Atoi(matches[1])
				port, _ := strconv.Atoi(matches[3])
				processName, err := getProcessName(osType, pid)
				if err != nil {
					log.Printf("Failed to get process name for PID %d: %v", pid, err)
					continue
				}
				processes = append(processes, ProcessSystemInfo{
					PID:   pid,
					Port:  port,
					Proto: matches[2],
					Name:  processName,
				})
			}
		}
	case "windows":
		scanner := bufio.NewScanner(bytes.NewBuffer(output))
		for scanner.Scan() {
			line := scanner.Text()
			// 需要针对Windows的netstat输出格式进行解析
			// 示例：TCP    0.0.0.0:80          0.0.0.0:*               LISTENING       1234
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[3] == "LISTENING" {
				_, port, err := net.SplitHostPort(fields[1])
				if err != nil {
					log.Printf("Failed to parse port: %v", err)
					continue
				}
				pid, _ := strconv.Atoi(fields[4])
				processName, err := getProcessName(osType, pid)
				if err != nil {
					log.Printf("Failed to get process name for PID %d: %v", pid, err)
					continue
				}
				portInt, err := strconv.Atoi(port)
				if err != nil {
					log.Printf("Failed to convert port to integer: %v", err)
					continue
				}
				processes = append(processes, ProcessSystemInfo{
					PID:   pid,
					Port:  portInt,
					Proto: fields[0],
					Name:  processName,
				})
			}
		}
	}

	return processes, nil
}

func getProcessName(osType string, pid int) (string, error) {
	var cmd *exec.Cmd
	var output []byte
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	case "windows":
		cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH")
		cmd = SetHideConsoleCursor(cmd)
	default:
		return "", fmt.Errorf("unsupported operating system")
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get process name: %v", err)
	}
	// log.Printf("output: %s", output)
	switch osType {
	case "darwin", "linux":
		return strings.TrimSpace(string(output)), nil
	case "windows":
		parts := strings.Fields(string(output))
		if len(parts) >= 1 {
			return parts[0], nil
		}
		return "", fmt.Errorf("no process name found in output")
	}

	return "", fmt.Errorf("unknown error getting process name")
}

func killProcessByName(name string) error {
	osType := runtime.GOOS

	var cmd *exec.Cmd
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("pkill", name)
	case "windows":
		cmd = exec.Command("taskkill", "/IM", name, "/F") // /F 表示强制结束
		cmd = SetHideConsoleCursor(cmd)
	default:
		return fmt.Errorf("unsupported operating system")
	}

	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to kill process with name %s: %v", name, err)
	}

	return err
}

func KillProcessByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if err := killProcessByName(name); err != nil {
		http.Error(w, fmt.Sprintf("Failed to kill process: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Process '%s' has been killed", name)
}

func ListAllProcessesHandler(w http.ResponseWriter, r *http.Request) {
	processes, err := listAllProcesses()
	if err != nil {
		http.Error(w, "Failed to list all processes", http.StatusInternalServerError)
		return
	}

	response := AllProcessesResponse{
		Processes: processes,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		return
	}
}
