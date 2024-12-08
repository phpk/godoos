//go:build linux
// +build linux

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

	stat := fileinfo.Sys().(*syscall.Stat_t)
	data.Createtime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	data.Modifytime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	data.Accesstime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)

	return true, nil
}
