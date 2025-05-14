package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// styleSheet defines a struct containing the information we care about from the styles.xml file.
type styleSheet struct {
	NumberFormats []numberFormat `xml:"numFmts>numFmt,omitempty"`
	CellStyles    []cellStyle    `xml:"cellXfs>xf,omitempty"`
}

// numberFormat defines a struct containing the format strings for numerical styles.
type numberFormat struct {
	NumberFormatID int    `xml:"numFmtId,attr,omitempty"`
	FormatCode     string `xml:"formatCode,attr,omitempty"`
}

// cellStyle defines a struct containing style information for a cell.
type cellStyle struct {
	NumberFormatID int `xml:"numFmtId,attr"`
}

// getFormatCode returns the format string for a given format ID.
// If the format code is not found, it returns an empty string.
func getFormatCode(ID int, numberFormats []numberFormat) string {
	for _, nf := range numberFormats {
		if nf.NumberFormatID == ID {
			return nf.FormatCode
		}
	}

	return ""
}

var formatGroup = regexp.MustCompile(`\[.+?\]|\\.|".*?"`)

// isDateFormatCode determines whether a format code is for a date.
func isDateFormatCode(formatCode string) bool {
	c := formatGroup.ReplaceAllString(formatCode, "")
	return strings.ContainsAny(c, "dmhysDMHYS")
}

// getDateStylesFromStyleSheet populates a map of all date related styles, based on their
// style sheet index.
func getDateStylesFromStyleSheet(ss *styleSheet) *map[int]bool {
	dateStyles := map[int]bool{}

	for i, style := range ss.CellStyles {
		if 14 <= style.NumberFormatID && style.NumberFormatID <= 22 {
			dateStyles[i] = true
		}
		if 164 <= style.NumberFormatID {
			formatCode := getFormatCode(style.NumberFormatID, ss.NumberFormats)
			if isDateFormatCode(formatCode) {
				dateStyles[i] = true
			}
		}
	}

	return &dateStyles
}

// getDateFormatStyles reads the styles XML, and returns a map of all styles that relate to date
// fields.
// If the styles sheet cannot be found, or cannot be read, then an error is returned.
func getDateFormatStyles(files []*zip.File) (*map[int]bool, error) {
	stylesFile, err := getFileForName(files, "xl/styles.xml")
	if err != nil {
		return nil, fmt.Errorf("unable to get styles file: %w", err)
	}

	data, err := readFile(stylesFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read styles file: %w", err)
	}

	var ss styleSheet
	err = xml.Unmarshal(data, &ss)
	if err != nil {
		return nil, fmt.Errorf("unable to parse styles file: %w", err)
	}

	return getDateStylesFromStyleSheet(&ss), nil
}
