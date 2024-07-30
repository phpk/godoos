//go:build !windows

package store

import (
	// 导入包

	"os/exec"
)

func SetHideConsoleCursor(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
