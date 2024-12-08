//go:build windows
// +build windows

package office

import (
	"os"
	"syscall"
	"time"
)

func getFileInfoData(data *Document) (bool, error) {
	fileinfo, err := os.Stat(data.path)
	if err != nil {
		return false, err
	}
	data.Filename = fileinfo.Name()
	data.Title = data.Filename
	data.Size = int(fileinfo.Size())

	stat := fileinfo.Sys().(*syscall.Win32FileAttributeData)
	data.Createtime = time.Unix(0, stat.CreationTime.Nanoseconds())
	data.Modifytime = time.Unix(0, stat.LastWriteTime.Nanoseconds())
	data.Accesstime = time.Unix(0, stat.LastAccessTime.Nanoseconds())

	return true, nil
}
