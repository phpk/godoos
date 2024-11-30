// /*
//   - GodoOS - A lightweight cloud desktop
//   - Copyright (C) 2024 https://godoos.com
//     *
//   - This program is free software: you can redistribute it and/or modify
//   - it under the terms of the GNU Lesser General Public License as published by
//   - the Free Software Foundation, either version 2.1 of the License, or
//   - (at your option) any later version.
//     *
//   - This program is distributed in the hope that it will be useful,
//   - but WITHOUT ANY WARRANTY; without even the implied warranty of
//   - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   - GNU Lesser General Public License for more details.
//     *
//   - You should have received a copy of the GNU Lesser General Public License
//   - along with this program.  If not, see <http://www.gnu.org/licenses/>.
//     */
package convert

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"godo/libs"
// 	"io"
// 	"log"
// 	"mime"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// // UploadInfo 用于表示上传文件的信息
// type UploadInfo struct {
// 	Name      string    `json:"name"`
// 	SavePath  string    `json:"save_path"`
// 	Content   string    `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // SaveContentToFile 保存内容到文件并返回UploadInfo结构体
// func SaveContentToFile(content, fileName string) (UploadInfo, error) {
// 	uploadBaseDir, err := libs.GetUploadDir()
// 	if err != nil {
// 		return UploadInfo{}, err
// 	}

// 	// 去除文件名中的空格
// 	fileNameWithoutSpaces := strings.ReplaceAll(fileName, " ", "_")
// 	fileNameWithoutSpaces = strings.ReplaceAll(fileNameWithoutSpaces, "/", "")
// 	fileNameWithoutSpaces = strings.ReplaceAll(fileNameWithoutSpaces, `\`, "")
// 	// 提取文件名和扩展名
// 	// 查找最后一个点的位置
// 	lastDotIndex := strings.LastIndexByte(fileNameWithoutSpaces, '.')

// 	// 如果找到点，则提取扩展名，否则视为没有扩展名
// 	ext := ""
// 	if lastDotIndex != -1 {
// 		ext = fileNameWithoutSpaces[lastDotIndex:]
// 		fileNameWithoutSpaces = fileNameWithoutSpaces[:lastDotIndex]
// 	} else {
// 		ext = ""
// 	}
// 	randFileName := fmt.Sprintf("%s_%s%s", fileNameWithoutSpaces, strconv.FormatInt(time.Now().UnixNano(), 10), ext)
// 	savePath := filepath.Join(uploadBaseDir, time.Now().Format("2006-01-02"), randFileName)

// 	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
// 		return UploadInfo{}, err
// 	}

// 	if err := os.WriteFile(savePath, []byte(content), 0644); err != nil {
// 		return UploadInfo{}, err
// 	}

// 	return UploadInfo{
// 		Name:     fileNameWithoutSpaces,
// 		SavePath: savePath,
// 		//Content:   content,
// 		CreatedAt: time.Now(),
// 	}, nil
// }

// // MultiUploadHandler 处理多文件上传请求
// func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseMultipartForm(10000 << 20); err != nil {
// 		libs.Error(w, "Failed to parse multipart form")
// 		return
// 	}

// 	files := r.MultipartForm.File["files"]
// 	if len(files) == 0 {
// 		libs.Error(w, "No file parts in the request")
// 		return
// 	}

// 	fileInfoList := make([]UploadInfo, 0, len(files))

// 	for _, fileHeader := range files {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			libs.Error(w, "Failed to open uploaded file")
// 			continue
// 		}
// 		defer file.Close()

// 		content, err := io.ReadAll(file)
// 		if err != nil {
// 			libs.Error(w, "Failed to read uploaded file")
// 			continue
// 		}

// 		//log.Printf(string(content))
// 		// 保存上传的文件内容
// 		info, err := SaveContentToFile(string(content), fileHeader.Filename)
// 		if err != nil {
// 			libs.Error(w, "Failed to save uploaded file")
// 			continue
// 		}
// 		log.Println(info.SavePath)
// 		// 对上传的文件进行转换处理
// 		convertData := Convert(info.SavePath) // Assuming convert.Convert expects a file path
// 		log.Printf("convertData: %v", convertData)
// 		if convertData.Data == "" {
// 			continue
// 		}
// 		images := []ImagesInfo{}
// 		resInfo := ResContentInfo{
// 			Content: convertData.Data,
// 			Images:  images,
// 		}
// 		// 将转换后的数据写入文件
// 		savePath := info.SavePath + "_result.json"
// 		// if err := WriteConvertedDataToFile(convertData.Data, savePath); err != nil {
// 		// 	serv.Err("Failed to write converted data to file", w)
// 		// 	continue
// 		// }
// 		// 使用 json.MarshalIndent 直接获取内容的字节切片
// 		contents, err := json.MarshalIndent(resInfo, "", "  ")
// 		if err != nil {
// 			libs.Error(w, "failed to marshal reqBodies to JSON:"+savePath)
// 			continue
// 		}
// 		// 将字节切片直接写入文件
// 		if err := os.WriteFile(savePath, contents, 0644); err != nil {
// 			libs.Error(w, "failed to write to file:"+savePath)
// 			continue
// 		}

// 		//info.SavePath = savePath
// 		fileInfoList = append(fileInfoList, info)
// 	}

// 	libs.Success(w, fileInfoList, "success")
// }

// // WriteConvertedDataToFile 将转换后的数据写入文件
// func WriteConvertedDataToFile(data, filePath string) error {
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(data)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Successfully wrote %d bytes to file %s.\n", len(data), filePath)
// 	return nil
// }

// // jsonParamHandler 处理JSON参数请求
// func JsonParamHandler(w http.ResponseWriter, r *http.Request) {
// 	type RequestBody struct {
// 		Path string `json:"path"`
// 	}

// 	var requestBody RequestBody
// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		libs.Error(w, "Invalid request body")
// 		return
// 	}

// 	path := requestBody.Path
// 	fmt.Printf("Parameter 'path' from JSON is: %s\n", path)

// 	if path != "" {
// 		resp := Convert(path)
// 		w.Header().Set("Content-Type", "application/json")
// 		if err := json.NewEncoder(w).Encode(resp); err != nil {
// 			libs.Error(w, "Error encoding JSON")
// 			return
// 		}
// 		return
// 	}
// }

// // HandleURLPost 接收一个POST请求，其中包含一个URL参数，然后处理该URL指向的内容并保存
// func HandleURLPost(w http.ResponseWriter, r *http.Request) {
// 	var requestBody struct {
// 		URL string `json:"url"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&requestBody); err != nil {
// 		libs.Error(w, "Invalid request body")
// 		return
// 	}
// 	resp, err := http.Get(requestBody.URL)
// 	if err != nil {
// 		libs.Error(w, "Invalid request url:"+requestBody.URL)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, errRead := io.ReadAll(resp.Body)
// 	if errRead != nil {
// 		libs.Error(w, "Invalid request body")
// 		return
// 	}
// 	reader := bytes.NewReader(body)
// 	res, err := ConvertHTML(reader)
// 	if err != nil {
// 		libs.Error(w, "Failed to convert content")
// 		return
// 	}
// 	log.Printf("Converted content: %s", res)
// 	// 使用通用的SaveContentToFile函数保存内容到文件
// 	//fileName := "converted_from_url"
// 	// 获取内容的第一行作为标题
// 	fileName := strings.SplitN(res, "\n", 2)[0]
// 	if fileName == "" {
// 		fileName = "未命名网页"
// 	}
// 	fileName = fileName + ".html"
// 	info, err := SaveContentToFile(res, fileName)
// 	if err != nil {
// 		libs.Error(w, "Failed to save converted content to file")
// 		return
// 	}
// 	// 将转换后的数据写入文件
// 	savePath := info.SavePath + "_result.json"
// 	// if err := WriteConvertedDataToFile(info.Content, savePath); err != nil {
// 	// 	serv.Err("Failed to write converted data to file", w)
// 	// 	return
// 	// }
// 	// 使用 json.MarshalIndent 直接获取内容的字节切片
// 	resInfo := ResContentInfo{
// 		Content: info.Content,
// 	}
// 	contents, err := json.MarshalIndent(resInfo, "", "  ")
// 	if err != nil {
// 		libs.Error(w, "failed to marshal reqBodies to JSON:"+savePath)
// 		return
// 	}
// 	// 将字节切片直接写入文件
// 	if err := os.WriteFile(savePath, contents, 0644); err != nil {
// 		libs.Error(w, "failed to write to file:"+savePath)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(info); err != nil {
// 		libs.Error(w, "Error encoding JSON")
// 		return
// 	}
// }
// func ShowDetailHandler(w http.ResponseWriter, r *http.Request) {
// 	// 从 URL 查询参数中获取图片路径
// 	filePath := r.URL.Query().Get("path")
// 	//log.Printf("imagePath: %s", imagePath)
// 	// 检查图片路径是否为空或无效
// 	if filePath == "" {
// 		libs.Error(w, "Invalid file path")
// 		return
// 	}
// 	var reqBodies ResContentInfo

// 	if libs.PathExists(filePath + "_result.json") {
// 		//log.Printf("ShowDetailHandler: %s", filePath)
// 		filePath = filePath + "_result.json"
// 		content, err := os.ReadFile(filePath)
// 		if err != nil {
// 			libs.Error(w, "Failed to open file")
// 			return
// 		}
// 		err = json.Unmarshal(content, &reqBodies)
// 		if err != nil {
// 			libs.Error(w, "Failed to read file")
// 			return
// 		}
// 		// 设置响应头
// 		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 		resContent := reqBodies.Content + "/n"
// 		for _, image := range reqBodies.Images {
// 			resContent += image.Content + "/n"
// 		}
// 		// 写入响应体
// 		_, err = w.Write([]byte(resContent))
// 		if err != nil {
// 			libs.Error(w, "Failed to write response")
// 			return
// 		}
// 	} else {
// 		// 确保图片路径是绝对路径
// 		absImagePath, err := filepath.Abs(filePath)
// 		//log.Printf("absImagePath: %s", absImagePath)
// 		if err != nil {
// 			libs.Error(w, err.Error())
// 			return
// 		}

// 		// 获取文件的 MIME 类型
// 		mimeType := mime.TypeByExtension(filepath.Ext(absImagePath))
// 		if mimeType == "" {
// 			mimeType = "application/octet-stream" // 如果无法识别，就用默认的二进制流类型
// 		}

// 		// 设置响应头的 MIME 类型
// 		w.Header().Set("Content-Type", mimeType)

// 		// 打开文件并读取内容
// 		file, err := os.Open(absImagePath)
// 		if err != nil {
// 			libs.Error(w, err.Error())
// 			return
// 		}
// 		defer file.Close()

// 		// 将文件内容写入响应体
// 		_, err = io.Copy(w, file)
// 		if err != nil {
// 			libs.Error(w, err.Error())
// 		}
// 	}

// }
// func ServeImage(w http.ResponseWriter, r *http.Request) {
// 	// 从 URL 查询参数中获取图片路径
// 	imagePath := r.URL.Query().Get("path")
// 	//log.Printf("imagePath: %s", imagePath)
// 	// 检查图片路径是否为空或无效
// 	if imagePath == "" {
// 		libs.Error(w, "Invalid image path")
// 		return
// 	}

// 	// 确保图片路径是绝对路径
// 	absImagePath, err := filepath.Abs(imagePath)
// 	//log.Printf("absImagePath: %s", absImagePath)
// 	if err != nil {
// 		libs.Error(w, err.Error())
// 		return
// 	}

// 	// 获取文件的 MIME 类型
// 	mimeType := mime.TypeByExtension(filepath.Ext(absImagePath))
// 	if mimeType == "" {
// 		mimeType = "application/octet-stream" // 如果无法识别，就用默认的二进制流类型
// 	}

// 	// 设置响应头的 MIME 类型
// 	w.Header().Set("Content-Type", mimeType)

// 	// 打开文件并读取内容
// 	file, err := os.Open(absImagePath)
// 	if err != nil {
// 		libs.Error(w, err.Error())
// 		return
// 	}
// 	defer file.Close()

// 	// 将文件内容写入响应体
// 	_, err = io.Copy(w, file)
// 	if err != nil {
// 		libs.Error(w, err.Error())
// 	}
// }
