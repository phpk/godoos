package progress

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func CheckActive(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
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
		// 检查进程是否还在运行
		pid := p.Cmd.Process.Pid
		_, err := os.FindProcess(pid)
		if err != nil {
			p.Running = false
			// 进程已不存在，跳过后续检查
			continue
		}

		// 如果进程还在运行，检查/ping接口
		if p.Running {
			resp, err := http.Get(p.PingURL)
			if err != nil {
				p.Running = false
				// 进程未响应
				log.Printf("Error checking ping for process %s: %v", name, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				p.Running = false
				// 进程未正常响应，停止进程
				if err := p.Cmd.Process.Signal(os.Kill); err != nil {
					log.Printf("Failed to stop process %s before restart: %v", name, err)
					continue
				}
			} else {
				// 进程正常响应，更新LastPingAt并继续检查下一个进程
				p.LastPing = time.Now()
			}
		}
	}
}
