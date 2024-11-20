package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OllamaModelsList struct {
	Models []OllamaModelsInfo `json:"models"`
}

type OllamaDetails struct {
	ParameterSize     string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}
type OllamaModelsInfo struct {
	Model   string        `json:"model"`
	Details OllamaDetails `json:"details"`
	Size    int64         `json:"size"`
}
type OllamaModelDetail struct {
	Parameters string                 `json:"parameters"`
	Template   string                 `json:"template"`
	Details    map[string]interface{} `json:"details"`
	ModelInfo  map[string]interface{} `json:"model_info"`
}
type ResModelInfo struct {
	Parameters      string `json:"parameters"`
	Template        string `json:"template"`
	ContextLength   int64  `json:"context_length"`
	EmbeddingLength int64  `json:"embedding_length"`
	Size            string `json:"size"`
	Quant           string `json:"quant"`
	Desk            string `json:"desk"`
	Cpu             string `json:"cpu"`
	Gpu             string `json:"gpu"`
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	From      string `json:"from,omitempty"`
	status    string
}
type ManifestV2 struct {
	SchemaVersion int      `json:"schemaVersion"`
	MediaType     string   `json:"mediaType"`
	Config        *Layer   `json:"config"`
	Layers        []*Layer `json:"layers"`
}
type OmodelPath struct {
	Space   string
	LibPath string
	Name    string
	Tag     string
}

const (
	KB = 1 << (10 * iota)
	MB
	GB
)
const (
	CPU_8GB  = "8GB"
	CPU_16GB = "16GB"
	CPU_32GB = "32GB"
	GPU_6GB  = "6GB"
	GPU_8GB  = "8GB"
	GPU_12GB = "12GB"
)

func humanReadableSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	unitIndex := 0 // Start with Bytes
	for size >= 1000 && unitIndex < len(units)-1 {
		size /= 1000
		unitIndex++
	}

	switch unitIndex {
	case 0, 1, 2, 3: // For B, KB, and MB, keep decimal points
		return fmt.Sprintf("%d%s", size, units[unitIndex])
	default:
		return fmt.Sprintf("%dB", size) // Fallback for sizes less than 1B or unhandled cases
	}
}

func Tagshandler(w http.ResponseWriter, r *http.Request) {
	err := LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error")
		return
	}
	var reqBodies []ReqBody
	reqBodyMap.Range(func(key, value interface{}) bool {
		rb, ok := value.(ReqBody)
		if ok {
			reqBodies = append(reqBodies, rb)
		}
		return true // 继续遍历
	})
	// 对reqBodies按CreatedAt降序排列
	sort.Slice(reqBodies, func(i, j int) bool {
		return reqBodies[i].CreatedAt.After(reqBodies[j].CreatedAt) // 降序排列
	})
	// 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 使用json.NewEncoder将reqBodies编码为JSON并写入响应体
	if err := json.NewEncoder(w).Encode(reqBodies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	err := LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "Load config error")
		return
	}
	model := r.URL.Query().Get("model")
	if model == "" {
		libs.ErrorMsg(w, "Model name is empty")
		return
	}
	//log.Printf("ShowHandler: %s", model)
	var reqBodies ReqBody
	reqBodyMap.Range(func(key, value interface{}) bool {
		rb, ok := value.(ReqBody)
		if ok && rb.Model == model {
			reqBodies = rb
			return false
		}
		return true
	})
	//log.Printf("ShowHandler: %s", reqBodies)
	// 设置响应内容类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 使用json.NewEncoder将reqBodies编码为JSON并写入响应体
	if err := json.NewEncoder(w).Encode(reqBodies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func extractParameterSize(sizeStr string, model string) (float64, bool) {
	// log.Printf("extractParameterSize: %s", sizeStr)
	// log.Printf("extractParameterModel: %s", model)
	// 尝试直接从原始sizeStr中提取数字，包括小数
	if size, err := strconv.ParseFloat(strings.TrimSuffix(sizeStr, "B"), 64); err == nil {
		return size, true
	}

	if parts := strings.Split(model, ":"); len(parts) > 1 {
		// 确保移除 "b" 或 "B" 后缀，并尝试转换为浮点数
		cleanedPart := strings.TrimSuffix(strings.ToLower(parts[1]), "b")
		if size, err := strconv.ParseFloat(cleanedPart, 64); err == nil {
			return size, true
		}
	}

	return 0, false
}

func parseOllamaInfo(info OllamaModelsInfo) ResModelInfo {
	res := ResModelInfo{
		Size:  humanReadableSize(info.Size),
		Quant: info.Details.QuantizationLevel,
	}
	res.Desk = res.Size
	paramSize, ok := extractParameterSize(info.Details.ParameterSize, info.Model)
	if !ok {
		res.Cpu = CPU_8GB
		res.Gpu = GPU_6GB
		return res
	}

	switch {
	case paramSize < 3:
		res.Cpu = CPU_8GB
		res.Gpu = GPU_6GB
	case paramSize < 9:
		res.Cpu = CPU_16GB
		res.Gpu = GPU_8GB
	default:
		res.Cpu = CPU_32GB
		res.Gpu = GPU_12GB
	}

	return res
}
func getOllamaModels() ([]OllamaModelsInfo, error) {
	req, err := http.Get(GetOllamaUrl() + "/api/tags")
	res := []OllamaModelsInfo{}
	if err != nil {
		return res, fmt.Errorf("failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return res, fmt.Errorf("failed to read response body")
	}
	rest := OllamaModelsList{}
	if err := json.Unmarshal(body, &rest); err != nil {
		return res, fmt.Errorf("failed to unmarshal response body")
	}
	return rest.Models, nil

}
func setOllamaInfo(w http.ResponseWriter, r *http.Request, reqBody ReqBody) {
	model := reqBody.Model
	postQuery := map[string]interface{}{
		"model": model,
	}
	url := GetOllamaUrl() + "/api/pull"
	ForwardHandler(w, r, postQuery, url, "POST")
	details, err := getOllamaInfo(r, model)
	//log.Printf("details is %v", details)
	if err != nil {
		libs.ErrorMsg(w, "get ollama info error: ")
		return
	}

	modelList, err := getOllamaModels()
	if err != nil {
		libs.ErrorMsg(w, "Load ollama error: ")
		return
	}
	if len(modelList) < 1 {
		libs.ErrorMsg(w, "Load ollama error: ")
		return
	}
	for _, model := range modelList {
		if model.Model == reqBody.Model {
			oinfo := parseOllamaInfo(model)
			architecture := details.ModelInfo["general.architecture"].(string)
			contextLength := details.ModelInfo[architecture+".context_length"]
			embeddingLength := details.ModelInfo[architecture+".embedding_length"]
			info := map[string]interface{}{
				"size":             oinfo.Size,
				"quant":            oinfo.Quant,
				"Desk":             oinfo.Desk,
				"cpu":              oinfo.Cpu,
				"gpu":              oinfo.Gpu,
				"pb":               model.Details.ParameterSize,
				"template":         details.Template,
				"parameters":       details.Parameters,
				"context_length":   contextLength,
				"embedding_length": embeddingLength,
			}
			paths, err := getManifests(model.Model)
			if err != nil {
				log.Printf("Error parsing Manifests: %v", err)
				continue
			}

			reqBody.Info = info
			reqBody.Paths = paths
			reqBody.Status = "success"
			reqBody.CreatedAt = time.Now()
			if err := SetModel(reqBody); err != nil {
				libs.ErrorMsg(w, "Set model error")
				return
			}
			return
		}
	}
}
func getOllamaInfo(r *http.Request, model string) (OllamaModelDetail, error) {
	infoQuery := map[string]interface{}{
		"name": model,
	}
	res := OllamaModelDetail{}
	url := GetOllamaUrl() + "/api/show"
	payloadBytes, err := json.Marshal(infoQuery)
	if err != nil {
		return res, fmt.Errorf("json payload error: %w", err)
	}
	// 创建POST请求，复用原始请求的上下文（如Cookies）
	req, err := http.NewRequestWithContext(r.Context(), "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return res, fmt.Errorf("couldn't create req context: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("couldn't read response body: %w", err)
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return res, fmt.Errorf("failed to unmarshal response body")
	}
	return res, nil
}
func getOpName(model string) OmodelPath {
	libPath := "library"
	modelName := model
	modelTags := "latest"
	if strings.Contains(modelName, ":") {
		names := strings.Split(model, ":")
		modelName = names[0]
		modelTags = names[1]
	}

	if strings.Contains(modelName, "/") {
		names := strings.Split(modelName, "/")
		libPath = names[0]
		modelName = names[1]
	}
	return OmodelPath{
		Space:   "registry.ollama.ai",
		LibPath: libPath,
		Name:    modelName,
		Tag:     modelTags,
	}
}
func getManifests(model string) ([]string, error) {
	res := []string{}
	opName := getOpName(model)
	modelsDir, err := getOModelsDir()
	if err != nil {
		return res, fmt.Errorf("failed to get user home directory: %w", err)
	}
	manifestsFile := filepath.Join(modelsDir, "manifests", opName.Space, opName.LibPath, opName.Name, opName.Tag)
	if !libs.PathExists(manifestsFile) {
		return res, fmt.Errorf("failed to get manifests file: %w", err)
	}
	res = append(res, manifestsFile)
	var manifest ManifestV2
	f, err := os.Open(manifestsFile)
	if err != nil {
		return res, err
	}
	defer f.Close()

	sha256sum := sha256.New()
	if err := json.NewDecoder(io.TeeReader(f, sha256sum)).Decode(&manifest); err != nil {
		return res, err
	}
	filename, err := GetBlobsPath(manifest.Config.Digest)
	if err != nil {
		return nil, err
	}
	res = append(res, filename)
	for _, layer := range manifest.Layers {
		filename, err := GetBlobsPath(layer.Digest)
		if err != nil {
			return nil, err
		}
		res = append(res, filename)
	}
	return res, nil
}

func GetBlobsPath(digest string) (string, error) {
	dir, err := getOModelsDir()
	if err != nil {
		return "", err
	}
	// only accept actual sha256 digests
	pattern := "^sha256[:-][0-9a-fA-F]{64}$"
	re := regexp.MustCompile(pattern)

	if digest != "" && !re.MatchString(digest) {
		return "", errors.New("invalid digest format")
	}

	digest = strings.ReplaceAll(digest, ":", "-")
	path := filepath.Join(dir, "blobs", digest)
	return path, nil
}
