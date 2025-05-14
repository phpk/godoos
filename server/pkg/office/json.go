package office

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// extractTextFromJSON 递归地从 JSON 数据中提取纯文本
func extractTextFromJSON(data interface{}) []string {
	var texts []string

	switch v := data.(type) {
	case map[string]interface{}:
		for _, value := range v {
			texts = append(texts, extractTextFromJSON(value)...)
		}
	case []interface{}:
		for _, item := range v {
			texts = append(texts, extractTextFromJSON(item)...)
		}
	case string:
		texts = append(texts, v)
	default:
		// 其他类型忽略
	}

	return texts
}

func json2txt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	var jsonData interface{}
	err = json.Unmarshal(byteValue, &jsonData)
	if err != nil {
		return "", err
	}

	plainText := extractTextFromJSON(jsonData)

	// 将切片中的所有字符串连接成一个字符串
	plainTextStr := strings.Join(plainText, " ")

	// 移除多余的空格
	re := regexp.MustCompile(`\s+`)
	plainTextStr = re.ReplaceAllString(plainTextStr, " ")

	// 移除开头和结尾的空格
	plainTextStr = strings.TrimSpace(plainTextStr)

	return plainTextStr, nil
}
