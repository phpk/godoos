package server

import (
	"encoding/json"
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
	headers, url, err := GetHeadersAndUrl(req, "chat")
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	ForwardHandler(w, r, req, url, headers, "POST")
}

func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	headers, url, err := GetHeadersAndUrl(req, "embeddings")
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	ForwardHandler(w, r, req, url, headers, "POST")
}
