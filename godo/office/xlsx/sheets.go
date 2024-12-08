package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"strings"
)

// workbook is a struct representing the data we care about from the workbook.xml file.
type workbook struct {
	Sheets []sheet `xml:"sheets>sheet"`
}

// sheet is a struct representing the sheet xml element.
type sheet struct {
	Name           string `xml:"name,attr,omitempty"`
	RelationshipID string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr,omitempty"`
}

// relationships is a struct representing the data we care about from the _rels/workboox.xml.rels file.
type relationships struct {
	Relationships []relationship `xml:"Relationship"`
}

type relationship struct {
	ID     string `xml:"Id,attr,omitempty"`
	Target string `xml:"Target,attr,omitempty"`
}

func getFileNameFromRelationships(rels []relationship, s sheet) (string, error) {
	for _, rel := range rels {
		if rel.ID == s.RelationshipID {
			if strings.HasPrefix(rel.Target, "/") {
				// path is absolute, take all but the leading slash
				return rel.Target[1:], nil
			}
			// path is relative, so needs xl/ adding
			return "xl/" + rel.Target, nil
		}
	}
	return "", fmt.Errorf("unable to find file with relationship %s", s.RelationshipID)
}

// getWorksheets loads the workbook.xml file and extracts a list of worksheets, along
// with a map of the canonical worksheet name to a file descriptor.
// This will return an error if it is not possible to read the workbook.xml file, or
// if a worksheet without a file is referenced.
func getWorksheets(files []*zip.File) ([]string, *map[string]*zip.File, error) {
	wbFile, err := getFileForName(files, "xl/workbook.xml")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get workbook file: %w", err)
	}
	data, err := readFile(wbFile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read workbook file: %w", err)
	}

	var wb workbook
	err = xml.Unmarshal(data, &wb)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse workbook file: %w", err)
	}

	relsFile, err := getFileForName(files, "xl/_rels/workbook.xml.rels")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get relationships file: %w", err)
	}
	relsData, err := readFile(relsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read relationships file: %w", err)
	}

	rels := relationships{}
	err = xml.Unmarshal(relsData, &rels)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse relationships file: %w", err)
	}

	wsFileMap := map[string]*zip.File{}
	sheetNames := make([]string, len(wb.Sheets))

	for i, sheet := range wb.Sheets {
		sheetFilename, err := getFileNameFromRelationships(rels.Relationships, sheet)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to get file name from relationships: %w", err)
		}
		sheetFile, err := getFileForName(files, sheetFilename)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to get file for sheet name %s: %w", sheetFilename, err)
		}

		wsFileMap[sheet.Name] = sheetFile
		sheetNames[i] = sheet.Name
	}

	return sheetNames, &wsFileMap, nil
}
