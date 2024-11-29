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

	"godo/convert/libs"

	"github.com/pbnjay/grate"
	_ "github.com/pbnjay/grate/simple" // tsv and csv support
	_ "github.com/pbnjay/grate/xls"
	_ "github.com/pbnjay/grate/xlsx"
)

func ConvertXlsx(r io.Reader) (string, error) {
	absFileFrom, tmpfromfile, err := libs.GetTempFile(r, "prefix-xlsx-from")
	if err != nil {
		return "", err
	}
	text := ""
	wb, _ := grate.Open(absFileFrom) // open the file
	sheets, _ := wb.List()           // list available sheets
	for _, s := range sheets {       // enumerate each sheet name
		sheet, _ := wb.Get(s) // open the sheet
		for sheet.Next() {    // enumerate each row of data
			row := sheet.Strings() // get the row's content as []string
			//fmt.Println(strings.Join(row, "\t"))
			// 跳过空记录
			if len(row) == 0 {
				continue
			}
			text += strings.Join(row, "\t") + "\n"
		}
	}
	wb.Close()
	libs.CloseTempFile(tmpfromfile)
	return text, nil
}
