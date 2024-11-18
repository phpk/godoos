package model

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetDownDir(modelPath string) (string, error) {
	baseDir, err := GetHfModelDir()
	if err != nil {
		return "", err
	}
	modelDir := filepath.Join(baseDir, modelPath)
	if !libs.PathExists(modelDir) {
		os.MkdirAll(modelDir, 0755)
	}
	return modelDir, nil
}
func GetModelDir(model string) (string, error) {
	modelName := ReplaceModelName(model)
	modelDir, err := GetDownDir(modelName)
	if err != nil {
		return "", err
	}
	return modelDir, nil
}
func GetModelPath(urls string, model string, reqType string) (string, error) {
	modelDir, err := GetModelDir(model)
	if err != nil {
		return "", err
	}
	//filePath := filepath.Join(modelDir, filepath.Base(reqBody.DownloadUrl))
	//log.Printf("====url: %s", urls)
	var fileName string
	pathParts := strings.Split(urls, "/")
	if len(pathParts) > 0 { // 确保路径有部分可分割
		fileName = pathParts[len(pathParts)-1] // 获取路径最后一部分
	} else {
		parsedUrl, err := url.Parse(urls)
		if err != nil {
			return "", fmt.Errorf("failed to parse URL: %w", err)
		}
		urlPath := parsedUrl.Path
		fileName = filepath.Base(urlPath)
	}
	// 构建完整的文件路径
	filePath := filepath.Join(modelDir, fileName)
	if reqType == "local" {
		dir, err := getOModelsDir()
		if err != nil {
			return "", err
		}
		if strings.Contains(fileName, "sha256-") && len(fileName) == 71 {
			filePath = filepath.Join(dir, "blobs", fileName)
			log.Printf("====filePath1: %s", filePath)
		} else {
			opName := getOpName(model)
			filePath = filepath.Join(dir, "manifests", opName.Space, opName.LibPath, opName.Name, opName.Tag)
			log.Printf("====filePath2: %s", filePath)
		}
	}
	return filePath, nil
}
func GetHfModelDir() (string, error) {
	dataDir := libs.GetDataDir()
	return filepath.Join(dataDir, "hfmodels"), nil
}
func GetOllamaModelDir() string {
	dataDir := libs.GetDataDir()
	return filepath.Join(dataDir, "models")
}
func getOModelsDir() (string, error) {
	return GetOllamaModelDir(), nil
}
func GetOllamaUrl() string {
	return "http://localhost:11434"
}
func ReplaceModelName(modelName string) string {
	reg := regexp.MustCompile(`[/\s:]`)
	return reg.ReplaceAllString(modelName, "")
}

// ModelConfigFromRequest 解析HTTP请求中的JSON数据并填充ModelConfig，如果请求中没有'modelconfig'键或解析出错，则返回一个空的ModelConfig
func ModelConfigFromRequest(r *http.Request) ModelConfig {
	// 初始化一个空的ModelConfig
	var config ModelConfig

	// 尝试解析请求体中的JSON数据
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err == nil {
		// 检查请求数据中是否存在'modelconfig'键
		if modelConfigData, ok := requestData["options"].(map[string]interface{}); ok {
			// 尝试将modelconfig数据转换为ModelConfig结构体
			jsonData, _ := json.Marshal(modelConfigData)
			if err := json.Unmarshal(jsonData, &config); err == nil {
				// 成功解析modelconfig数据到config
				return config
			}
		}
	}

	// 如果没有'modelconfig'键或者解析出错，直接返回一个空的ModelConfig
	return ModelConfig{}
}
func getIntInfo(val interface{}) int64 {
	if val, ok := val.(float64); ok {
		return int64(val)
	}
	return 0 // 如果键不存在或值不是期望的类型，则返回0
}
