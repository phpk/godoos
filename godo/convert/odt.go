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
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"time"
)

// ConvertODT converts a ODT file to text
func ConvertODT(r io.Reader) (string, error) {
	meta := make(map[string]string)
	var textBody string

	b, err := io.ReadAll(io.LimitReader(r, maxBytes))
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", fmt.Errorf("error unzipping data: %v", err)
	}

	for _, f := range zr.File {
		switch f.Name {
		case "meta.xml":
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("error extracting '%v' from archive: %v", f.Name, err)
			}
			defer rc.Close()

			info, err := XMLToMap(rc)
			if err != nil {
				return "", fmt.Errorf("error parsing '%v': %v", f.Name, err)
			}

			if tmp, ok := info["creator"]; ok {
				meta["Author"] = tmp
			}
			if tmp, ok := info["date"]; ok {
				if t, err := time.Parse("2006-01-02T15:04:05", tmp); err == nil {
					meta["ModifiedDate"] = fmt.Sprintf("%d", t.Unix())
				}
			}
			if tmp, ok := info["creation-date"]; ok {
				if t, err := time.Parse("2006-01-02T15:04:05", tmp); err == nil {
					meta["CreatedDate"] = fmt.Sprintf("%d", t.Unix())
				}
			}

		case "content.xml":
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("error extracting '%v' from archive: %v", f.Name, err)
			}
			defer rc.Close()

			textBody, err = XMLToText(rc, []string{"br", "p", "tab"}, []string{}, true)
			if err != nil {
				return "", fmt.Errorf("error parsing '%v': %v", f.Name, err)
			}
		}
	}
	// 在成功解析ZIP文件后，添加图片提取逻辑
	images, err := findImagesInZip(zr)
	if err != nil {
		fmt.Printf("Error extracting images: %v", err)
	}
	fmt.Printf("Images: %v", images)

	return textBody, nil
}
