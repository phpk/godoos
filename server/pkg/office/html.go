package office

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func html2txt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return ConvertHTML(file)

}

// 去除字符串中的HTML标签
func TrimHtml(text string) string {
	// 去除字符串中的HTML标签
	re := regexp.MustCompile(`</?\w+[^>]*>`)
	text = re.ReplaceAllString(text, "")
	// 先去除所有空格
	text = strings.ReplaceAll(text, "  ", "")
	// 合并多个连续换行为单个换行
	re = regexp.MustCompile(`\n+`)
	text = re.ReplaceAllString(text, "\n")
	return strings.TrimSpace(text)
}

func ConvertHTML(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return ``, err
	}

	var title string
	var content bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if firstChild := n.FirstChild; firstChild != nil && firstChild.Type == html.TextNode {
					title = firstChild.Data
				}
			case "script", "style", "img": // 忽略这些元素及其子元素
				return
			}
		}

		if n.Type == html.TextNode {
			content.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	contentStr := content.String()
	contentStr = TrimHtml(contentStr)

	if title == "" {
		title = "未命名网页"
	}

	return title + "\n" + strings.TrimSpace(contentStr), nil
}
