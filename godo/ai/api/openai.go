package api

import (
	"fmt"
	"godo/libs"
)

// 获取 OpenAI 聊天 API 的 URL
func GetOpenAIChatUrl() string {
	return "https://api.openai.com/v1/chat/completions"
}

// 获取 OpenAI 文本嵌入 API 的 URL
func GetOpenAIEmbeddingUrl() string {
	return "https://api.openai.com/v1/embeddings"
}

// 获取 OpenAI 文本转图像 API 的 URL
func GetOpenAIText2ImgUrl() string {
	return "https://api.openai.com/v1/images/generations"
}

// 获取 OpenAI 密钥
func GetOpenAISecret() (string, error) {
	secret, has := libs.GetConfig("openaiSecret")
	if !has {
		return "", fmt.Errorf("the openai secret is not set")
	}
	return secret.(string), nil
}
