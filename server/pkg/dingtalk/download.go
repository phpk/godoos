package dingtalk

import (
	"fmt"
	"godocms/common"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cast"
)

type WriteFileLog struct {
	Action   string            `json:"action"`
	FileName string            `json:"file_name"`
	Path     string            `json:"path"`
	Remote   string            `json:"remote"`
	Header   map[string]string `json:"header"`
	Unionid  string            `json:"unionid"`
	SpaceID  string            `json:"space_id"`
	FileID   string            `json:"file_id"`
	Time     string            `json:"time"`
}

func DownloadFile(httpClient *http.Client, fl *WriteFileLog, history map[string]struct{}) (*WriteFileLog, error) {
	slog.Info("开始下载文件", "url", fl.Remote, "路径", fl.Path)
	if _, ok := history[fl.Path]; ok {
		return nil, nil
	}
	req, err := http.NewRequest("GET", fl.Remote, nil)
	if err != nil {
		return fl, err
	}

	for key, value := range fl.Header {
		req.Header.Add(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fl, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fl, fmt.Errorf("failed to download file[%s], ststus[%s]", fl.Path, resp.Status)
	}

	// 检查文件是否存在
	dingSpaceFileInfo, err := os.Stat(fl.Path)
	if err == nil && dingSpaceFileInfo.Size() > 0 {
		slog.Info("文件已存在", "filepath", fl.Path)
		return nil, nil
	}

	// 如果文件不存在，则创建新文件
	outFile, err := os.Create(fl.Path)
	if err != nil {
		return fl, err
	}

	// 确保函数退出时删除文件
	defer func() {
		if err != nil {
			os.Remove(fl.Path)
		}
		outFile.Close()
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fl, err
	}

	history[fl.Path] = struct{}{}
	slog.Info("文件已下载", "filepath", fl.Path)
	return nil, nil
}

func ConverWriteFileLog(downLoadLogs any) []*WriteFileLog {
	logs := make([]*WriteFileLog, 0)
	if log, ok := downLoadLogs.([]any); ok {
		for _, file := range log {
			if fileLog, ok := file.(map[string]any); ok {
				logs = append(logs, &WriteFileLog{
					Action:   fileLog["action"].(string),
					FileID:   fileLog["file_id"].(string),
					SpaceID:  fileLog["space_id"].(string),
					FileName: fileLog["file_name"].(string),
					Header:   cast.ToStringMapString(fileLog["header"]),
					Path:     fileLog["path"].(string),
					Remote:   fileLog["remote"].(string),
					Time:     fileLog["time"].(string),
					Unionid:  fileLog["unionid"].(string),
				})
			}
		}
	}

	return logs
}

func ReadSpaceFileCache(unionid string) []*WriteFileLog {
	data, _ := common.Cache.Get(fmt.Sprintf("ding_space_%s", unionid))
	logs := ConverWriteFileLog(data)
	return logs
}

func ReadDownLoadFailCache() []*WriteFileLog {
	data, _ := common.Cache.Get("ding_space_failload")
	tasks := ConverWriteFileLog(data)

	return tasks
}

func ExecDownLoadSpaceTask(tasks []*WriteFileLog) ([]*WriteFileLog, error) {
	downloadHistory := make(map[string]struct{}, 0)
	client := &http.Client{}
	failHistory := make([]*WriteFileLog, 0)
	for _, task := range tasks {
		// 开始下载文件, 返回下载失败的任务
		if failLog, err := DownloadFile(client, task, downloadHistory); err != nil {
			failHistory = append(failHistory, failLog)
			slog.Error("下载文件失败", "err", err)
		}
	}

	originFails := ReadDownLoadFailCache()
	newFails := getNewFailCache(originFails, failHistory)

	// 每次下载失败的文件写入缓存
	common.Cache.Set("ding_space_failload", newFails, 60*24*time.Minute)

	return failHistory, nil
}

// 找出arr1 在 arr2 中 没有的，并加入arr2
func getNewFailCache(arr1, arr2 []*WriteFileLog) []*WriteFileLog {
	set := make(map[string]bool)
	for _, item := range arr2 {
		set[item.FileID] = true
	}

	// 将arr1中没有在arr2中出现的元素追加到arr2中
	for _, item := range arr1 {
		if !set[item.FileID] {
			arr2 = append(arr2, item)
		}
	}

	return arr2
}
