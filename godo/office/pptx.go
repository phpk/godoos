package office

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PowerPoint struct {
	Files        []*zip.File
	Slides       map[string]string
	NotesSlides  map[string]string
	Themes       map[string]string
	Images       map[string]string
	Presentation string
}

func ReadPowerPoint(path string) (*PowerPoint, error) {
	var p PowerPoint
	p.Slides = make(map[string]string)
	p.NotesSlides = make(map[string]string)
	p.Themes = make(map[string]string)
	p.Images = make(map[string]string)
	f, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file" + err.Error())
	}
	p.Files = f.File

	for _, file := range p.Files {
		if strings.Contains(file.Name, "ppt/slides/slide") {
			slideOpen, _ := file.Open()
			p.Slides[file.Name] = string(readCloserToByte(slideOpen))
		}
		if strings.Contains(file.Name, "ppt/notesSlides/notesSlide") {
			notesSlideOpen, _ := file.Open()
			p.NotesSlides[file.Name] = string(readCloserToByte(notesSlideOpen))
		}
		if strings.Contains(file.Name, "ppt/theme/theme") {
			themeOpen, _ := file.Open()
			p.Themes[file.Name] = string(readCloserToByte(themeOpen))
		}
		if strings.Contains(file.Name, "ppt/media/image") {
			imageOpen, _ := file.Open()
			p.Images[file.Name] = string(readCloserToByte(imageOpen))
		}
		if strings.Contains(file.Name, "ppt/presentation.xml") {
			presentationOpen, _ := file.Open()
			p.Presentation = string(readCloserToByte(presentationOpen))
		}
	}

	return &p, nil
}

func (p *PowerPoint) GetSlidesContent() []string {
	var slides []string
	for _, slide := range p.Slides {
		slides = append(slides, slide)
	}
	return slides
}

// 只能删除文本编辑密码
func (p *PowerPoint) DeletePassWord() {
	reg := regexp.MustCompile("<p:modifyVerifier(.*?)/>")
	p.Presentation = reg.ReplaceAllString(p.Presentation, "")
}

func (p *PowerPoint) GetSlideCount() int {
	return len(p.Slides)
}

func (p *PowerPoint) GetNotesSlideCount() int {
	return len(p.NotesSlides)
}

func (p *PowerPoint) GetThemeCount() int {
	return len(p.Themes)
}

func (p *PowerPoint) FindSlideString(findString string) []int {
	var nums []int
	reg := regexp.MustCompile(`\d+`)
	for k, v := range p.Slides {
		if strings.Contains(v, findString) {
			num := reg.FindString(k)
			n, _ := strconv.Atoi(num)
			nums = append(nums, n)
		}
	}
	return nums
}

func (p *PowerPoint) DeleteSlide(index int) error {
	if index <= 0 {
		index = len(p.Slides)
	}
	if index > len(p.Slides) {
		return fmt.Errorf("index out of range")
	}

	p.Slides[fmt.Sprintf("ppt/slides/slide%d.xml", index)] = " "
	for {
		if index == len(p.Slides) {
			break
		} else {
			p.Slides[fmt.Sprintf("ppt/slides/slide%d.xml", index)], p.Slides[fmt.Sprintf("ppt/slides/slide%d.xml", index+1)] = p.Slides[fmt.Sprintf("ppt/slides/slide%d.xml", index+1)], p.Slides[fmt.Sprintf("ppt/slides/slide%d.xml", index)]
			index++
		}
	}
	//通过在p.slides、p.notesslides、p.files中删除对应的页面。删除是成功的，但生成的pptx无法打开。目前只能把需要删除的页置空并放到最后一页。
	//delete(p.Slides, fmt.Sprintf("ppt/slides/slide%d.xml", len(p.Slides)))
	//delete(p.NotesSlides, fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", len(p.NotesSlides)))
	//for k, v := range p.Files {
	//	if strings.Contains(v.Name, fmt.Sprintf("ppt/slides/slide%d.xml", len(p.Slides)+1)) {
	//		p.Files = append(p.Files[:k], p.Files[k+1:]...)
	//	}
	//	if strings.Contains(v.Name, fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", len(p.NotesSlides)+1)) {
	//		p.Files= append(p.Files[:k], p.Files[k+1:]...)
	//	}
	//}

	return nil
}

func (p *PowerPoint) ReplaceSlideContent(oldString string, newString string, num int) {
	for k, v := range p.Slides {
		p.Slides[k] = strings.Replace(v, oldString, newString, num)
	}
}

func (p *PowerPoint) ReplaceNotesSlideContent(oldString string, newString string, num int) {
	for k, v := range p.NotesSlides {
		p.NotesSlides[k] = strings.Replace(v, oldString, newString, num)
	}
}

func (p *PowerPoint) ReplaceThemeName(oldString string, newString string, num int) {
	for k, v := range p.Themes {
		p.Themes[k] = strings.Replace(v, oldString, newString, num)
	}
}

func (p *PowerPoint) ReplaceImage(newImagePath string, index int) error {
	if index > len(p.Images) {
		return fmt.Errorf("index out of range")
	}
	newImageOpen, _ := os.ReadFile(newImagePath)
	newImageStr := string(newImageOpen)
	for k := range p.Images {
		if strings.Contains(k, fmt.Sprintf("ppt/media/image%d.", index)) {
			p.Images[k] = newImageStr
		}
	}
	return nil
}

func (p *PowerPoint) WriteToFile(path string) (err error) {
	var target *os.File
	target, err = os.Create(path)
	if err != nil {
		return
	}

	defer target.Close()
	err = p.Write(target)
	return
}

func (p *PowerPoint) Write(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	defer w.Close()
	for _, file := range p.Files {
		var writer io.Writer
		var readCloser io.ReadCloser
		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}

		if strings.Contains(file.Name, "ppt/slides/slide") && p.Slides[file.Name] != "" {
			writer.Write([]byte(p.Slides[file.Name]))
		} else if strings.Contains(file.Name, "ppt/notesSlides/notesSlide") && p.NotesSlides[file.Name] != "" {
			writer.Write([]byte(p.NotesSlides[file.Name]))
		} else if strings.Contains(file.Name, "ppt/theme/theme") && p.Themes[file.Name] != "" {
			writer.Write([]byte(p.Themes[file.Name]))
		} else if file.Name == "ppt/presentation.xml" {
			writer.Write([]byte(p.Presentation))
		} else if strings.Contains(file.Name, "ppt/media/image") && p.Images[file.Name] != "" {
			writer.Write([]byte(p.Images[file.Name]))
		} else {
			readCloser, _ = file.Open()
			writer.Write(readCloserToByte(readCloser))
		}
	}
	return
}
func readCloserToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
