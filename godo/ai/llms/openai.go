package llms

import (
	"fmt"
	"godo/libs"
)

// ollama openai deepseek bigmodel alibaba 01ai cloudflare groq mistral anthropic llamafamily
var OpenAIApiMaps = map[string]string{
	//"openai":      GetOpenAIUrl(),
	"deepseek":    "https://api.deepseek.com/v1",
	"bigmodel":    "https://open.bigmodel.cn/api/paas/v4",
	"alibaba":     "https://dashscope.aliyuncs.com/compatible-mode/v1",
	"01ai":        "https://api.lingyiwanwu.com/v1",
	"groq":        "https://api.groq.com/openai/v1",
	"mistral":     "https://api.mistral.ai/v1",
	"anthropic":   "https://api.anthropic.com/v1",
	"llamafamily": "https://api.atomecho.cn/v1",
}

// 获取 OpenAI 聊天 API 的 URL
func GetOpenAIChatUrl(types string) (map[string]string, string, error) {
	aiUrl, err := GetAIUrl(types)
	if err != nil {
		return nil, "", err
	}
	headers, err := GetOpenAIHeaders(types)
	if err != nil {
		return nil, "", err
	}
	return headers, aiUrl + "/chat/completions", nil
}

// 获取 OpenAI 文本嵌入 API 的 URL
func GetOpenAIEmbeddingUrl(types string) (map[string]string, string, error) {
	aiUrl, err := GetAIUrl(types)
	if err != nil {
		return nil, "", err
	}
	headers, err := GetOpenAIHeaders(types)
	if err != nil {
		return nil, "", err
	}
	return headers, aiUrl + "/embeddings", nil
}

// 获取 OpenAI 文本转图像 API 的 URL
func GetOpenAIText2ImgUrl(types string) (map[string]string, string, error) {
	aiUrl, err := GetAIUrl(types)
	if err != nil {
		return nil, "", err
	}
	headers, err := GetOpenAIHeaders(types)
	if err != nil {
		return nil, "", err
	}
	return headers, aiUrl + "/images/generations", nil
}

func GetAIUrl(types string) (string, error) {
	if types == "openai" {
		return GetOpenAIUrl(), nil
	} else if types == "cloudflare" {
		return GetCloudflareUrl()
	} else {
		url, exists := OpenAIApiMaps[types]
		if !exists {
			return "", fmt.Errorf("the " + types + " url is not set")
		} else {
			return url, nil
		}
	}
}
func GetOpenAIHeaders(types string) (map[string]string, error) {
	secret, err := GetOpenAISecret(types)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + secret,
	}, nil
}

// 获取 OpenAI 密钥
func GetOpenAISecret(types string) (string, error) {
	secret, has := libs.GetConfig(types + "Secret")
	if !has {
		return "", fmt.Errorf("the " + types + " secret is not set")
	}
	return secret.(string), nil
}
func GetOpenAIUserId(types string) (string, error) {
	userId, has := libs.GetConfig(types + "UserId")
	if !has {
		return "", fmt.Errorf("the " + types + " user id is not set")
	}
	return userId.(string), nil
}
func GetOpenAIUrl() string {
	openaiUrl, ok := libs.GetConfig("openaiUrl")
	if ok {
		return openaiUrl.(string)
	} else {
		return "https://api.openai.com/v1"
	}
}
func GetCloudflareUrl() (string, error) {
	userId, err := GetOpenAIUserId("cloudflare")
	if err != nil {
		return "", err
	}
	return "https://api.cloudflare.com/client/v4/accounts/" + userId + "/ai/v1", nil
}
