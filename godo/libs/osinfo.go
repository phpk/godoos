package libs

import "runtime"

const (
	OSDarwin  = "macOS"
	OSWindows = "Windows"
	OSLinux   = "Linux"

	ArchAMD64 = "x86_64"
	ArchARM64 = "ARM64"
	ArchARM   = "ARM"
)

// GetOSAndArch 返回当前操作系统的名称和架构的名称
func GetOSAndArch() (osName, archName string) {
	// 获取操作系统
	switch runtime.GOOS {
	case "darwin":
		osName = OSDarwin
	case "windows":
		osName = OSWindows
	case "linux":
		osName = OSLinux
	default:
		osName = "Unknown OS"
	}

	// 获取架构
	switch runtime.GOARCH {
	case "amd64":
		archName = ArchAMD64
	case "arm64":
		archName = ArchARM64
	case "arm":
		archName = ArchARM
	default:
		archName = "Unknown Architecture"
	}

	return osName, archName
}
