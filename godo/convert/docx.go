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
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"
)

type typeOverride struct {
	XMLName     xml.Name `xml:"Override"`
	ContentType string   `xml:"ContentType,attr"`
	PartName    string   `xml:"PartName,attr"`
}

type contentTypeDefinition struct {
	XMLName   xml.Name       `xml:"Types"`
	Overrides []typeOverride `xml:"Override"`
}

// ConvertDocx converts an MS Word docx file to text.
func ConvertDocx(r io.Reader) (string, error) {
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
		b, err := io.ReadAll(io.LimitReader(r, maxBytes))
		if err != nil {
			return "", fmt.Errorf("error read data: %v", err)
		}
		size = int64(len(b))
		ra = bytes.NewReader(b)
	}

	zr, err := zip.NewReader(ra, size)
	if err != nil {
		return "", fmt.Errorf("error unzipping data: %v", err)
	}

	zipFiles := mapZipFiles(zr.File)

	contentTypeDefinition, err := getContentTypeDefinition(zipFiles["[Content_Types].xml"])
	if err != nil {
		return "", err
	}

	meta := make(map[string]string)
	var textHeader, textBody, textFooter string
	for _, override := range contentTypeDefinition.Overrides {
		f := zipFiles[override.PartName]

		switch {
		case override.ContentType == "application/vnd.openxmlformats-package.core-properties+xml":
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("error opening '%v' from archive: %v", f.Name, err)
			}
			defer rc.Close()

			meta, err = XMLToMap(rc)
			if err != nil {
				return "", fmt.Errorf("error parsing '%v': %v", f.Name, err)
			}

			if tmp, ok := meta["modified"]; ok {
				if t, err := time.Parse(time.RFC3339, tmp); err == nil {
					meta["ModifiedDate"] = fmt.Sprintf("%d", t.Unix())
				}
			}
			if tmp, ok := meta["created"]; ok {
				if t, err := time.Parse(time.RFC3339, tmp); err == nil {
					meta["CreatedDate"] = fmt.Sprintf("%d", t.Unix())
				}
			}
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml":
			body, err := parseDocxText(f)
			if err != nil {
				return "", err
			}
			textBody += body + "\n"
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml":
			footer, err := parseDocxText(f)
			if err != nil {
				return "", err
			}
			textFooter += footer + "\n"
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml":
			header, err := parseDocxText(f)
			if err != nil {
				return "", err
			}
			textHeader += header + "\n"
		}

	}
	// 在成功解析ZIP文件后，添加图片提取逻辑
	images, err := findImagesInZip(zr)
	if err != nil {
		fmt.Printf("Error extracting images: %v", err)
	}
	fmt.Printf("Images: %v", images)

	return textHeader + "\n" + textBody + "\n" + textFooter, nil
}

func getContentTypeDefinition(zf *zip.File) (*contentTypeDefinition, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x := &contentTypeDefinition{}
	if err := xml.NewDecoder(io.LimitReader(f, maxBytes)).Decode(x); err != nil {
		return nil, err
	}
	return x, nil
}

func mapZipFiles(files []*zip.File) map[string]*zip.File {
	filesMap := make(map[string]*zip.File, 2*len(files))
	for _, f := range files {
		filesMap[f.Name] = f
		filesMap["/"+f.Name] = f
	}
	return filesMap
}

func parseDocxText(f *zip.File) (string, error) {
	r, err := f.Open()
	if err != nil {
		return "", fmt.Errorf("error opening '%v' from archive: %v", f.Name, err)
	}
	defer r.Close()

	text, err := DocxXMLToText(r)
	if err != nil {
		return "", fmt.Errorf("error parsing '%v': %v", f.Name, err)
	}
	return text, nil
}

// DocxXMLToText converts Docx XML into plain text.
func DocxXMLToText(r io.Reader) (string, error) {
	return XMLToText(r, []string{"br", "p", "tab"}, []string{"instrText", "script"}, true)
}
