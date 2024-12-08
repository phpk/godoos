package xlsx

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// rawRow represent the raw XML element for parsing a row of data.
type rawRow struct {
	Index    int       `xml:"r,attr,omitempty"`
	RawCells []rawCell `xml:"c"`
}

func (rr *rawRow) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local != "r" {
			continue
		}

		var err error

		if rr.Index, err = strconv.Atoi(attr.Value); err != nil {
			return fmt.Errorf("unable to parse row index: %w", err)
		}
	}

	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("error retrieving xml token: %w", err)
		}

		var se xml.StartElement

		switch el := tok.(type) {
		case xml.StartElement:
			se = el
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
		default:
			continue
		}

		if se.Name.Local != "c" {
			continue
		}

		var rc rawCell
		if err = rc.unmarshalXML(d, se); err != nil {
			return fmt.Errorf("unable to unmarshal cell: %w", err)
		}

		rr.RawCells = append(rr.RawCells, rc)
	}
}

// rawCell represents the raw XML element for parsing a cell.
type rawCell struct {
	Reference    string  `xml:"r,attr"` // E.g. A1
	Type         string  `xml:"t,attr,omitempty"`
	Value        *string `xml:"v,omitempty"`
	Style        int     `xml:"s,attr"`
	InlineString *string `xml:"is>t"`
}

func (rc *rawCell) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// unmarshal attributes
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "r":
			rc.Reference = attr.Value
		case "t":
			rc.Type = attr.Value
		case "s":
			var err error

			if rc.Style, err = strconv.Atoi(attr.Value); err != nil {
				return err
			}
		}
	}

	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("error retrieving xml token: %w", err)
		}

		var se xml.StartElement

		switch el := tok.(type) {
		case xml.StartElement:
			se = el
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
			continue
		default:
			continue
		}

		switch se.Name.Local {
		case "is":
			err = rc.unmarshalInlineString(d, se)
		case "v":
			var v string

			if v, err = getCharData(d); err != nil {
				return err
			}

			rc.Value = &v
		default:
			continue
		}

		if err != nil {
			return fmt.Errorf("unable to parse cell data: %w", err)
		}
	}
}

func (rc *rawCell) unmarshalInlineString(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("error retrieving xml token: %w", err)
		}

		var se xml.StartElement

		switch el := tok.(type) {
		case xml.StartElement:
			se = el
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
			continue
		default:
			continue
		}

		if se.Name.Local != "t" {
			continue
		}

		v, err := getCharData(d)
		if err != nil {
			return fmt.Errorf("unable to parse string: %w", err)
		}

		rc.InlineString = &v
		return nil
	}
}

// Row represents a row of data read from an Xlsx file, in a consumable format
type Row struct {
	Error error
	Index int
	Cells []Cell
}

// Cell represents the data in a single cell as a consumable format.
type Cell struct {
	Column string // E.G   A, B, C
	Row    int
	Value  string
	Type   CellType
}

// CellType defines the data type of an excel cell
type CellType string

const (
	// TypeString is for text cells
	TypeString CellType = "string"
	// TypeNumerical is for numerical values
	TypeNumerical CellType = "numerical"
	// TypeDateTime is for date values
	TypeDateTime CellType = "datetime"
	// TypeBoolean is for true/false values
	TypeBoolean CellType = "boolean"
)

// ColumnIndex gives a number, representing the column the cell lies beneath.
func (c Cell) ColumnIndex() int {
	return asIndex(c.Column)
}

// getCellValue interrogates a raw cell to get a textual representation of the cell's contents.
// Numerical values are returned in their string format.
// Dates are returned as an ISO YYYY-MM-DD formatted string.
// Datetimes are returned in RFC3339 (ISO-8601) YYYY-MM-DDTHH:MM:SSZ formatted string.
func (x *XlsxFile) getCellValue(r rawCell) (string, error) {
	if r.Type == "inlineStr" {
		if r.InlineString == nil {
			return "", fmt.Errorf("cell had type of InlineString, but the InlineString attribute was missing")
		}
		return *r.InlineString, nil
	}

	if r.Value == nil {
		return "", fmt.Errorf("unable to get cell value for cell %s - no value element found", r.Reference)
	}

	if r.Type == "s" {
		index, err := strconv.Atoi(*r.Value)
		if err != nil {
			return "", err
		}
		if len(x.sharedStrings) <= index {
			return "", fmt.Errorf("attempted to index value %d in shared strings of length %d",
				index, len(x.sharedStrings))
		}

		return x.sharedStrings[index], nil
	}

	if x.dateStyles[r.Style] && r.Type != "d" {
		formattedDate, err := convertExcelDateToDateString(*r.Value)
		if err != nil {
			return "", err
		}
		return formattedDate, nil
	}

	return *r.Value, nil
}

func (x *XlsxFile) getCellType(r rawCell) CellType {
	if x.dateStyles[r.Style] {
		return TypeDateTime
	}

	switch r.Type {
	case "b":
		return TypeBoolean
	case "d":
		return TypeDateTime
	case "n", "":
		return TypeNumerical
	case "s", "inlineStr":
		return TypeString
	default:
		return TypeString
	}
}

// readSheetRows iterates over "row" elements within a worksheet,
// pushing a parsed Row struct into a channel for each one.
func (x *XlsxFile) readSheetRows(sheet string, ch chan<- Row) {
	defer close(ch)

	xmlFile, err := x.openSheetFile(sheet)
	if err != nil {
		select {
		case <-x.doneCh:
		case ch <- Row{Error: err}:
		}
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	for {
		token, _ := decoder.Token()
		if token == nil {
			return
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "row" {
				row := x.parseRow(decoder, &startElement)
				if len(row.Cells) < 1 && row.Error == nil {
					continue
				}
				select {
				case <-x.doneCh:
					return
				case ch <- row:
				}
			}
		}
	}
}

func (x *XlsxFile) openSheetFile(sheet string) (io.ReadCloser, error) {
	file, ok := x.sheetFiles[sheet]
	if !ok {
		return nil, fmt.Errorf("unable to open sheet %s", sheet)
	}
	return file.Open()
}

// parseRow parses the raw XML of a row element into a consumable Row struct.
// The Row struct returned will contain any errors that occurred either in
// interrogating values, or in parsing the XML.
func (x *XlsxFile) parseRow(decoder *xml.Decoder, startElement *xml.StartElement) Row {
	var r rawRow
	err := r.unmarshalXML(decoder, *startElement)
	if err != nil {
		return Row{
			Error: err,
			Index: r.Index,
		}
	}

	cells, err := x.parseRawCells(r.RawCells, r.Index)
	if err != nil {
		return Row{
			Error: err,
			Index: r.Index,
		}
	}
	return Row{
		Cells: cells,
		Index: r.Index,
	}
}

// parseRawCells converts a slice of structs containing a raw representation of the XML into
// a standardised slice of Cell structs. An error will be returned if it is not possible
// to interpret the value of any of the cells.
func (x *XlsxFile) parseRawCells(rawCells []rawCell, index int) ([]Cell, error) {
	cells := []Cell{}
	for _, rawCell := range rawCells {
		if rawCell.Value == nil && rawCell.InlineString == nil {
			// This cell is empty, so ignore it
			continue
		}
		column := strings.Map(removeNonAlpha, rawCell.Reference)
		val, err := x.getCellValue(rawCell)
		if err != nil {
			return nil, err
		}

		cells = append(cells, Cell{
			Column: column,
			Row:    index,
			Value:  val,
			Type:   x.getCellType(rawCell),
		})
	}

	return cells, nil
}

// ReadRows provides an interface allowing rows from a specific worksheet to be streamed
// from an xlsx file.
// In order to provide a simplistic interface, this method returns a channel that can be
// range-d over.
//
// If you want to read only some of the values, please ensure that the Close() method is
// called after processing the entire file to stop all active goroutines and prevent any
// potential goroutine leaks.
//
// Notes:
// Xlsx sheets may omit cells which are empty, meaning a row may not have continuous cell
// references. This function makes no attempt to fill/pad the missing cells.
func (x *XlsxFile) ReadRows(sheet string) chan Row {
	rowChannel := make(chan Row)
	go x.readSheetRows(sheet, rowChannel)
	return rowChannel
}

// removeNonAlpha is used in combination with strings.Map to remove any non alpha-numeric
// characters from a cell reference, returning just the column name in a consistent uppercase format.
// For example, a11 -> A, AA1 -> AA
func removeNonAlpha(r rune) rune {
	if 'A' <= r && r <= 'Z' {
		return r
	}
	if 'a' <= r && r <= 'z' {
		// make it uppercase
		return r - 32
	}
	// drop the rune
	return -1
}

// cell name to cell index. 'A' -> 0, 'Z' -> 25, 'AA' -> 26
func asIndex(s string) int {
	index := 0
	for _, c := range s {
		index *= 26
		index += int(c) - 'A' + 1
	}
	return index - 1
}
