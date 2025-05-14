//go:build darwin
// +build darwin

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
	data.Createtime = time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	data.Modifytime = time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
	data.Accesstime = time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)

	return true, nil
}
