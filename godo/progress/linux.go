//go:build !windows

package progress

import (
	// 导入包
	"os/exec"
)

func setHideConsoleCursor(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
