package deps

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// extractZip 解压嵌入的 zip 文件到目标目录
func extractZip(zipData []byte, targetDir string) error {
	log.Println("Extracting embedded ZIP file...")
	// 使用内存缓冲区来读取嵌入的ZIP数据
	reader := bytes.NewReader(zipData)
	zipReader, err := zip.NewReader(reader, int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %v", err)
	}

	// 遍历ZIP文件中的每个条目并解压
	for _, zipEntry := range zipReader.File {
		// 忽略 __MACOSX 文件夹
		if strings.HasPrefix(zipEntry.Name, "__MACOSX") {
			continue
		}

		// 检查条目名称是否以"."开头，如果是，则跳过
		if strings.HasPrefix(zipEntry.Name, ".") {
			fmt.Printf("Skipping hidden entry: %s\n", zipEntry.Name)
			continue
		}

		// 去掉 ZIP 文件中的顶层目录(zip文件名本身)
		entryName := zipEntry.Name
		if idx := strings.Index(entryName, "/"); idx != -1 {
			entryName = entryName[idx+1:] // 去掉顶层目录
		}

		// 如果 entryName 为空，说明是顶层目录本身，跳过
		if entryName == "" {
			continue
		}

		// 构建解压后的文件或目录路径
		entryPath := filepath.Join(targetDir, entryName)
		fmt.Println("Extracting:", entryName)

		// 如果是目录，则创建目录
		if zipEntry.FileInfo().IsDir() {
			if err := os.MkdirAll(entryPath, zipEntry.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		// 如果是文件，则解压文件
		zipFile, err := zipEntry.Open()
		if err != nil {
			return fmt.Errorf("failed to open zip file entry: %v", err)
		}
		defer zipFile.Close()

		// 确保目标文件的父目录存在
		if err := os.MkdirAll(filepath.Dir(entryPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %v", err)
		}

		dstFile, err := os.OpenFile(entryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %v", err)
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, zipFile); err != nil {
			return fmt.Errorf("failed to copy content to destination file: %v", err)
		}
	}

	log.Println("Embedded ZIP extracted to", targetDir)

	return nil
}
