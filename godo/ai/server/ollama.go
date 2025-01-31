package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"godo/ai/config"
	"godo/ai/types"
	"godo/libs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

func extractParameterSize(sizeStr string, model string) (float64, bool) {
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

func parseOllamaInfo(info types.OllamaModelsInfo) types.ModelInfo {
	res := types.ModelInfo{
		Size:  humanReadableSize(info.Size),
		Quant: info.Details.QuantizationLevel,
	}
	res.Desk = res.Size
	paramSize, ok := extractParameterSize(info.Details.ParameterSize, info.Model)
	if !ok {
		res.CPU = CPU_8GB
		res.GPU = GPU_6GB
		return res
	}

	switch {
	case paramSize < 3:
		res.CPU = CPU_8GB
		res.GPU = GPU_6GB
	case paramSize < 9:
		res.CPU = CPU_16GB
		res.GPU = GPU_8GB
	default:
		res.CPU = CPU_32GB
		res.GPU = GPU_12GB
	}

	return res
}
func getOllamaModels() ([]types.OllamaModelsInfo, error) {
	req, err := http.Get(GetOllamaUrl() + "/api/tags")
	res := []types.OllamaModelsInfo{}
	if err != nil {
		return res, fmt.Errorf("failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return res, fmt.Errorf("failed to read response body")
	}
	rest := types.OllamaModelsList{}
	if err := json.Unmarshal(body, &rest); err != nil {
		return res, fmt.Errorf("failed to unmarshal response body")
	}
	return rest.Models, nil

}
func RefreshOllamaHandler(w http.ResponseWriter, r *http.Request) {
	err := refreshOllamaModels(r)
	if err != nil {
		libs.ErrorMsg(w, "Refresh Ollama Models error")
		return
	}
	//libs.SuccessMsg(w, nil, "Refresh Ollama Models success")
	Tagshandler(w, r)
}
func refreshOllamaModels(r *http.Request) error {
	modelList, err := getOllamaModels()
	if err != nil {
		return fmt.Errorf("load ollama error: %v", err)
	}
	// 将modelList中的数据写入ReqBodyMap
	for _, modelInfo := range modelList {
		model := modelInfo.Model
		if _, exists := config.ReqBodyMap.Load(model); !exists {
			// 创建一个新的ReqBody对象并填充相关信息
			oinfo := parseOllamaInfo(modelInfo)
			details, err := getOllamaInfo(r, model)
			if err != nil {
				log.Printf("Error getting ollama info: %v", err)
				continue
			}
			architecture := details.ModelInfo["general.architecture"].(string)
			contextLength := convertInt(details.ModelInfo, architecture+".context_length")
			embeddingLength := convertInt(details.ModelInfo, architecture+".embedding_length")
			paths, err := getManifests(model)
			if err != nil {
				log.Printf("Error parsing Manifests: %v", err)
				continue
			}
			reqBody := types.ReqBody{
				Model:     model,
				Status:    "success",
				CreatedAt: time.Now(),
			}
			reqBody.Info = types.ModelInfo{
				Engine:          "ollama",
				From:            "ollama",
				Path:            paths,
				Size:            oinfo.Size,
				Quant:           oinfo.Quant,
				Desk:            oinfo.Desk,
				CPU:             oinfo.CPU,
				GPU:             oinfo.GPU,
				Template:        details.Template,
				Parameters:      details.Parameters,
				ContextLength:   contextLength,
				EmbeddingLength: embeddingLength,
			}

			// 将新的ReqBody对象写入ReqBodyMap
			config.ReqBodyMap.Store(model, reqBody)
		}
	}
	return nil
}
func setOllamaInfo(w http.ResponseWriter, r *http.Request, reqBody types.ReqBody) {
	model := reqBody.Model
	postQuery := map[string]interface{}{
		"model": model,
	}
	url := GetOllamaUrl() + "/api/pull"
	ForwardHandler(w, r, postQuery, url, nil, "POST")
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
			contextLength := convertInt(details.ModelInfo, architecture+".context_length")
			embeddingLength := convertInt(details.ModelInfo, architecture+".embedding_length")
			paths, err := getManifests(model.Model)
			if err != nil {
				log.Printf("Error parsing Manifests: %v", err)
				continue
			}

			reqBody.Info = types.ModelInfo{
				Engine:          reqBody.Info.Engine,
				From:            reqBody.Info.From,
				Path:            paths,
				Size:            oinfo.Size,
				Quant:           oinfo.Quant,
				Desk:            oinfo.Desk,
				CPU:             oinfo.CPU,
				GPU:             oinfo.GPU,
				Template:        details.Template,
				Parameters:      details.Parameters,
				ContextLength:   contextLength,
				EmbeddingLength: embeddingLength,
			}
			//reqBody.Paths = paths
			reqBody.Status = "success"
			reqBody.CreatedAt = time.Now()
			if err := config.SetModel(reqBody); err != nil {
				libs.ErrorMsg(w, "Set model error")
				return
			}
			return
		}
	}
}
func convertInt(data map[string]interface{}, str string) int {
	res := 0
	if val, ok := data[str]; ok {
		switch v := val.(type) {
		case int:
			res = v
		case float64:
			res = int(v)
		default:
			log.Printf("Unexpected type for embedding_length: %T", v)
		}
	}
	return res
}
func getOllamaInfo(r *http.Request, model string) (types.OllamaModelDetail, error) {
	infoQuery := map[string]interface{}{
		"name": model,
	}
	res := types.OllamaModelDetail{}
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
func GetOpName(model string) types.OmodelPath {
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
	return types.OmodelPath{
		Space:   "registry.ollama.ai",
		LibPath: libPath,
		Name:    modelName,
		Tag:     modelTags,
	}
}
func getManifests(model string) ([]string, error) {
	res := []string{}
	opName := GetOpName(model)
	modelsDir := GetOllamaModelDir()
	manifestsFile := filepath.Join(modelsDir, "manifests", opName.Space, opName.LibPath, opName.Name, opName.Tag)
	if !libs.PathExists(manifestsFile) {
		return res, fmt.Errorf("failed to get manifests file: %s", manifestsFile)
	}
	res = append(res, manifestsFile)
	var manifest types.ManifestV2
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
	dir := GetOllamaModelDir()
	// only accept actual sha256 digests
	pattern := "^sha256[:-][0-9a-fA-F]{64}$"
	re := regexp.MustCompile(pattern)

	if digest != "" && !re.MatchString(digest) {
		return "", fmt.Errorf("invalid digest format")
	}

	digest = strings.ReplaceAll(digest, ":", "-")
	path := filepath.Join(dir, "blobs", digest)
	return path, nil
}

func ConvertOllama(w http.ResponseWriter, r *http.Request, req types.ReqBody) {
	modelFile := "FROM " + req.Info.Path[0] + "\n"
	modelFile += `TEMPLATE """` + req.Info.Template + `"""`
	if req.Info.Parameters != "" {
		parameters := strings.Split(req.Info.Parameters, "\n")
		for _, param := range parameters {
			modelFile += "\nPARAMETER " + param
		}
	}

	url := GetOllamaUrl() + "/api/create"
	postParams := map[string]string{
		"name":      req.Model,
		"modelfile": modelFile,
	}
	ForwardHandler(w, r, postParams, url, nil, "POST")
	modelDir, err := config.GetModelDir(req.Model)
	if err != nil {
		libs.ErrorMsg(w, "GetModelDir")
		return
	}
	// modelFilePath := filepath.Join(modelDir, "Modelfile")
	// if err := os.WriteFile(modelFilePath, []byte(modelFile), 0644); err != nil {
	// 	ErrMsg("WriteFile", err, w)
	// 	return
	// }
	err = os.RemoveAll(modelDir)
	if err != nil {
		libs.ErrorMsg(w, "Error removing directory")
		return
	}
}

func Var(key string) string {
	return strings.Trim(strings.TrimSpace(os.Getenv(key)), "\"'")
}
func GetOllamaModelDir() string {
	if s := Var("OLLAMA_MODELS"); s != "" {
		return s
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ollama", "models")
}
func GetOllamaUrl() string {
	if s := strings.TrimSpace(Var("OLLAMA_HOST")); s != "" {
		return "http://" + s
	}
	ollamaUrl := libs.GetConfigString("ollamaUrl")
	if ollamaUrl != "" {
		return ollamaUrl
	} else {
		return "http://localhost:11434"
	}
}
func GetModelDir(fileName string, model string) string {
	var filePath string
	dir := GetOllamaModelDir()
	if strings.Contains(fileName, "sha256-") && len(fileName) == 71 {
		filePath = filepath.Join(dir, "blobs", fileName)
		//log.Printf("====filePath1: %s", filePath)
	} else {
		opName := GetOpName(model)
		filePath = filepath.Join(dir, "manifests", opName.Space, opName.LibPath, opName.Name, opName.Tag)
		//log.Printf("====filePath2: %s", filePath)
	}
	return filePath
}
