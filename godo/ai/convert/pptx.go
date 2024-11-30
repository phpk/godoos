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
	"os"
	"strings"
)

// ConvertPptx converts an MS PowerPoint pptx file to text.
func ConvertPptx(r io.Reader) (string, error) {
	var size int64

	// Common case: if the reader is a file (or trivial wrapper), avoid
	// loading it all into memory.
	var ra io.ReaderAt
	if f, ok := r.(interface {
		io.ReaderAt
		Stat() (os.FileInfo, error)
	}); ok {
		si, err := f.Stat()
		if err != nil {
			return "", err
		}
		size = si.Size()
		ra = f
	} else {
		b, err := io.ReadAll(r)
		if err != nil {
			return "", nil
		}
		size = int64(len(b))
		ra = bytes.NewReader(b)
	}

	zr, err := zip.NewReader(ra, size)
	if err != nil {
		return "", fmt.Errorf("could not unzip: %v", err)
	}

	zipFiles := mapZipFiles(zr.File)

	contentTypeDefinition, err := getContentTypeDefinition(zipFiles["[Content_Types].xml"])
	if err != nil {
		return "", err
	}

	var textBody string
	for _, override := range contentTypeDefinition.Overrides {
		f := zipFiles[override.PartName]

		switch override.ContentType {
		case "application/vnd.openxmlformats-officedocument.presentationml.slide+xml",
			"application/vnd.openxmlformats-officedocument.drawingml.diagramData+xml":
			body, err := parseDocxText(f)
			if err != nil {
				return "", fmt.Errorf("could not parse pptx: %v", err)
			}
			textBody += body + "\n"
		}
	}
	// 在成功解析ZIP文件后，添加图片提取逻辑
	images, err := findImagesInZip(zr)
	if err != nil {
		fmt.Printf("Error extracting images: %v", err)
	}
	fmt.Printf("Images: %v", images)

	return strings.TrimSuffix(textBody, "\n"), nil
}
