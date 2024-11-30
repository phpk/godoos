package config

import (
	"encoding/json"
	"fmt"
	"godo/ai/types"
	"godo/libs"
	"os"
	"path/filepath"
	"sync"
)

var ReqBodyMap = sync.Map{}

func GetConfigFile() (string, error) {
	modelDir, err := libs.GetAppDir()
	if err != nil {
		return "", err
	}
	if !libs.PathExists(modelDir) {
		os.MkdirAll(modelDir, 0755)
	}
	configFile := filepath.Join(modelDir, "ai_model.json")
	if !libs.PathExists(configFile) {
		// 如果文件不存在，则创建一个空的配置文件
		err := os.WriteFile(configFile, []byte("[]"), 0644)
		if err != nil {
			return "", err
		}
	}
	return configFile, nil
}

// LoadConfig 从文件加载所有ReqBody到映射中，如果文件不存在则创建一个空文件
func LoadConfig() error {
	filePath, err := GetConfigFile()
	if err != nil {
		return err
	}
	var reqBodies []types.ReqBody
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &reqBodies)
	if err != nil {
		return err
	}
	for _, reqBody := range reqBodies {
		ReqBodyMap.Store(reqBody.Model, reqBody)
	}
	//log.Printf("Load config file success %v", ReqBodyMap)
	return nil
}

// SaveReqBodiesToFile 将映射中的所有ReqBody保存回文件
func SaveConfig() error {
	filePath, err := GetConfigFile()
	if err != nil {
		return err
	}
	var reqBodies []types.ReqBody
	ReqBodyMap.Range(func(key, value interface{}) bool {
		rb := value.(types.ReqBody)
		reqBodies = append(reqBodies, rb)
		return true
	})

	// 使用 json.MarshalIndent 直接获取内容的字节切片
	content, err := json.MarshalIndent(reqBodies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
	}
	// log.Printf("====content: %s", string(content))
	// 将字节切片直接写入文件，避免了 string(content) 的冗余转换
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func GetModel(Model string) (types.ReqBody, bool) {
	value, ok := ReqBodyMap.Load(Model)
	if ok {
		return value.(types.ReqBody), true
	}
	return types.ReqBody{}, false
}
func ExistModel(Model string) bool {
	_, exists := ReqBodyMap.Load(Model)
	return exists
}
func SetModel(reqBody types.ReqBody) error {

	ReqBodyMap.Store(reqBody.Model, reqBody)

	//log.Println("=====SetModel", reqBody.Model)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated model configuration: %w", err)
	}

	return nil
}

func UpdateModel(reqBody types.ReqBody) error {
	_, loaded := ReqBodyMap.Load(reqBody.Model)
	if !loaded {
		return fmt.Errorf("model directory %s not found", reqBody.Model)
	}

	ReqBodyMap.Store(reqBody.Model, reqBody)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated model configuration: %w", err)
	}

	return nil
}

func AddModel(Model string, reqBody types.ReqBody) error {
	_, loaded := ReqBodyMap.Load(Model)
	if loaded {
		return fmt.Errorf("model directory %s already exists", Model)
	}

	ReqBodyMap.Store(Model, reqBody)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save new model configuration: %w", err)
	}

	return nil
}

func DeleteModel(Model string) error {
	_, loaded := ReqBodyMap.Load(Model)
	if loaded {
		ReqBodyMap.Delete(Model)
	}
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to delete model configuration: %w", err)
	}

	return nil
}
