package office

import (
	"bytes"
	"fmt"
	pdf "godo/office/pdf"
	xlsx "godo/office/xlsx"
	"html"
	"os"
	"regexp"
	"strings"
)

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

func removeStrangeChars(input string) string {
	// Define the regex pattern for allowed characters
	re := regexp.MustCompile("[�\x13\x0b]+")
	// Replace all disallowed characters with an empty string
	return re.ReplaceAllString(input, " ")
}
