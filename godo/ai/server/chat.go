package server

import (
	"encoding/json"
	"godo/ai/llms"
	"godo/libs"
	"net/http"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// url := GetOllamaUrl() + "/v1/chat/completions"
	var url string
	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	engine, ok := req["engine"].(string)
	if !ok {
		libs.ErrorMsg(w, "Invalid engine field in request")
		return
	}
	model, ok := req["model"].(string)
	if !ok {
		libs.ErrorMsg(w, "Invalid model field in request")
		return
	}
	var headers map[string]string
	switch engine {
	case "ollama":
		ollamaUrl := GetOllamaUrl()
		url = llms.GetOllamaChatUrl(ollamaUrl)
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	case "gitee":
		url = llms.GetGiteeChatUrl(model)
	case "openai":
		headers, url, err = llms.GetOpenAIChatUrl("openai")
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	case "cloudflare":
		headers, url, err = llms.GetOpenAIChatUrl("cloudflare")
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	default:
		headers, url, err = llms.GetOpenAIChatUrl("openai") // 默认URL
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	}
	ForwardHandler(w, r, req, url, headers, "POST")
}

func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	url := GetOllamaUrl() + "/api/embeddings"
	var request interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	ForwardHandler(w, r, request, url, nil, "POST")
}
