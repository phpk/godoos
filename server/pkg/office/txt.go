package office

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func text2txt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 移除行内的换行符
		line = strings.ReplaceAll(line, "\r", "")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// 合并所有行
	content := strings.Join(lines, " ")

	// 移除多余的空格
	re := regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	// 移除开头和结尾的空格
	content = strings.TrimSpace(content)

	return content, nil
}
