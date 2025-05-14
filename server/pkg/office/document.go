package office

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"godocms/common"
	"godocms/libs"
	lb "godocms/pkg/office/libs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type DocResult struct {
	filePath    string
	newFilePath string
	err         error
}

func SetDocument(dirPath string, knowledgeId uint, splitSize int) error {
	if !libs.PathExists(dirPath) {
		return nil
	}
	var wg sync.WaitGroup
	results := make(chan DocResult, 100) // 缓冲通道
	if splitSize == 0 {
		splitSize = 512
	}
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 获取文件名
			fileName := filepath.Base(filePath)
			// 检查文件名是否以点开头
			if len(fileName) > 0 && fileName[0] == '.' {
				return nil // 跳过以点开头的文件
			}
			// 获取文件扩展名
			ext := filepath.Ext(filePath)
			// 检查文件扩展名是否为 .exe
			if ext == ".exe" {
				return nil // 跳过 .exe 文件
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				result := ProcessFile(filePath, knowledgeId, splitSize)
				results <- result
			}(filePath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			fmt.Printf("Failed to process file %s: %v\n", result.filePath, result.err)
		} else {
			fmt.Printf("Processed file %s and saved JSON to %s\n", result.filePath, result.newFilePath)
		}
	}

	return nil
}

func ProcessFile(filePath string, knowledgeId uint, splitSize int) DocResult {
	doc, err := GetDocument(filePath, splitSize)
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	//newFileName := ".godocloud." + filepath.Base(filePath)
	newFileName := fmt.Sprintf(".godocloud.%d.%s.json", knowledgeId, filepath.Base(filePath))
	newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)
	// log.Printf("New file name: %s", newFileName)
	// log.Printf("New file path: %s", newFilePath)
	err = os.WriteFile(newFilePath, jsonData, 0644)
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	return DocResult{filePath: filePath, newFilePath: newFilePath, err: nil}
}

// ProcessBase64File 处理解码后的文件并提取文本信息
func ProcessBase64File(base64String string, fileName string) (string, error) {
	// 解码 Base64 字符串
	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %v", err)
	}
	// 获取文件后缀并转换为小写
	// fileExt := strings.ToLower(filepath.Ext(fileName))
	// if fileExt == "" {
	// 	return "", fmt.Errorf("file extension not found in filename: %s", fileName)
	// }
	// log.Printf("File type: %s\n", fileExt)
	cacheDir := common.GetCacheDir()
	tempFilePath := filepath.Join(cacheDir, fileName)

	// 创建临时文件
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// 将解码后的数据写入临时文件
	_, err = tempFile.Write(decodedBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}

	// 获取文档内容
	doc, err := GetDocument(tempFilePath, 300)
	if err != nil {
		return "", fmt.Errorf("failed to get document: %v", err)
	}
	log.Printf("Document content: %s\n", doc.Content)

	// 删除临时文件
	defer os.Remove(tempFilePath)

	// 提取文本内容
	return doc.Content, nil
}
func GetDocument(pathname string, splitSize int) (*Document, error) {
	if !libs.PathExists(pathname) {
		return nil, fmt.Errorf("file does not exist: %s", pathname)
	}
	abPath, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}
	filename := path.Base(pathname)
	data := Document{path: pathname, RePath: abPath, Title: filename}

	extension := path.Ext(pathname)

	_, err = getFileInfoData(&data)
	if err != nil {
		return &data, err
	}
	res, err := getMetaData(&data, extension)
	if !res {
		fmt.Printf("⚠️ %s", err.Error())
		return &data, err
	}
	switch extension {
	case ".docx":
		_, err = getContentData(&data, docx2txt)
	case ".pptx":
		_, err = getContentData(&data, pptx2txt)
	case ".xlsx":
		_, err = getContentData(&data, xlsx2txt)
	case ".pdf":
		_, err = getContentData(&data, pdf2txt)
	case ".doc":
		_, err = getContentData(&data, doc2txt)
	case ".ppt":
		_, err = getContentData(&data, ppt2txt)
	case ".xls":
		_, err = getContentData(&data, xls2txt)
	case ".epub":
		_, err = getContentData(&data, epub2txt)
	case ".odt":
		_, err = getContentData(&data, odt2txt)
	case ".xml":
		_, err = getContentData(&data, xml2txt)
	case ".rtf":
		_, err = getContentData(&data, rtf2txt)
	case ".md":
		_, err = getContentData(&data, md2txt)
	case ".txt":
		_, err = getContentData(&data, text2txt)
	case ".xhtml", ".html", ".htm":
		_, err = getContentData(&data, html2txt)
	case ".json":
		_, err = getContentData(&data, json2txt)
	case ".jpg", ".jpeg", ".jpe", ".jfif", ".jfif-tbnl", ".png", ".gif", ".bmp", ".webp", ".tif", ".tiff":
		_, err = getContentData(&data, lb.GetImageContent)
	default:
		_, err = getContentData(&data, text2txt)
	}
	if splitSize == 0 {
		splitSize = 512
	}
	data.Split = SplitText(data.Content, splitSize)
	fileTitle := strings.TrimSuffix(filepath.Base(pathname), extension)
	log.Printf("Split fileTitle: %v\n", fileTitle)
	for i := range data.Split {
		data.Split[i] = fileTitle + ":" + data.Split[i]
	}
	if err != nil {
		return &data, err
	}
	return &data, nil
}
func GetTxtDoc(pathname string) (*Document, error) {
	if !libs.PathExists(pathname) {
		return nil, fmt.Errorf("file does not exist: %s", pathname)
	}
	abPath, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}
	filename := path.Base(pathname)
	data := Document{path: pathname, RePath: abPath, Title: filename}
	_, err = getFileInfoData(&data)
	if err != nil {
		return &data, err
	}
	extension := path.Ext(pathname)
	res, err := getMetaData(&data, extension)
	if !res {
		fmt.Printf("⚠️ %s", err.Error())
		return &data, err
	}
	_, err = getContentData(&data, text2txt)
	if err != nil {
		return &data, err
	}
	return &data, nil
}

// Read the content of office files and insert into the interface
func getContentData(data *Document, reader DocReader) (bool, error) {
	content, err := reader(data.path)
	if err != nil {
		return false, err
	}
	data.Content = content

	return true, nil
}
