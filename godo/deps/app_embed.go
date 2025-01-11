package deps

import (
	_ "embed"
	"errors"
	"godo/libs"
)

//go:embed darwin/amd64.zip
var embeddedDarwinAMD64Zip []byte

//go:embed darwin/arm64.zip
var embeddedDarwinARM64Zip []byte

//go:embed linux/amd64.zip
var embeddedLinuxAMD64Zip []byte

//go:embed linux/arm64.zip
var embeddedLinuxARM64Zip []byte

//go:embed windows/amd64.zip
var embeddedWindowsAMD64Zip []byte

//go:embed windows/arm64.zip
var embeddedWindowsARM64Zip []byte

const (
	FRPCAPP = "frpc"
)

var (
	ErrUnsupportedOSArch = errors.New("unsupported OS and architecture combination")
	embeddedZips         = map[string]map[string]map[string][]byte{
		libs.OSDarwin: {
			libs.ArchAMD64: {FRPCAPP: embeddedDarwinAMD64Zip},
			libs.ArchARM64: {FRPCAPP: embeddedDarwinARM64Zip},
			libs.ArchARM:   {FRPCAPP: embeddedDarwinARM64Zip},
		},
		libs.OSWindows: {
			libs.ArchAMD64: {FRPCAPP: embeddedWindowsAMD64Zip},
			libs.ArchARM64: {FRPCAPP: embeddedWindowsARM64Zip},
			libs.ArchARM:   {FRPCAPP: embeddedWindowsARM64Zip},
		},
		libs.OSLinux: {
			libs.ArchAMD64: {FRPCAPP: embeddedLinuxAMD64Zip},
			libs.ArchARM64: {FRPCAPP: embeddedLinuxARM64Zip},
			libs.ArchARM:   {FRPCAPP: embeddedLinuxARM64Zip},
		},
	}
)

func ExtractZip(appName, targetDir string) error {
	os, arch := libs.GetOSAndArch()
	archMap, ok := embeddedZips[os]
	if !ok {
		return ErrUnsupportedOSArch
	}

	zipData, ok := archMap[arch][appName]
	if !ok {
		return ErrUnsupportedOSArch
	}

	return extractZip(zipData, targetDir)
}
