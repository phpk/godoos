//go:build windows
// +build windows

package progress

import (
	// 导入包
	"os/exec"
	"syscall"
)

func SetHideConsoleCursor(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
