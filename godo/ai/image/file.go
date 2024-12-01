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
package image

import (
	"fmt"
	"godo/ai/config"
	"godo/libs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetRandImgs(num int) ([]string, error) {
	var seedList []string

	imagePath, err := GetImageDir()
	if err != nil {
		return seedList, err
	}
	//根据requestBody.Num 生成num个随机数切片
	nowNum := strconv.FormatInt(time.Now().UnixNano(), 10)
	for i := 1; i <= num; i++ {
		filename := fmt.Sprintf("txt2img_%s.png", nowNum)
		if i > 1 {
			filename = fmt.Sprintf("txt2img_%s_%d.png", nowNum, i)
		}
		numPath := filepath.Join(imagePath, filename)
		seedList = append(seedList, numPath)
	}

	return seedList, nil
}
func GetOutputFiles(num int) ([]string, error) {
	if num < 1 {
		return nil, fmt.Errorf("num must be at least 1")
	}
	prefix := "output"
	var filenames []string
	for i := 1; i <= num; i++ {
		suffix := ""
		if i > 1 {
			suffix = fmt.Sprintf("_%d", i)
		}
		filename := fmt.Sprintf("%s%s.png", prefix, suffix)
		tmpfile, err := os.CreateTemp("", filename)
		if err != nil {
			// If any creation fails, clean up and return the error.
			for _, file := range filenames {
				os.Remove(file)
			}
			return nil, err
		}
		defer tmpfile.Close() // Defer closing until after the loop or error handling.

		absFilePath, _ := filepath.Abs(tmpfile.Name())
		filenames = append(filenames, absFilePath)
	}

	return filenames, nil
}

func GetModelPath(modelPath string, fileName string) (string, error) {
	baseDir, err := config.GetHfModelDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(baseDir, modelPath, fileName)
	if !libs.PathExists(filePath) {
		return "", fmt.Errorf("model not found")
	}
	return filePath, nil
}

func GetImageDir() (string, error) {
	userDir, err := libs.GetUserDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	savePath := filepath.Join(userDir, "Pictures", time.Now().Format("2006-01-02"))
	if !libs.PathExists(savePath) {
		os.MkdirAll(savePath, 0755)
	}
	return savePath, nil
}
