/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func CheckActive(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	activeCheckWg := sync.WaitGroup{}
	activeCheckWg.Add(1)

	go func() {
		defer activeCheckWg.Done()
		for {
			select {
			case <-ticker.C:
				checkInactiveProcesses()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func checkInactiveProcesses() {
	processesMu.RLock()
	defer processesMu.RUnlock()

	for name, p := range processes {
		pid := p.Pid
		_, err := os.FindProcess(pid)
		if err != nil {
			if os.IsNotExist(err) {
				p.Running = false
				log.Printf("Process %s (PID: %d) not found.", name, pid)
				if p.IsOn {
					go ExecuteScript(name)
				}
			} else {
				log.Printf("Error finding process %s (PID: %d): %v", name, pid, err)
			}
			continue
		}

		// 进程仍然运行，更新LastPing
		p.LastPing = time.Now()
	}
}

func IsProcessRunning(appName string) (bool, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "Get-Process", appName)
	case "linux":
		fallthrough
	case "darwin":
		cmd = exec.Command("pgrep", "-f", appName)
	default:
		return false, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	var output []byte
	var err error
	// if cmd.SysProcAttr == nil {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{}
	// }
	//cmd.SysProcAttr.HideWindow = true // For Windows to hide the console window
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	if output, err = cmd.CombinedOutput(); err != nil {
		return false, fmt.Errorf("error checking process: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}
