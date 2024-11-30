/*
 * GodoOS - A lightweight cloud desktop
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
package convert

import (
	"crypto/md5"
	"fmt"
	lib "godo/ai/convert/libs"
	"godo/libs"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ResContentInfo struct {
	Content string       `json:"content"`
	Images  []ImagesInfo `json:"image"`
}
type ImagesInfo struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// 计算文件的MD5哈希值
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// 复制图片并检查大小和MD5
func CopyImages(destDir string) ([]ImagesInfo, error) {
	copiedFiles := []ImagesInfo{}
	srcDir, err := libs.GetTrueCacheDir()
	if !libs.PathExists(srcDir) {
		return copiedFiles, fmt.Errorf("source directory does not exist: %s", srcDir)
	}
	if err != nil {
		return copiedFiles, fmt.Errorf("failed to create temporary cache directory: %w", err)
	}
	if !libs.PathExists(destDir) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return copiedFiles, fmt.Errorf("failed to create destination directory: %w", err)
		}
	}
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if isImageExtension(ext) {
				destPath := filepath.Join(destDir, info.Name())

				// 检查目标文件是否存在且大小相同
				if fileInfo, err := os.Stat(destPath); err == nil {
					if fileInfo.Size() == info.Size() {
						// 文件大小相同，进一步检查MD5
						srcHash, err := calculateFileHash(path)
						if err != nil {
							log.Printf("Error calculating source hash for %s: %v", path, err)
							return err
						}
						destHash, err := calculateFileHash(destPath)
						if err != nil {
							log.Printf("Error calculating destination hash for %s: %v", destPath, err)
							return err
						}
						if srcHash == destHash {
							fmt.Printf("Skipping %s because a file with the same size and content already exists.\n", path)
							return nil
						}
					}
				}
				paths := []string{path}
				content, err := lib.RunRapid(paths)
				if err != nil {
					content = ""
				}
				// 复制文件
				if err := copyImagesFile(path, destPath); err != nil {
					return err
				}

				copiedFiles = append(copiedFiles, ImagesInfo{Path: destPath, Content: content}) // 记录复制成功的文件路径
				fmt.Printf("Copied %s to %s\n", path, destPath)
			}
		}
		return nil
	})
	defer func() {
		os.RemoveAll(srcDir)
	}()
	if len(copiedFiles) < 1 {
		os.RemoveAll(destDir)
	}
	if err != nil {
		return copiedFiles, err
	}
	return copiedFiles, nil
}

// 辅助函数检查文件扩展名是否为图片
func isImageExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".jpe", ".jfif", ".jfif-tbnl", ".png", ".gif", ".bmp", ".webp", ".tif", ".tiff":
		return true
	default:
		return false
	}
}
func isConvertImageFile(ext string) bool {
	switch ext {
	case ".docx", ".pdf", ".pptx", ".odt":
		return true
	default:
		return false
	}
}

// 复制单个文件
func copyImagesFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, in, 0644)
}
