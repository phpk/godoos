package common

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

func GetInfoFile() (string, error) {
	baseDir, err := GetAppDir()
	if err != nil {
		return "", err
	}
	if !PathExists(baseDir) {
		os.MkdirAll(baseDir, 0755)
	}
	InfoFile := filepath.Join(baseDir, "info")
	//log.Printf("Info file path: %s", InfoFile)
	if !PathExists(InfoFile) {
		// 如果文件不存在，则创建一个空的配置文件
		err = os.WriteFile(InfoFile, []byte(""), 0644)
		if err != nil {
			return "", err
		}
	}
	return InfoFile, nil
}

// LoadInfo 从文件加载所有ReqBody到映射中，如果文件不存在则创建一个空文件
func LoadInfo() error {
	filePath, err := GetInfoFile()
	if err != nil {
		return fmt.Errorf("error getting info file path: %w", err)
	}
	var reqBodies []ReqBody
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(content) < 1 {
		return nil
	}

	err = json.Unmarshal(content, &reqBodies)
	if err != nil {
		return err
	}
	for _, reqBody := range reqBodies {
		reqBodyMap.Store(reqBody.Name, reqBody)
	}
	//log.Printf("Load Info file success %v", reqBodyMap)
	return nil
}

// SaveReqBodiesToFile 将映射中的所有ReqBody保存回文件
func SaveInfo() error {
	filePath, err := GetInfoFile()
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

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func GetInfo(Name string) (any, bool) {
	value, ok := reqBodyMap.Load(Name)
	if ok {
		return value.(ReqBody).Value, true
	}
	return "", false
}
func ExistInfo(Name string) bool {
	_, exists := reqBodyMap.Load(Name)
	return exists
}
func SetInfo(reqBody ReqBody) error {

	reqBodyMap.Store(reqBody.Name, reqBody)
	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to save updated Infouration: %w", err)
	}

	return nil
}
func SetInfoByName(Name string, Value any) error {
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

	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to save updated Infouration: %w", err)
	}
	return nil
}
func SetInfos(reqBody []ReqBody) error {
	for _, rb := range reqBody {
		reqBodyMap.Store(rb.Name, rb)
	}

	//log.Println("=====SetName", reqBody.Name)
	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to save updated Infouration: %w", err)
	}

	return nil
}
func UpdateInfo(reqBody ReqBody) error {
	_, loaded := reqBodyMap.Load(reqBody.Name)
	if !loaded {
		return fmt.Errorf("Info directory %s not found", reqBody.Name)
	}

	reqBodyMap.Store(reqBody.Name, reqBody)
	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to save updated Infouration: %w", err)
	}

	return nil
}

func AddInfo(Name string, reqBody ReqBody) error {
	_, loaded := reqBodyMap.Load(Name)
	if loaded {
		return fmt.Errorf("Info directory %s already exists", Name)
	}

	reqBodyMap.Store(Name, reqBody)
	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to save new Infouration: %w", err)
	}

	return nil
}

func DeleteInfo(Name string) error {
	_, loaded := reqBodyMap.Load(Name)
	if loaded {
		reqBodyMap.Delete(Name)
	}
	if err := SaveInfo(); err != nil {
		return fmt.Errorf("failed to delete Infouration: %w", err)
	}

	return nil
}
