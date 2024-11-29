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
	"fmt"
	"os"
	"path"
	"strings"
)

const maxBytes = 1024 << 20 // 1GB

type Res struct {
	Status int
	Data   string
}

// Convert 函数根据文件类型，将文件内容转换为字符串格式。
// 支持的文件类型包括：.doc, .docx, .odt, .pdf, .csv, .xls, .xlsx, .tsv,
// .pptx, .rtf, .epub, .xml, .xhtml, .html, .htm, .jpg, .jpeg, .jpe, .jfif,
// .jfif-tbnl, .png, .gif, .bmp, .webp, .tif, .tiff, .txt, .md。
// 如果文件以 http 开头，将直接调用 ConvertHttp 函数进行处理。
// 参数：
//
//	filename string - 文件名或文件URL。
//
// 返回值：
//
//	Res - 包含转换结果的状态码和数据。
func Convert(filename string) Res {
	//libs.InitConvertDir()
	// 检查文件名是否以 http 开头，是则调用 ConvertHttp 处理
	if strings.HasPrefix(filename, "http") {
		return ConvertHttp(filename)
	}
	// 尝试打开文件
	r, err := os.Open(filename)
	if err != nil {
		// 打开文件失败，返回错误信息
		return Res{
			Status: 201,
			Data:   fmt.Sprintf("error opening file: %v", err),
		}

	}
	// 确保文件在函数返回前被关闭
	defer r.Close()

	// 获取文件扩展名，并转为小写
	ext := strings.ToLower(path.Ext(filename))

	var body string
	// 根据文件扩展名，调用相应的转换函数
	switch ext {
	case ".doc":
		body, err = ConvertDoc(r)
	case ".docx":
		body, err = ConvertDocx(r)
	case ".odt":
		body, err = ConvertODT(r)
	// .pages 类型文件的处理暂不支持
	// case ".pages":
	// 	return "application/vnd.apple.pages"
	case ".pdf":
		body, err = ConvertPDF(r)
	case ".csv", ".xls", ".xlsx", ".tsv":
		body, err = ConvertXlsx(r)
	case ".pptx":
		body, err = ConvertPptx(r)
	case ".rtf":
		body, err = ConvertRTF(r)
	case ".epub":
		body, err = ConvetEpub(r)
	case ".xml":
		body, err = ConvertXML(r)
	case ".xhtml", ".html", ".htm":
		body, err = ConvertHTML(r)
	case ".jpg", ".jpeg", ".jpe", ".jfif", ".jfif-tbnl", ".png", ".gif", ".bmp", ".webp", ".tif", ".tiff":
		body, err = ConvertImage(r)
	case ".md":
		body, err = ConvertMd(r)
	case ".txt":
		body, err = ConvertTxt(r)
	}

	// 转换过程中若发生错误，返回错误信息
	if err != nil {
		return Res{
			Status: 204,
			Data:   fmt.Sprintf("error opening file: %v", err),
		}
	}
	return Res{
		Status: 0,
		Data:   body,
	}
}
