// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package libs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var reqBodyMap = sync.Map{}

type ReqBody struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func GetConfigFile() (string, error) {
	baseDir, err := GetAppDir()
	if err != nil {
		return "", err
	}
	if !PathExists(baseDir) {
		os.MkdirAll(baseDir, 0755)
	}
	configFile := filepath.Join(baseDir, "os_config.json")
	//log.Printf("config file path: %s", configFile)
	if !PathExists(configFile) {
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
	var reqBodies []ReqBody
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &reqBodies)
	if err != nil {
		return err
	}
	for _, reqBody := range reqBodies {
		reqBodyMap.Store(reqBody.Name, reqBody)
	}
	//log.Printf("Load config file success %v", reqBodyMap)
	return nil
}

// SaveReqBodiesToFile 将映射中的所有ReqBody保存回文件
func SaveConfig() error {
	filePath, err := GetConfigFile()
	if err != nil {
		return err
	}
	// 创建一个 map 用来存储已遇到的 Name，以防止重复
	seenNames := make(map[string]bool)
	var reqBodies []ReqBody
	reqBodyMap.Range(func(key, value interface{}) bool {
		rb := value.(ReqBody)
		// 如果 rb.Name 还没有出现在 seenNames 中，才添加到 reqBodies
		if _, exists := seenNames[rb.Name]; !exists {
			seenNames[rb.Name] = true
			reqBodies = append(reqBodies, rb)
		}
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
func GetConfig(Name string) (any, bool) {
	value, ok := reqBodyMap.Load(Name)
	if ok {
		return value.(ReqBody).Value, true
	}
	return "", false
}
func ExistConfig(Name string) bool {
	_, exists := reqBodyMap.Load(Name)
	return exists
}
func SetConfig(reqBody ReqBody) error {

	reqBodyMap.Store(reqBody.Name, reqBody)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}

	return nil
}
func SetConfigByName(Name string, Value any) error {
	// 尝试从 reqBodyMap 中加载 Name
	value, ok := reqBodyMap.Load(Name)
	if !ok {
		// 如果 Name 不存在，则创建一个新的 ReqBody
		newReqBody := ReqBody{Name: Name, Value: Value}
		reqBodyMap.Store(Name, newReqBody)
	} else {
		// 如果 Name 存在，则更新现有 ReqBody 的 Value
		existingReqBody, _ := value.(ReqBody)
		existingReqBody.Value = Value
		reqBodyMap.Store(Name, existingReqBody)
	}

	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}
	return nil
}
func SetConfigs(reqBody []ReqBody) error {
	for _, rb := range reqBody {
		reqBodyMap.Store(rb.Name, rb)
	}

	//log.Println("=====SetName", reqBody.Name)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}

	return nil
}
func UpdateConfig(reqBody ReqBody) error {
	_, loaded := reqBodyMap.Load(reqBody.Name)
	if !loaded {
		return fmt.Errorf("config directory %s not found", reqBody.Name)
	}

	reqBodyMap.Store(reqBody.Name, reqBody)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}

	return nil
}

func AddConfig(Name string, reqBody ReqBody) error {
	_, loaded := reqBodyMap.Load(Name)
	if loaded {
		return fmt.Errorf("config directory %s already exists", Name)
	}

	reqBodyMap.Store(Name, reqBody)
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to save new configuration: %w", err)
	}

	return nil
}

func DeleteConfig(Name string) error {
	_, loaded := reqBodyMap.Load(Name)
	if loaded {
		reqBodyMap.Delete(Name)
	}
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("failed to delete configuration: %w", err)
	}

	return nil
}
