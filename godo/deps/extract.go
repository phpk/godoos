/*
 * GodoAI - A software focused on localizing AI applications
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package deps

import (
	"archive/zip"
	"bytes"
	"fmt"
	"godo/libs"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func InitDir() error {
	// 获取当前用户主目录
	runDir := libs.GetAiExeDir()
	if !libs.PathExists(runDir) {
		if err := os.MkdirAll(runDir, 0o755); err != nil {
			return fmt.Errorf("failed to create user directory: %v", err)
		}
		err := ExtractEmbeddedZip(runDir)
		if err != nil {
			return fmt.Errorf("failed to extract embedded zip: %v", err)
		}

	}

	return nil
}

// ExtractEmbeddedZip 解压嵌入的ZIP文件到指定目录
func ExtractEmbeddedZip(exeDir string) error {
	// 使用内存缓冲区来读取嵌入的ZIP数据
	reader := bytes.NewReader(embeddedZip)
	zipReader, err := zip.NewReader(reader, int64(len(embeddedZip)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %v", err)
	}

	// 遍历ZIP文件中的每个条目并解压
	for _, zipEntry := range zipReader.File {
		// 检查条目名称是否以"."开头，如果是，则跳过
		if strings.HasPrefix(zipEntry.Name, ".") {
			fmt.Printf("Skipping hidden entry: %s\n", zipEntry.Name)
			continue
		}

		// 构建解压后的文件或目录路径
		entryPath := filepath.Join(exeDir, zipEntry.Name)

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

	fmt.Println("Embedded ZIP extracted to", exeDir)
	return nil
}
