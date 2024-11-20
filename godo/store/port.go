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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
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

//var processInfoRegex = regexp.MustCompile(`(\d+)\s+.*:\s*(\d+)\s+.*LISTEN\s+.*:(\d+)`)

func listAllProcesses() ([]ProcessSystemInfo, error) {
	osType := runtime.GOOS

	var cmd *exec.Cmd
	var output []byte
	var err error

	switch osType {
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-i")
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
	//lsofRegex := regexp.MustCompile(`(\d+)\s+(\S+\s+\S+)\s+(\d+)\s+(\w+)\s+IPv4\s+(\w+)\s+(\w+)\s+TCP\s+(\S+)\->\s+(\S+)\s+\((\w+)\)`)

	processes := make([]ProcessSystemInfo, 0)
	// 初始化映射用于去重
	seenPIDs := make(map[int]bool)
	// 解析输出
	switch osType {
	case "darwin", "linux":
		scanner := bufio.NewScanner(bytes.NewBuffer(output))
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) >= 9 && fields[7] == "TCP" {
				// 解析 PID 和其他字段
				pid, _ := strconv.Atoi(fields[1])
				name := fields[0]
				// 解析本地端口
				localPortStr := strings.Split(fields[8], ":")[1]
				localPort, _ := strconv.Atoi(localPortStr)
				if !seenPIDs[pid] {
					seenPIDs[pid] = true
					processes = append(processes, ProcessSystemInfo{
						PID:   pid,
						Port:  localPort,
						Proto: "TCP",
						Name:  name,
					})
				}
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
				if !seenPIDs[pid] {
					seenPIDs[pid] = true
					processes = append(processes, ProcessSystemInfo{
						PID:   pid,
						Port:  portInt,
						Proto: fields[0],
						Name:  processName,
					})
				}
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
		// 使用 pgrep 查找进程 PID
		pgrepCmd := exec.Command("pgrep", "-f", name)
		pgrepOutput, err := pgrepCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to find process with name %s: %v", name, err)
		}

		// 将输出转换为字符串并按行分割
		pids := strings.Split(strings.TrimSpace(string(pgrepOutput)), "\n")

		// 对于每个找到的 PID，使用 kill 命令杀死进程
		for _, pidStr := range pids {
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				log.Printf("Failed to convert PID to integer: %v", err)
				continue
			}
			// 每次循环使用新的 *exec.Cmd 实例
			killCmd := exec.Command("kill", "-9", strconv.Itoa(pid))
			if err := killCmd.Run(); err != nil {
				log.Printf("Failed to kill process with PID %d: %v", pid, err)
			}
		}
	case "windows":
		cmd = exec.Command("taskkill", "/IM", name, "/F") // /F 表示强制结束
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to kill process with name %s: %v", name, err)
		}
	default:
		return fmt.Errorf("unsupported operating system")
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
