package localchat

import (
	"encoding/base64"
	"fmt"
	"godoos/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 合并文件分片
func mergeFiles(fileName string, totalParts int, chatDir string) error {
	var parts []io.Reader
	for i := 1; i <= totalParts; i++ {
		filePath := fmt.Sprintf("%v%v_%v.part", chatDir, fileName, i)
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open part %d: %w", i, err)
		}
		defer file.Close()
		parts = append(parts, file)
	}

	mergedFilePath := fmt.Sprintf("%vmerged_%v", chatDir, fileName)
	mergedFile, err := os.Create(mergedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create merged file: %w", err)
	}
	defer mergedFile.Close()

	_, err = io.Copy(mergedFile, io.MultiReader(parts...))
	if err != nil {
		return fmt.Errorf("failed to merge files: %w", err)
	}

	// 合并后清理分片文件
	for i := 1; i <= totalParts; i++ {
		os.Remove(fmt.Sprintf("%v%v_%v.part", chatDir, fileName, i))
	}

	return nil
}

func UploadBigFileHandler(w http.ResponseWriter, r *http.Request, msg Message) {

	chatDir, err := GetChatPath()
	if err != nil {
		http.Error(w, "Failed to get chat path", http.StatusInternalServerError)
		return
	}

	// 创建或打开临时文件以写入分片
	tempFilePath := fmt.Sprintf("%v%v_%v.part", chatDir, msg.FileInfo.FileName, msg.FileInfo.PartNumber)
	out, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Failed to create temp file: %v", err)
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将请求体的内容写入临时文件
	_, err = io.Copy(out, r.Body)
	if err != nil {
		log.Printf("Failed to write file part: %v", err)
		http.Error(w, "Failed to write file part", http.StatusInternalServerError)
		return
	}

	// 更新上传状态
	uploadStatus.Lock()
	uploadStatus.Status[msg.FileInfo.FileName]++
	if uploadStatus.Status[msg.FileInfo.FileName] == msg.FileInfo.TotalParts {
		// 所有分片上传完成，触发合并
		go func() {
			err := mergeFiles(msg.FileInfo.FileName, msg.FileInfo.TotalParts, chatDir)
			if err != nil {
				log.Printf("Failed to merge files for %v: %v", msg.FileInfo.FileName, err)
			} else {
				log.Printf("Merged file %v successfully", msg.FileInfo.FileName)
			}
			// 清理状态记录
			delete(uploadStatus.Status, msg.FileInfo.FileName)
			msg.Content = "uploaded"
			messageChan <- msg
		}()
	}
	uploadStatus.Unlock()

	// 返回成功响应
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "File part uploaded successfully")
}

// SaveContentToFile 保存内容到文件并返回UploadInfo结构体
func SaveContentToFile(content, fileName string) (UploadInfo, error) {
	uploadBaseDir, err := GetChatPath()
	if err != nil {
		return UploadInfo{}, err
	}
	appDir, err := libs.GetOsDir()
	if err != nil {
		return UploadInfo{}, err
	}

	// 去除文件名中的空格
	fileNameWithoutSpaces := strings.ReplaceAll(fileName, " ", "_")
	fileNameWithoutSpaces = strings.ReplaceAll(fileNameWithoutSpaces, "/", "")
	fileNameWithoutSpaces = strings.ReplaceAll(fileNameWithoutSpaces, `\`, "")
	// 提取文件名和扩展名
	// 查找最后一个点的位置
	lastDotIndex := strings.LastIndexByte(fileNameWithoutSpaces, '.')

	// 如果找到点，则提取扩展名，否则视为没有扩展名
	ext := ""
	if lastDotIndex != -1 {
		ext = fileNameWithoutSpaces[lastDotIndex:]
		fileNameWithoutSpaces = fileNameWithoutSpaces[:lastDotIndex]
	} else {
		ext = ""
	}
	randFileName := fmt.Sprintf("%s_%s%s", fileNameWithoutSpaces, strconv.FormatInt(time.Now().UnixNano(), 10), ext)
	savePath := filepath.Join(uploadBaseDir, randFileName)

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return UploadInfo{}, err
	}

	if err := os.WriteFile(savePath, []byte(content), 0644); err != nil {
		return UploadInfo{}, err
	}
	content = string(content)
	// 检查文件内容是否以"link::"开头
	if !strings.HasPrefix(content, "link::") {
		content = base64.StdEncoding.EncodeToString([]byte(content))
	}
	return UploadInfo{
		Name:      fileNameWithoutSpaces,
		SavePath:  strings.TrimPrefix(savePath, appDir),
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// MultiUploadHandler 处理多文件上传请求
func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10000 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No file parts in the request", http.StatusBadRequest)
		return
	}

	fileInfoList := make([]UploadInfo, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusBadRequest)
			continue
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
			continue
		}
		//log.Printf(string(content))
		// 保存上传的文件内容
		info, err := SaveContentToFile(string(content), fileHeader.Filename)
		if err != nil {
			http.Error(w, "Failed to save uploaded file", http.StatusBadRequest)
			continue
		}
		log.Println(info.SavePath)

		//info.SavePath = savePath
		fileInfoList = append(fileInfoList, info)
	}
	user := UserInfo{
		IP:       r.FormValue("ip"),
		Hostname: r.FormValue("hostname"),
	}
	msg := Message{
		Type:       "file",
		Content:    "file recieved",
		SenderInfo: user,
		FileList:   fileInfoList,
	}
	messageChan <- msg
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File send successfully")
	//serv.Res(serv.Response{Code: 0, Data: fileInfoList}, w)
}
