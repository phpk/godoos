package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"godo/ai/types"
	"godo/libs"
	"io"
	"log"
	"net/http"
)

func SendChat(w http.ResponseWriter, r *http.Request, reqBody interface{}, url string, headers map[string]string) (types.OpenAIResponse, error) {
	var res types.OpenAIResponse
	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		return res, fmt.Errorf("Error marshaling payload")
	}
	// 创建POST请求，复用原始请求的上下文（如Cookies）
	req, err := http.NewRequestWithContext(r.Context(), "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return res, fmt.Errorf("Failed to create request")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	//req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, fmt.Errorf("Failed to send request")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("Failed to decode response body")
	}
	return res, nil
}
func ForwardHandler(w http.ResponseWriter, r *http.Request, reqBody interface{}, url string, headers map[string]string, method string) {
	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		libs.ErrorMsg(w, "Error marshaling payload")
		return
	}
	// 创建POST请求，复用原始请求的上下文（如Cookies）
	req, err := http.NewRequestWithContext(r.Context(), method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		libs.ErrorMsg(w, "Failed to create request")
		return
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	//req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		libs.ErrorMsg(w, "Failed to send request")
		return
	}
	defer resp.Body.Close()
	// 将外部服务的响应内容原封不动地转发给客户端
	for k, v := range resp.Header {
		for _, value := range v {
			w.Header().Add(k, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	//log.Printf("resp.Body: %v", resp.Body)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// 如果Copy过程中出错，尝试发送一个错误响应给客户端
		http.Error(w, "Error forwarding response", http.StatusInternalServerError)
		log.Printf("Error forwarding response body: %v", err)
		return
	}

}
