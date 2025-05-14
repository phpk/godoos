package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// sharedStringsValue is a struct that holds the value of the shared strings.
type sharedStringsValue struct {
	Text     string   `xml:"t"`
	RichText []string `xml:"r>t"`
}

func (sv *sharedStringsValue) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("error retrieving xml token: %w", err)
		}

		var se xml.StartElement

		switch el := tok.(type) {
		case xml.EndElement:
			if el == start.End() {
				return nil
			}
			continue
		case xml.StartElement:
			se = el
		default:
			continue
		}

		switch se.Name.Local {
		case "t":
			sv.Text, err = getCharData(d)
		case "r":
			err = sv.decodeRichText(d, se)
		default:
			continue
		}

		if err != nil {
			return fmt.Errorf("unable to parse string: %w", err)
		}
	}
}

func (sv *sharedStringsValue) decodeRichText(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("unable to get shared strings value token: %w", err)
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

		var s string

		if s, err = getCharData(d); err != nil {
			return fmt.Errorf("unable to parse string: %w", err)
		}

		sv.RichText = append(sv.RichText, s)
	}
}

// String gets a string value from the raw sharedStringsValue struct.
// Since the values can appear in many different places in the xml structure, we need to normalise this.
// They can either be:
// <si> <t> value </t> </si>
// or
// <si> <r> <t> val </t> </r> <r> <t> ue </t> </r> </si>
func (sv *sharedStringsValue) String() string {
	// fast path: no rich text, just return text
	if len(sv.RichText) == 0 {
		return sv.Text
	}

	var sb strings.Builder
	for _, t := range sv.RichText {
		sb.WriteString(t)
	}

	return sb.String()
}

// Reset zeroes data inside struct.
func (sv *sharedStringsValue) Reset() {
	sv.Text = ""
	sv.RichText = sv.RichText[:0]
}

// Sentinel error to indicate that no shared strings file can be found
var errNoSharedStrings = errors.New("no shared strings file exists")

// getSharedStringsFile attempts to find and return the zip.File struct associated with the
// shared strings section of an xlsx file. An error is returned if the sharedStrings file
// does not exist, or cannot be found.
func getSharedStringsFile(files []*zip.File) (*zip.File, error) {
	for _, file := range files {
		if file.Name == "xl/sharedStrings.xml" || file.Name == "xl/SharedStrings.xml" {
			return file, nil
		}
	}

	return nil, errNoSharedStrings
}

// getSharedStrings loads the contents of the shared string file into memory.
// This serves as a large lookup table of values, so we can efficiently parse rows.
func getSharedStrings(files []*zip.File) ([]string, error) {
	ssFile, err := getSharedStringsFile(files)
	if err != nil && errors.Is(err, errNoSharedStrings) {
		// Valid to contain no shared strings
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to get shared strings file: %w", err)
	}

	f, err := ssFile.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open shared strings file: %w", err)
	}

	defer f.Close()

	var (
		sharedStrings []string
		value         sharedStringsValue
	)

	dec := xml.NewDecoder(f)
	for {
		token, err := dec.Token()
		if err == io.EOF {
			return sharedStrings, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error decoding token: %w", err)
		}

		startElement, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		if sharedStrings == nil { // don't use len() == 0 here!
			sharedStrings = makeSharedStringsSlice(startElement)
			continue
		}

		value.Reset()
		if err := value.unmarshalXML(dec, startElement); err != nil {
			return nil, fmt.Errorf("error unmarshaling shared strings value %+v: %w", startElement, err)
		}

		sharedStrings = append(sharedStrings, value.String())
	}
}

// makeSharedStringsSlice allocates shared strings slice according to 'count' attribute of root tag
// absence of attribute doesn't break flow because make(..., 0) is valid
func makeSharedStringsSlice(rootElem xml.StartElement) []string {
	var count int
	for _, attr := range rootElem.Attr {
		if attr.Name.Local != "count" {
			continue
		}

		var err error

		count, err = strconv.Atoi(attr.Value)
		if err != nil {
			return []string{}
		}
	}

	return make([]string, 0, count)
}
