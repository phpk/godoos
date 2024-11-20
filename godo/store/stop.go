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
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/process"
)

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
func StopCmd(name string) error {
	processesMu.Lock()
	defer processesMu.Unlock()
	cmd, ok := processes[name]
	if !ok {
		return fmt.Errorf("Process not found")
	}
	//processes[name].Running = false
	if cmd.ProgressName != "" {
		err := KillProcessByName(cmd.ProgressName)
		if err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	} else {
		err := KillByPid(cmd.Pid)
		if err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}
	delete(processes, name)
	return nil
}
func StopProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := StopCmd(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s stopped.", name)
}
func StopAllHandler() error {
	processesMu.Lock()
	defer processesMu.Unlock()

	for name, cmd := range processes {
		if err := KillByPid(cmd.Cmd.Process.Pid); err != nil {
			return fmt.Errorf("failed to stop process %s: %v", name, err)
		}
		//processes[name].Running = false
		delete(processes, name)
	}
	return nil
}
func StopAll(w http.ResponseWriter, r *http.Request) {
	err := StopAllHandler()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All processes stopped.")
}

// Restart a specific process by stopping and then starting it.
func ReStartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Stop the process first
	if _, ok := processes[name]; ok {
		if err := StopCmd(name); err != nil {
			// respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop process %s before restart: %v", name, err))
			// return
			log.Printf("Failed to stop process %s before restart: %v", name, err)
		}
		// processesMu.Lock()
		// defer processesMu.Unlock()
		//processes[name].Running = false
	} else {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Process %s not found to restart", name))
		return
	}
	time.Sleep(time.Second * 2)
	// Start the process again
	err := ExecuteScript(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restart process %s: %v", name, err))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Process %s restarted.", name)
}
func killByPids(pids []int) error {
	for _, pid := range pids {
		if err := KillByPid(pid); err != nil {
			return fmt.Errorf("failed to kill process with PID %d: %w", pid, err)
		}
		fmt.Printf("Killed process with PID %d\n", pid)
	}
	return nil
}

// KillProcessByName 终止所有与给定名称匹配的进程
func KillProcessByName(processName string) error {
	switch runtime.GOOS {
	case "windows":
		// Windows 系统下的实现
		pids, err := findPidsWindows(processName)
		if err != nil {
			return err
		}
		err = killByPids(pids)
		if err != nil {
			return err
		}
	default:
		// Unix/Linux 系统下的实现
		pids, err := findPidsUnix(processName)
		if err != nil {
			return err
		}
		err = killByPids(pids)
		if err != nil {
			return err
		}
	}
	return nil
}

// findPidsUnix 在 Unix/Linux 系统下查找具有指定名称的进程的 PID
func findPidsUnix(name string) ([]int, error) {
	cmd := exec.Command("pgrep", "-f", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute pgrep: %w", err)
	}
	pids := strings.Fields(string(output))
	result := make([]int, len(pids))
	for i, pidStr := range pids {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert PID to integer: %w", err)
		}
		result[i] = pid
	}
	return result, nil
}

// findPidsWindows 在 Windows 系统下查找具有指定名称的进程的 PID
func findPidsWindows(name string) ([]int, error) {
	cmd := exec.Command("tasklist")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute tasklist: %w", err)
	}
	lines := strings.Split(string(output), "\r\n")
	var pids []int
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			if strings.Contains(strings.ToLower(fields[0]), strings.ToLower(name)) {
				pid, err := strconv.Atoi(fields[1])
				if err != nil {
					return nil, fmt.Errorf("failed to convert PID to integer: %w", err)
				}
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}
