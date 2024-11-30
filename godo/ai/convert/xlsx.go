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
	"strings"

	"godo/ai/convert/libs"

	"github.com/pbnjay/grate"
	_ "github.com/pbnjay/grate/simple" // tsv and csv support
	_ "github.com/pbnjay/grate/xls"
	_ "github.com/pbnjay/grate/xlsx"
)

// 返回行索引 列索引
func ConvertXlsx(r io.Reader) (string, error) {
	absFileFrom, tmpfromfile, err := libs.GetTempFile(r, "prefix-xlsx-from")
	if err != nil {
		return "", err
	}
	textByRow := ""
	textByColumn := ""

	wb, _ := grate.Open(absFileFrom) // open the file
	sheets, _ := wb.List()           // list available sheets

	// 用于存储每一列的内容
	columns := make([][]string, 0)

	for _, s := range sheets { // enumerate each sheet name
		sheet, _ := wb.Get(s) // open the sheet
		maxColumns := 0
		for sheet.Next() { // enumerate each row of data
			row := sheet.Strings() // get the row's content as []string

			// 更新最大列数
			if len(row) > maxColumns {
				maxColumns = len(row)
			}

			// 跳过空记录
			if len(row) == 0 {
				continue
			}

			textByRow += strings.Join(row, "\t") + "\n"

			// 初始化列切片
			if len(columns) < maxColumns {
				columns = make([][]string, maxColumns)
			}

			// 将每一列的内容添加到对应的列切片中
			for i, cell := range row {
				columns[i] = append(columns[i], cell)
			}
		}
	}

	// 拼接每一列的内容
	for _, col := range columns {
		textByColumn += strings.Join(col, "\n") + "\n"
	}

	wb.Close()
	libs.CloseTempFile(tmpfromfile)
	return textByRow + "\n\n" + textByColumn, nil
}
