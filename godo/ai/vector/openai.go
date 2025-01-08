package vector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const BaseURLOpenAI = "https://api.openai.com/v1"

type EmbeddingModelOpenAI string

const (
	EmbeddingModelOpenAI2Ada   EmbeddingModelOpenAI = "text-embedding-ada-002"
	EmbeddingModelOpenAI3Small EmbeddingModelOpenAI = "text-embedding-3-small"
	EmbeddingModelOpenAI3Large EmbeddingModelOpenAI = "text-embedding-3-large"
)

type openAIResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// NewEmbeddingFuncDefault 返回一个函数，使用 OpenAI 的 "text-embedding-3-small" 模型通过 API 创建文本嵌入向量。
// 该模型支持的最大文本长度为 8191 个标记。
// API 密钥从环境变量 "OPENAI_API_KEY" 中读取。
func NewEmbeddingFuncDefault() EmbeddingFunc {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return NewEmbeddingFuncOpenAI(apiKey, EmbeddingModelOpenAI3Small)
}

// NewEmbeddingFuncOpenAI 返回一个函数，使用 OpenAI API 创建文本嵌入向量。
func NewEmbeddingFuncOpenAI(apiKey string, model EmbeddingModelOpenAI) EmbeddingFunc {
	// OpenAI 嵌入向量已归一化
	normalized := true
	return NewEmbeddingFuncOpenAICompat(BaseURLOpenAI, apiKey, string(model), &normalized)
}

// NewEmbeddingFuncOpenAICompat 返回一个函数，使用兼容 OpenAI 的 API 创建文本嵌入向量。
// 例如：
//   - Azure OpenAI: https://azure.microsoft.com/en-us/products/ai-services/openai-service
//   - LitLLM: https://github.com/BerriAI/litellm
//   - Ollama: https://github.com/ollama/ollama/blob/main/docs/openai.md
//
// `normalized` 参数表示嵌入模型返回的向量是否已经归一化。如果为 nil，则会在首次请求时自动检测（有小概率向量恰好长度为 1）。
func NewEmbeddingFuncOpenAICompat(baseURL, apiKey, model string, normalized *bool) EmbeddingFunc {
	client := &http.Client{}

	var checkedNormalized bool
	checkNormalized := sync.Once{}

	return func(ctx context.Context, text string) ([]float32, error) {
		// 准备请求体
		reqBody, err := json.Marshal(map[string]string{
			"input": text,
			"model": model,
		})
		if err != nil {
			return nil, fmt.Errorf("无法序列化请求体: %w", err)
		}

		// 创建带有上下文的请求以支持超时
		req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/embeddings", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("无法创建请求: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("无法发送请求: %w", err)
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("嵌入 API 返回错误响应: " + resp.Status)
		}

		// 读取并解码响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("无法读取响应体: %w", err)
		}
		var embeddingResponse openAIResponse
		err = json.Unmarshal(body, &embeddingResponse)
		if err != nil {
			return nil, fmt.Errorf("无法反序列化响应体: %w", err)
		}

		// 检查响应中是否包含嵌入向量
		if len(embeddingResponse.Data) == 0 || len(embeddingResponse.Data[0].Embedding) == 0 {
			return nil, errors.New("响应中未找到嵌入向量")
		}

		v := embeddingResponse.Data[0].Embedding
		if normalized != nil {
			if *normalized {
				return v, nil
			}
			return normalizeVector(v), nil
		}
		checkNormalized.Do(func() {
			if isNormalized(v) {
				checkedNormalized = true
			} else {
				checkedNormalized = false
			}
		})
		if !checkedNormalized {
			v = normalizeVector(v)
		}

		return v, nil
	}
}
