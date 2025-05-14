package office

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	reHTML       = regexp.MustCompile(`<[^>]*>`)
	reMarkdown   = regexp.MustCompile(`[\*_|#]{1,4}`)
	reWhitespace = regexp.MustCompile(`\s+`)
)

func md2txt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 移除 HTML 标签
		line = reHTML.ReplaceAllString(line, "")
		// 移除 Markdown 格式符号
		line = reMarkdown.ReplaceAllString(line, "")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// 合并所有行
	content := strings.Join(lines, " ")

	// 移除多余的空格
	content = reWhitespace.ReplaceAllString(content, " ")

	// 移除开头和结尾的空格
	content = strings.TrimSpace(content)

	return content, nil
}
