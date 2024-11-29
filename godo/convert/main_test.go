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
package convert_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"godo/convert"
)

func TestConvert(t *testing.T) {
	tempDir := "../testdata"
	cases := []struct {
		name        string
		filename    string
		expectedRes convert.Res
	}{
		{
			name:        "HTTP",
			filename:    "https://www.baidu.com",
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "docx",
			filename:    filepath.Join(tempDir, "test.docx"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "xls",
			filename:    filepath.Join(tempDir, "test.xls"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "pdf",
			filename:    filepath.Join(tempDir, "test.pdf"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "pptx",
			filename:    filepath.Join(tempDir, "test.pptx"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "rtf",
			filename:    filepath.Join(tempDir, "test.rtf"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "odt",
			filename:    filepath.Join(tempDir, "test.odt"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "txt",
			filename:    filepath.Join(tempDir, "test.txt"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "md",
			filename:    filepath.Join(tempDir, "test.md"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "html",
			filename:    filepath.Join(tempDir, "test.html"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "jpg",
			filename:    filepath.Join(tempDir, "test.jpg"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "xml",
			filename:    filepath.Join(tempDir, "test.xml"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
		{
			name:        "epub",
			filename:    filepath.Join(tempDir, "test.epub"),
			expectedRes: convert.Res{Status: 0, Data: ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			// 调用 Convert 函数并检查结果
			res := convert.Convert(tc.filename)
			fmt.Printf("res: %v\n", tc.filename)
			// 比较结果
			if res.Status != tc.expectedRes.Status {
				t.Errorf("For case '%s', expected status %d, got %d", tc.name, tc.expectedRes.Status, res.Status)
			}
			// 如果需要，也可以比较 Data 字段
			// 注意：根据实际情况调整比较逻辑，此处省略了对 Data 的直接比较
		})
	}
}
