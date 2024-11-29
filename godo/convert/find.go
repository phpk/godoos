package convert

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"godo/libs"
)

// 常见图片扩展名列表
var imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tif", ".tiff"}

// 直接遍历ZIP并查找可能的图片文件
func findImagesInZip(zr *zip.Reader) ([]string, error) {
	var images []string
	cacheDir := libs.GetCacheDir()

	for _, f := range zr.File {
		if isImageFile(f.Name) {
			images = append(images, f.Name)
			if err := extractImageToCache(zr, f.Name, cacheDir); err != nil {
				log.Printf("Error extracting image %s to cache: %v", f.Name, err)
			}
		}
	}

	return images, nil
}

// 判断文件是否为图片文件
func isImageFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// 从zip中提取图片到缓存目录
func extractImageToCache(zr *zip.Reader, imageName, cacheDir string) error {
	fileInZip, err := getFileByName(zr.File, imageName)
	if err != nil {
		return err
	}

	rc, err := fileInZip.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s in zip: %w", imageName, err)
	}
	defer rc.Close()

	justFileName := filepath.Base(imageName)
	outFilePath := filepath.Join(cacheDir, justFileName)

	outFile, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outFilePath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return fmt.Errorf("failed to copy image content to %s: %w", outFilePath, err)
	}

	// 提取图片周围的文本内容
	textContent, err := getSurroundingTextForOffice(zr, imageName)
	if err != nil {
		log.Printf("Error getting surrounding text for image %s: %v", imageName, err)
	} else {
		textFilePath := filepath.Join(cacheDir, strings.TrimSuffix(justFileName, filepath.Ext(justFileName))+".txt")
		if err := saveTextToFile(textContent, textFilePath); err != nil {
			log.Printf("Error saving text to file %s: %v", textFilePath, err)
		}
	}

	return nil
}

// 从zip.File数组中根据文件名查找并返回对应的文件
func getFileByName(files []*zip.File, name string) (*zip.File, error) {
	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, fmt.Errorf("file %s not found in zip archive", name)
}

// 获取 .pptx, .xlsx 或 .docx 文件中图片周围的文本内容
func getSurroundingTextForOffice(zr *zip.Reader, imageName string) (string, error) {
	imageDir := filepath.Dir(imageName)
	xmlFiles, err := findRelevantXMLFiles(zr, imageDir)
	if err != nil {
		return "", err
	}

	for _, xmlFile := range xmlFiles {
		fileInZip, err := getFileByName(zr.File, xmlFile)
		if err != nil {
			continue
		}

		rc, err := fileInZip.Open()
		if err != nil {
			continue
		}
		defer rc.Close()

		doc, err := parseXMLDocument(rc, imageDir)
		if err != nil {
			continue
		}

		surroundingText := getSurroundingText(doc, filepath.Base(imageName))
		if surroundingText != "" {
			return truncateText(surroundingText), nil
		}
	}

	return "", fmt.Errorf("no surrounding text found for image %s", imageName)
}

// 查找相关的XML文件
func findRelevantXMLFiles(zr *zip.Reader, imageDir string) ([]string, error) {
	switch {
	case strings.Contains(imageDir, "ppt/media"):
		return findFilesByPattern(zr, "ppt/slides/slide*.xml"), nil
	case strings.Contains(imageDir, "xl/media"):
		return findFilesByPattern(zr, "xl/worksheets/sheet*.xml"), nil
	case strings.Contains(imageDir, "word/media"):
		return []string{"word/document.xml"}, nil
	default:
		return nil, fmt.Errorf("unknown image directory %s", imageDir)
	}
}

// 解析XML文档
func parseXMLDocument(rc io.ReadCloser, imageDir string) (interface{}, error) {
	var doc interface{}
	switch {
	case strings.Contains(imageDir, "ppt/media"):
		doc = &PPTXDocument{}
	case strings.Contains(imageDir, "xl/media"):
		doc = &XLSXDocument{}
	case strings.Contains(imageDir, "word/media"):
		doc = &DOCXDocument{}
	default:
		return nil, fmt.Errorf("unknown image directory %s", imageDir)
	}

	if err := xml.NewDecoder(rc).Decode(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// 获取图片周围的文本内容
func getSurroundingText(doc interface{}, imagePath string) string {
	switch d := doc.(type) {
	case *PPTXDocument:
		for _, slide := range d.Slides {
			for _, shape := range slide.Shapes {
				if shape.Type == "pic" && shape.ImagePath == imagePath {
					return getTextFromSlide(slide)
				}
			}
		}
	case *XLSXDocument:
		for _, sheet := range d.Sheets {
			for _, drawing := range sheet.Drawings {
				for _, image := range drawing.Images {
					if image.ImagePath == imagePath {
						return getTextFromSheet(sheet)
					}
				}
			}
		}
	case *DOCXDocument:
		for _, paragraph := range d.Body.Paragraphs {
			for _, run := range paragraph.Runs {
				for _, pic := range run.Pictures {
					if pic.ImagePath == imagePath {
						return getTextFromParagraph(paragraph)
					}
				}
			}
		}
	}
	return ""
}

// 查找符合模式的文件
func findFilesByPattern(zr *zip.Reader, pattern string) []string {
	var files []string
	for _, f := range zr.File {
		if matched, _ := filepath.Match(pattern, f.Name); matched {
			files = append(files, f.Name)
		}
	}
	return files
}

// 将文本内容保存到文件
func saveTextToFile(text, filePath string) error {
	return os.WriteFile(filePath, []byte(text), 0644)
}

// 截断文本，确保不超过80个字符
func truncateText(text string) string {
	if len(text) > 80 {
		return text[:80]
	}
	return text
}

// PPTXDocument 结构体定义
type PPTXDocument struct {
	Slides []Slide `xml:"p:sld"`
}

type Slide struct {
	Shapes []Shape `xml:"p:cSld>p:spTree>p:sp"`
}

type Shape struct {
	Type      string    `xml:"p:pic"`
	ImagePath string    `xml:"p:pic>p:blipFill>a:blip/@r:embed"`
	Elements  []Element `xml:"p:txBody>a:p>a:r"`
}

type Element struct {
	Type  string `xml:"a:t"`
	Value string `xml:",chardata"`
}

// XLSXDocument 结构体定义
type XLSXDocument struct {
	Sheets []Sheet `xml:"worksheet"`
}

type Sheet struct {
	Rows     []Row     `xml:"sheetData>row"`
	Drawings []Drawing `xml:"drawing"`
}

type Row struct {
	Cells []Cell `xml:"c"`
}

type Cell struct {
	Value string `xml:"v"`
}

type Drawing struct {
	Images []Image `xml:"xdr:pic"`
}

type Image struct {
	ImagePath string `xml:"xdr:pic>xdr:blipFill>a:blip/@r:embed"`
}

// DOCXDocument 结构体定义
type DOCXDocument struct {
	Body struct {
		Paragraphs []Paragraph `xml:"w:p"`
	} `xml:"w:body"`
}

type Paragraph struct {
	Runs []Run `xml:"w:r"`
}

type Run struct {
	Pictures []Picture `xml:"w:drawing"`
	Text     []Text    `xml:"w:t"`
}

type Text struct {
	Value string `xml:",chardata"`
}

type Picture struct {
	ImagePath string `xml:"wp:docPr/@name"`
}

// 从幻灯片中提取文本
func getTextFromSlide(slide Slide) string {
	var text string
	for _, shape := range slide.Shapes {
		if shape.Type != "pic" {
			text += getTextFromShape(shape)
		}
	}
	return text
}

// 从形状中提取文本
func getTextFromShape(shape Shape) string {
	var text string
	for _, element := range shape.Elements {
		text += element.Value
	}
	return text
}

// 从工作表中提取文本
func getTextFromSheet(sheet Sheet) string {
	var text string
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			text += cell.Value
		}
	}
	return text
}

// 从段落中提取文本
func getTextFromParagraph(paragraph Paragraph) string {
	var text string
	for _, run := range paragraph.Runs {
		for _, t := range run.Text {
			text += t.Value
		}
	}
	return text
}
