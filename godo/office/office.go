/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/

package office

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	pdf "godo/office/pdf"
	xlsx "godo/office/xlsx"
	"html"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func GetDocument(pathname string) (*Document, error) {
	abPath, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}
	filename := path.Base(pathname)
	data := Document{path: pathname, RePath: abPath, Title: filename}
	extension := path.Ext(pathname)
	_, err = getFileInfoData(&data)
	if err != nil {
		return &data, err
	}
	switch extension {
	case ".docx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, docx2txt)
	case ".pptx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, pptx2txt)
	case ".xlsx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, xlsx2txt)
	case ".pdf":
		_, err = getContentData(&data, pdf2txt)
	case ".doc":
		_, err = getContentData(&data, doc2txt)
	case ".ppt":
		_, err = getContentData(&data, ppt2txt)
	case ".xls":
		_, err = getContentData(&data, xls2txt)
	case ".epub":
		_, err = getContentData(&data, epub2txt)
	case ".odt":
		_, err = getContentData(&data, odt2txt)
	case ".xml":
		_, err = getContentData(&data, xml2txt)
	case ".rtf":
		_, err = getContentData(&data, rtf2txt)
	case ".md":
		_, err = getContentData(&data, md2txt)
	case ".txt":
		_, err = getContentData(&data, text2txt)
	case ".xhtml", ".html", ".htm":
		_, err = getContentData(&data, html2txt)
	case ".json":
		_, err = getContentData(&data, json2txt)
	}
	if err != nil {
		return &data, err
	}
	return &data, nil
}

// Read the meta data of office files (only *.docx, *.xlsx, *.pptx) and insert into the interface
func getMetaData(data *Document) (bool, error) {
	file, err := os.Open(data.path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	meta, err := GetContent(file)
	if err != nil {
		return false, errors.New("failed to get office meta data")
	}
	if meta.Title != "" {
		data.Title = meta.Title
	}
	data.Subject = meta.Subject
	data.Creator = meta.Creator
	data.Keywords = meta.Keywords
	data.Description = meta.Description
	data.Lastmodifiedby = meta.LastModifiedBy
	data.Revision = meta.Revision
	data.Category = meta.Category
	data.Content = meta.Category
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

// Read the content of office files and insert into the interface
func getContentData(data *Document, reader DocReader) (bool, error) {
	content, err := reader(data.path)
	if err != nil {
		return false, err
	}
	data.Content = content
	return true, nil
}

// Read the file information of any files and insert into the interface

func removeStrangeChars(input string) string {
	// Define the regex pattern for allowed characters
	re := regexp.MustCompile("[�\x13\x0b]+")
	// Replace all disallowed characters with an empty string
	return re.ReplaceAllString(input, " ")
}

func docx2txt(filename string) (string, error) {
	data_docx, err := ReadDocxFile(filename) // Read data from docx file
	if err != nil {
		return "", err
	}
	defer data_docx.Close()
	text_docx := data_docx.Editable().GetContent()        // Get whole docx data as XML formated text
	text_docx = PARA_RE.ReplaceAllString(text_docx, "\n") // Replace the end of paragraphs (</w:p) with /n
	text_docx = TAG_RE.ReplaceAllString(text_docx, "")    // Remove all the tags to extract the content
	text_docx = html.UnescapeString(text_docx)            // Replace all the html entities (e.g. &amp)

	// fmt.Println(text_docx)
	return text_docx, nil
}

func pptx2txt(filename string) (string, error) {
	data_pptx, err := ReadPowerPoint(filename) // Read data from pptx file
	if err != nil {
		return "", err
	}

	data_pptx.DeletePassWord()
	slides_pptx := data_pptx.GetSlidesContent() // Get pptx slides data as an array of XML formated text
	var text_pptx string
	for i := range slides_pptx {
		slide_text_pptx := PARA_RE.ReplaceAllString(slides_pptx[i], "\n") // Replace the end of paragraphs (</w:p) with /n
		slide_text_pptx = TAG_RE.ReplaceAllString(slide_text_pptx, "")    // Remove all the tags to extract the content
		slide_text_pptx = html.UnescapeString(slide_text_pptx)            // Replace all the html entities (e.g. &amp)
		if slide_text_pptx != "" {                                        // Save all slides as ONE string
			if text_pptx != "" {
				text_pptx = fmt.Sprintf("%s\n%s", text_pptx, slide_text_pptx)
			} else {
				text_pptx = fmt.Sprintf("%s%s", text_pptx, slide_text_pptx)
			}
		}
	}
	// fmt.Println(text_pptx)
	return text_pptx, nil
}

func xlsx2txt(filename string) (string, error) {
	data_xlsx, err := xlsx.OpenFile(filename) // Read data from xlsx file
	if err != nil {
		return "", err
	}
	defer data_xlsx.Close()

	var rows_xlsx string
	for _, sheet := range data_xlsx.Sheets { // For each sheet of the file
		for row := range data_xlsx.ReadRows(sheet) { // For each row of the sheet
			text_row := ""
			for i, col := range row.Cells { // Concatenate cells of the row with tab separator
				if i > 0 {
					text_row = fmt.Sprintf("%s\t%s", text_row, col.Value)
				} else {
					text_row = fmt.Sprintf("%s%s", text_row, col.Value)
				}
			}
			if rows_xlsx != "" { // Save all rows as ONE string
				rows_xlsx = fmt.Sprintf("%s\n%s", rows_xlsx, text_row)
			} else {
				rows_xlsx = fmt.Sprintf("%s%s", rows_xlsx, text_row)
			}
		}
	}
	// fmt.Println(rows_xlsx)
	return rows_xlsx, nil
}

func pdf2txt(filename string) (string, error) { // BUG: Cannot get text from specific (or really malformed?) pages
	file_pdf, data_pdf, err := pdf.Open(filename) // Read data from pdf file
	if err != nil {
		return "", err
	}
	defer file_pdf.Close()

	var buff_pdf bytes.Buffer
	bytes_pdf, err := data_pdf.GetPlainText() // Get text of entire pdf file
	if err != nil {
		return "", err
	}

	buff_pdf.ReadFrom(bytes_pdf)
	text_pdf := buff_pdf.String()
	// fmt.Println(text_pdf)
	return text_pdf, nil
}

func doc2txt(filename string) (string, error) {
	file_doc, _ := os.Open(filename)    // Open doc file
	data_doc, err := DOC2Text(file_doc) // Read data from a doc file
	if err != nil {
		return "", err
	}
	defer file_doc.Close()

	actual := data_doc.(*bytes.Buffer) // Buffer for hold line text of doc file
	text_doc := ""
	for aline, err := actual.ReadString('\r'); err == nil; aline, err = actual.ReadString('\r') { // Get text by line
		aline = strings.Trim(aline, " \n\r")
		if aline != "" {
			if text_doc != "" {
				text_doc = fmt.Sprintf("%s\n%s", text_doc, removeStrangeChars(aline))
			} else {
				text_doc = fmt.Sprintf("%s%s", text_doc, removeStrangeChars(aline))
			}
		}
	}
	text_doc = removeStrangeChars(text_doc)
	// fmt.Println(text_doc)
	return text_doc, nil
}

func ppt2txt(filename string) (string, error) {
	file_ppt, err := os.Open(filename) // Open ppt file
	if err != nil {
		return "", err
	}
	defer file_ppt.Close()

	text_ppt, err := ExtractPPTText(file_ppt) // Read text from a ppt file
	if err != nil {
		return "", err
	}
	text_ppt = removeStrangeChars(text_ppt)
	// fmt.Println(text_ppt)
	return text_ppt, nil
}

func xls2txt(filename string) (string, error) {
	file_xls, err := os.Open(filename) // Open xls file
	if err != nil {
		return "", err
	}
	defer file_xls.Close()

	text_xls, err := XLS2Text(file_xls) // Convert xls data to an array of rows (include all sheets)
	if err != nil {
		return "", err
	}
	text_xls = removeStrangeChars(text_xls)
	// fmt.Println(text_xls)
	return text_xls, nil
}
