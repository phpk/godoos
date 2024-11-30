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
	"io"

	"godo/ai/convert/libs"
)

func ConvertImage(r io.Reader) (string, error) {
	// 获取临时文件的绝对路径
	absFilePath, tmpfile, err := libs.GetTempFile(r, "prefix-image")
	if err != nil {
		return "", err
	}
	paths := []string{absFilePath}
	// 识别文本
	output, err := libs.RunRapid(paths)
	if err != nil {
		return "", err
	}
	libs.CloseTempFile(tmpfile)
	// resultString, err := libs.ExtractText(output)
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println(resultString)
	return output, nil
}
