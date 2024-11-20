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

package files

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"godo/libs"
	"io"
	"os"
	"path/filepath"
	"strings"

	"archive/tar"
	"archive/zip"
)

// Decompress 解压文件到指定目录，支持多种压缩格式
func Decompress(compressedFilePath, destPath string) (string, error) {
	ext := filepath.Ext(compressedFilePath)
	switch ext {
	case ".zip":
		return Unzip(compressedFilePath, destPath)
	case ".tar":
		tarFile, err := os.Open(compressedFilePath)
		if err != nil {
			return "", err
		}
		defer tarFile.Close()
		return Untar(tarFile, destPath)
	case ".gz":
		return Untargz(compressedFilePath, destPath)
	case ".bz2":
		return Untarbz2(compressedFilePath, destPath)
	default:
		return "", fmt.Errorf("unsupported compression format: %s", ext)
	}
}

// Untar 解压.tar文件到指定目录
func Untar(reader io.Reader, destPath string) (string, error) {
	var topLevelDirName string
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Determine the top-level directory name
		if header.Typeflag == tar.TypeDir && topLevelDirName == "" {
			topLevelDirName = header.Name
		} else {
			topLevelDirName = destPath
		}

		target := filepath.Join(destPath, header.Name)

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return "", err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return "", err
		}
	}

	// Return the full path of the top-level directory
	topLevelDestPath := filepath.Join(destPath, topLevelDirName)
	return topLevelDestPath, nil
}

// Untargz 解压.tar.gz文件到指定目录
func Untargz(targzFilePath, destPath string) (string, error) {
	gzFile, err := os.Open(targzFilePath)
	if err != nil {
		return "", err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	return Untar(gzReader, destPath)
}

// Untarbz2 解压.tar.bz2文件到指定目录
func Untarbz2(tarbz2FilePath, destPath string) (string, error) {
	bz2File, err := os.Open(tarbz2FilePath)
	if err != nil {
		return "", err
	}
	defer bz2File.Close()

	bz2Reader := bzip2.NewReader(bz2File)
	return Untar(bz2Reader, destPath)
}

// DecompressZip 解压zip文件到指定目录
func Unzip(zipFilePath, destPath string) (string, error) {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// 初始化顶层目录名称为空
	var topLevelDirName string
	var topLevelDestPath string
	// 遍历ZIP文件中的每个条目
	for _, f := range r.File {
		// 如果是目录且不是根目录，则更新顶层目录名称
		if f.FileInfo().IsDir() && f.Name != "" {
			// 去掉末尾的斜杠
			cleanedName := strings.TrimSuffix(f.Name, string(os.PathSeparator))
			if topLevelDirName == "" || cleanedName < topLevelDirName {
				topLevelDirName = cleanedName
			}
		}
	}

	// 检查是否找到了顶层目录
	if topLevelDirName == "" {
		topLevelDestPath = destPath
	} else {
		// 完整的顶层目录路径
		topLevelDestPath = filepath.Join(destPath, topLevelDirName)
	}

	// 再次遍历ZIP文件中的每个条目，这次进行解压操作
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		// 调整路径以确保所有文件都在顶层目录内
		//path := filepath.Join(topLevelDestPath, strings.TrimPrefix(f.Name, topLevelDirName+string(os.PathSeparator)))
		path := filepath.Join(topLevelDestPath, strings.TrimPrefix(f.Name, topLevelDirName))

		// 确保父目录存在
		if !libs.PathExists(filepath.Dir(path)) {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return "", err
			}
		}

		if f.FileInfo().IsDir() {
			// 创建目录
			if !libs.PathExists(path) {
				if err := os.Mkdir(path, f.Mode()); err != nil {
					return "", err
				}
			}
		} else {
			// 创建并写入文件
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return "", err
			}
		}
	}

	// 返回完整的顶层目录路径
	return topLevelDestPath, nil
}

// IsSupportedCompressionFormat 判断文件后缀是否为支持的压缩格式
func IsSupportedCompressionFormat(filename string) bool {
	supportedFormats := map[string]bool{
		".zip": true,
		".tar": true,
		".gz":  true,
		".bz2": true,
	}

	for format := range supportedFormats {
		if strings.HasSuffix(strings.ToLower(filename), format) {
			return true
		}
	}
	return false
}

// HandlerFile 根据文件后缀判断是否支持解压，支持则解压，否则移动文件
func HandlerFile(filePath string, destDir string) (string, error) {
	// log.Printf("HandlerFile: %s", filePath)
	// log.Printf("destDir: %s", destDir)
	if IsSupportedCompressionFormat(filePath) {
		return Decompress(filePath, destDir)
	} else {
		return "", fmt.Errorf("unsupported compression format: %s", filePath)
	}

}
