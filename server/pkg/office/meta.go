package office

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"errors"
	"os"
)

func getMetaData(data *Document, ext string) (bool, error) {
	file, err := os.Open(data.path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	switch ext {
	case ".docx", ".xlsx", ".pptx":
		meta, err := GetContent(file)
		if err != nil {
			return true, errors.New("failed to get office meta data")
		}
		data.Title = meta.Title
		data.Subject = meta.Subject
		data.Creator = meta.Creator
		data.Keywords = meta.Keywords
		data.Description = meta.Description
		data.Lastmodifiedby = meta.LastModifiedBy
		data.Revision = meta.Revision
		data.Category = meta.Category
	default:
		return true, nil
	}

	return true, nil
}

func GetContent(document *os.File) (fields XMLContent, err error) {
	// Attempt to read the document file directly as a zip file.
	z, err := zip.OpenReader(document.Name())
	if err != nil {
		return fields, errors.New("failed to open the file as zip")
	}
	defer z.Close()

	var xmlFile string
	for _, file := range z.File {
		if file.Name == "docProps/core.xml" {
			rc, err := file.Open()
			if err != nil {
				return fields, errors.New("failed to open docProps/core.xml")
			}
			defer rc.Close()

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				xmlFile += scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				return fields, errors.New("failed to read from docProps/core.xml")
			}
			break // Exit loop after finding and reading core.xml
		}
	}

	// Unmarshal the collected XML content into the XMLContent struct
	if err := xml.Unmarshal([]byte(xmlFile), &fields); err != nil {
		return fields, errors.New("failed to Unmarshal")
	}

	return fields, nil
}
