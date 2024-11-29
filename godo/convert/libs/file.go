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
package libs

import (
	"errors"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func getXpdfDir(exename string) (string, error) {
	convertDir, err := getConvertDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(convertDir, "pdf", exename+".exe")
	} else {
		path = filepath.Join(convertDir, "pdf", exename)
	}
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("pdf convert exe not found")
	}
}
func getRapidDir() (string, error) {
	convertDir, err := getConvertDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(convertDir, "rapid", "RapidOcrOnnx.exe")
	} else {
		path = filepath.Join(convertDir, "rapid", "RapidOcrOnnx")
	}
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("RapidOcrOnnx not found")
	}
}

func getRapidModelDir() (string, error) {
	convertDir, err := getConvertDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	path := filepath.Join(convertDir, "rapid", "models")
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("RapidOcrOnnx model not found")
	}
}
func getConvertDir() (string, error) {
	runDir, err := libs.GetAiRunDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(runDir, "goconv"), nil
}

func GetTempDir(pathname string) (string, error) {
	tempDir, err := os.MkdirTemp("", pathname)
	if err != nil {
		log.Println("Failed to create temporary directory:", err)
		return "./", err
	}

	log.Println("Temporary directory created:", tempDir)
	// defer func() {
	// 	os.RemoveAll(tempDir)
	// }()
	return tempDir, nil
}
func GetTempFile(r io.Reader, prename string) (string, *os.File, error) {
	// 创建临时文件
	tmpfile, err := os.CreateTemp("", prename)

	if err != nil {
		return "", tmpfile, err
	}

	// 将Reader内容写入临时文件
	if _, err := io.Copy(tmpfile, r); err != nil {
		return "", tmpfile, err
	}

	// 获取临时文件的绝对路径
	absFilePath, err := filepath.Abs(tmpfile.Name())
	if err != nil {
		return "", tmpfile, err
	}
	return absFilePath, tmpfile, nil
}
func CloseTempFile(tmpfile *os.File) {
	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name()) // 根据需要决定是否删除临时文件
	}()
}
