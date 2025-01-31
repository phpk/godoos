/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/

package office

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"godo/libs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type DocResult struct {
	filePath    string
	newFilePath string
	err         error
}

func SetDocument(dirPath string, knowledgeId uint) error {
	if !libs.PathExists(dirPath) {
		return nil
	}
	var wg sync.WaitGroup
	results := make(chan DocResult, 100) // 缓冲通道

	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 获取文件名
			fileName := filepath.Base(filePath)
			// 检查文件名是否以点开头
			if len(fileName) > 0 && fileName[0] == '.' {
				return nil // 跳过以点开头的文件
			}
			// 获取文件扩展名
			ext := filepath.Ext(filePath)
			// 检查文件扩展名是否为 .exe
			if ext == ".exe" {
				return nil // 跳过 .exe 文件
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				result := ProcessFile(filePath, knowledgeId)
				results <- result
			}(filePath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			fmt.Printf("Failed to process file %s: %v\n", result.filePath, result.err)
		} else {
			fmt.Printf("Processed file %s and saved JSON to %s\n", result.filePath, result.newFilePath)
		}
	}

	return nil
}

func ProcessFile(filePath string, knowledgeId uint) DocResult {
	doc, err := GetDocument(filePath)
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	//newFileName := ".godoos." + filepath.Base(filePath)
	newFileName := fmt.Sprintf(".godoos.%d.%s.json", knowledgeId, filepath.Base(filePath))
	newFilePath := filepath.Join(filepath.Dir(filePath), newFileName)
	// log.Printf("New file name: %s", newFileName)
	// log.Printf("New file path: %s", newFilePath)
	err = os.WriteFile(newFilePath, jsonData, 0644)
	if err != nil {
		return DocResult{filePath: filePath, err: err}
	}

	return DocResult{filePath: filePath, newFilePath: newFilePath, err: nil}
}

// ProcessBase64File 处理解码后的文件并提取文本信息
func ProcessBase64File(base64String string, fileName string) (string, error) {
	// 解码 Base64 字符串
	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %v", err)
	}
	// 获取文件后缀并转换为小写
	// fileExt := strings.ToLower(filepath.Ext(fileName))
	// if fileExt == "" {
	// 	return "", fmt.Errorf("file extension not found in filename: %s", fileName)
	// }
	// log.Printf("File type: %s\n", fileExt)
	cacheDir := libs.GetCacheDir()
	tempFilePath := filepath.Join(cacheDir, fileName)

	// 创建临时文件
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// 将解码后的数据写入临时文件
	_, err = tempFile.Write(decodedBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}

	// 获取文档内容
	doc, err := GetDocument(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get document: %v", err)
	}
	log.Printf("Document content: %s\n", doc.Content)

	// 删除临时文件
	defer os.Remove(tempFilePath)

	// 提取文本内容
	return doc.Content, nil
}
func GetDocument(pathname string) (*Document, error) {
	if !libs.PathExists(pathname) {
		return nil, fmt.Errorf("file does not exist: %s", pathname)
	}
	abPath, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}
	filename := path.Base(pathname)
	data := Document{path: pathname, RePath: abPath, Title: filename}
	extension := path.Ext(pathname)
	_, err = getFileInfoData(&data)
	if err != nil {
		return &data, err
	}
	switch extension {
	case ".docx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, docx2txt)
	case ".pptx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, pptx2txt)
	case ".xlsx":
		_, e := getMetaData(&data)
		if e != nil {
			fmt.Printf("⚠️ %s", e.Error())
		}
		_, err = getContentData(&data, xlsx2txt)
	case ".pdf":
		_, err = getContentData(&data, pdf2txt)
	case ".doc":
		_, err = getContentData(&data, doc2txt)
	case ".ppt":
		_, err = getContentData(&data, ppt2txt)
	case ".xls":
		_, err = getContentData(&data, xls2txt)
	case ".epub":
		_, err = getContentData(&data, epub2txt)
	case ".odt":
		_, err = getContentData(&data, odt2txt)
	case ".xml":
		_, err = getContentData(&data, xml2txt)
	case ".rtf":
		_, err = getContentData(&data, rtf2txt)
	case ".md":
		_, err = getContentData(&data, md2txt)
	case ".txt":
		_, err = getContentData(&data, text2txt)
	case ".xhtml", ".html", ".htm":
		_, err = getContentData(&data, html2txt)
	case ".json":
		_, err = getContentData(&data, json2txt)
	// case ".md", ".txt", ".py", ".java", ".c", ".cpp", ".h", ".hpp", ".js", ".ts", ".go", ".rb", ".php", ".swift", ".kt", ".scala", ".rust", ".perl", ".bash", ".sh", ".lua", ".dart", ".r", ".matlab", ".pl", ".pm", ".tcl", ".sql", ".groovy", ".cs", ".vb", ".fs", ".hs", ".erl", ".elixir", ".crystal", ".nim", ".d", ".coffeescript", ".typescript", ".vue", ".svelte", ".jsx", ".tsx", ".html", ".css", ".scss", ".less", ".json", ".xml", ".yaml", ".yml", ".toml", ".ini", ".makefile", ".dockerfile", ".gitignore", ".editorconfig", ".prettierrc", ".eslintrc", ".babelrc", ".jsonp", ".graphql", ".proto", ".plist", ".edn":
	// 	_, err = getContentData(&data, text2txt)
	default:
		_, err = getContentData(&data, text2txt)
	}
	if err != nil {
		return &data, err
	}
	return &data, nil
}
func GetTxtDoc(pathname string) (*Document, error) {
	if !libs.PathExists(pathname) {
		return nil, fmt.Errorf("file does not exist: %s", pathname)
	}
	abPath, err := filepath.Abs(pathname)
	if err != nil {
		return nil, err
	}
	filename := path.Base(pathname)
	data := Document{path: pathname, RePath: abPath, Title: filename}
	_, err = getFileInfoData(&data)
	if err != nil {
		return &data, err
	}
	_, err = getContentData(&data, text2txt)
	if err != nil {
		return &data, err
	}
	return &data, nil
}

// Read the meta data of office files (only *.docx, *.xlsx, *.pptx) and insert into the interface
func getMetaData(data *Document) (bool, error) {
	file, err := os.Open(data.path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	meta, err := GetContent(file)
	if err != nil {
		return false, errors.New("failed to get office meta data")
	}
	if meta.Title != "" {
		data.Title = meta.Title
	}
	data.Subject = meta.Subject
	data.Creator = meta.Creator
	data.Keywords = meta.Keywords
	data.Description = meta.Description
	data.Lastmodifiedby = meta.LastModifiedBy
	data.Revision = meta.Revision
	data.Category = meta.Category
	data.Content = meta.Category
	return true, nil
}
func GetContent(document *os.File) (fields XMLContent, err error) {
	// Attempt to read the document file directly as a zip file.
	z, err := zip.OpenReader(document.Name())
	if err != nil {
		return fields, errors.New("failed to open the file as zip")
	}
	defer z.Close()

	var xmlFile string
	for _, file := range z.File {
		if file.Name == "docProps/core.xml" {
			rc, err := file.Open()
			if err != nil {
				return fields, errors.New("failed to open docProps/core.xml")
			}
			defer rc.Close()

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				xmlFile += scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				return fields, errors.New("failed to read from docProps/core.xml")
			}
			break // Exit loop after finding and reading core.xml
		}
	}

	// Unmarshal the collected XML content into the XMLContent struct
	if err := xml.Unmarshal([]byte(xmlFile), &fields); err != nil {
		return fields, errors.New("failed to Unmarshal")
	}

	return fields, nil
}

// Read the content of office files and insert into the interface
func getContentData(data *Document, reader DocReader) (bool, error) {
	content, err := reader(data.path)
	if err != nil {
		return false, err
	}
	data.Content = content
	data.Split = SplitText(content, 256)
	return true, nil
}
