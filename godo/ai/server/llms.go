package server

import (
	"fmt"
	"godo/ai/types"
	"godo/libs"
	"log"
)

// ollama openai deepseek bigmodel alibaba 01ai cloudflare groq mistral anthropic llamafamily
var OpenAIApiMaps = map[string]string{
	"ollama":      "",
	"openai":      "",
	"gitee":       "",
	"cloudflare":  "",
	"deepseek":    "https://api.deepseek.com/v1",
	"volces":      "https://ark.cn-beijing.volces.com/api/v3",
	"bigmodel":    "https://open.bigmodel.cn/api/paas/v4",
	"alibaba":     "https://dashscope.aliyuncs.com/compatible-mode/v1",
	"01ai":        "https://api.lingyiwanwu.com/v1",
	"groq":        "https://api.groq.com/openai/v1",
	"mistral":     "https://api.mistral.ai/v1",
	"anthropic":   "https://api.anthropic.com/v1",
	"llamafamily": "https://api.atomecho.cn/v1",
	"siliconflow": "https://api.siliconflow.cn/v1",
}

func GetHeadersAndUrl(req types.ChatRequest, chattype string) (map[string]string, string, error) {
	// engine, ok := req["engine"].(string)
	// if !ok {
	// 	return nil, "", fmt.Errorf("invalid engine field in request")
	// }
	// model, ok := req["model"].(string)
	// if !ok {
	// 	return nil, "", fmt.Errorf("invalid model field in request")
	// }
	// 获取url
	url, has := OpenAIApiMaps[req.Engine]
	if !has {
		return nil, "", fmt.Errorf("invalid engine map in request")
	}
	var err error
	if url == "" {
		if req.Engine == "openai" {
			url = GetOpenAIUrl()
		} else if req.Engine == "cloudflare" {
			url, err = GetCloudflareUrl()
			if err != nil {
				return nil, "", err
			}
		} else if req.Engine == "gitee" {
			url = GetGiteeUrl(req.Model, chattype)
		} else if req.Engine == "ollama" {
			url = GetOllamaUrl() + "/v1"
			log.Printf("get ollama url is %v", url)
		}
	}

	headers, err := GetOpenAIHeaders(req.Engine)
	if err != nil {
		return nil, "", err
	}
	typeUrl := "/chat/completions"
	if chattype == "embeddings" {
		typeUrl = "/embeddings"
	} else if chattype == "text2img" {
		if req.Engine == "gitee" {
			typeUrl = "/text-to-image"
		} else {
			typeUrl = "/images/generations"
		}

	} else if chattype == "text2voice" {

	} else if chattype == "voice2text" {

	}
	return headers, url + typeUrl, nil

}

func GetOpenAIHeaders(types string) (map[string]string, error) {
	if types == "ollama" {
		return map[string]string{
			"Content-Type": "application/json",
		}, nil
	}
	secret, err := GetOpenAISecret(types)
	if types == "gitee" {
		return map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + secret,
			"X-Package":     "1910",
		}, nil
	}
	if err != nil {
		return map[string]string{
			"Content-Type": "application/json",
		}, nil
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
func GetGiteeUrl(model string, chatType string) string {
	return "https://ai.gitee.com/api/serverless/" + model
}
