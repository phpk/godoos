package office

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"godo/office/xls"
)

func XLS2Text(reader io.ReadSeeker) (string, error) {

	xlFile, err := xls.OpenReader(reader, "utf-8")
	if err != nil || xlFile == nil {
		return "", err
	}

	extracted_text := ""
	for n := 0; n < xlFile.NumSheets(); n++ {
		if sheet1 := xlFile.GetSheet(n); sheet1 != nil {
			if extracted_text != "" {
				extracted_text = fmt.Sprintf("%s\n%s", extracted_text, xlGenerateSheetTitle(sheet1.Name, n, int(sheet1.MaxRow)))
			} else {
				extracted_text = fmt.Sprintf("%s%s", extracted_text, xlGenerateSheetTitle(sheet1.Name, n, int(sheet1.MaxRow)))
			}

			for m := 0; m <= int(sheet1.MaxRow); m++ {
				row1 := sheet1.Row(m)
				if row1 == nil {
					continue
				}

				rowText := ""

				// go through all columns
				for c := row1.FirstCol(); c < row1.LastCol(); c++ {
					if text := row1.Col(c); text != "" {
						text = cleanCell(text)

						if c > row1.FirstCol() {
							rowText += ", "
						}
						rowText += text
					}
				}
				if extracted_text != "" {
					extracted_text = fmt.Sprintf("%s\n%s", extracted_text, rowText)
				} else {
					extracted_text = fmt.Sprintf("%s%s", extracted_text, rowText)
				}
			}
		}
	}

	return extracted_text, nil
}

// cleanCell returns a cleaned cell text without new-lines
func cleanCell(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.TrimSpace(text)

	return text
}

func xlGenerateSheetTitle(name string, number, rows int) (title string) {
	if number > 0 {
		title += "\n"
	}

	title += fmt.Sprintf("Sheet \"%s\" (%d rows):\n", name, rows)

	return title
}

// func writeOutput(writer io.Writer, output []byte, alreadyWritten *int64, size *int64) (err error) {

// 	if int64(len(output)) > *size {
// 		output = output[:*size]
// 	}

// 	*size -= int64(len(output))

// 	writtenOut, err := writer.Write(output)
// 	*alreadyWritten += int64(writtenOut)

// 	return err
// }

// IsFileXLS checks if the data indicates a XLS file
// XLS has a signature of D0 CF 11 E0 A1 B1 1A E1
func IsFileXLS(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
}

// XLS2Cells converts an XLS file to individual cells
func XLS2Cells(reader io.ReadSeeker) (cells []string, err error) {

	xlFile, err := xls.OpenReader(reader, "utf-8")
	if err != nil || xlFile == nil {
		return nil, err
	}

	for n := 0; n < xlFile.NumSheets(); n++ {
		if sheet1 := xlFile.GetSheet(n); sheet1 != nil {
			for m := 0; m <= int(sheet1.MaxRow); m++ {
				row1 := sheet1.Row(m)
				if row1 == nil {
					continue
				}

				for c := row1.FirstCol(); c < row1.LastCol(); c++ {
					if text := row1.Col(c); text != "" {
						text = cleanCell(text)
						cells = append(cells, text)
					}
				}
			}
		}
	}

	return
}
