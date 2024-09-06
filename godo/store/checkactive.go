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
