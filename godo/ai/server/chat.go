package server

import (
	"encoding/json"
	"fmt"
	"godo/ai/search"
	"godo/libs"
	"godo/office"
	"log"
	"net/http"
	"time"
)

type ChatRequest struct {
	Model       string                 `json:"model"`
	Engine      string                 `json:"engine"`
	Stream      bool                   `json:"stream"`
	WebSearch   bool                   `json:"webSearch"`
	FileContent string                 `json:"fileContent"`
	FileName    string                 `json:"fileName"`
	Options     map[string]interface{} `json:"options"`
	Messages    []Message              `json:"messages"`
}

type Message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var url string
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	if req.WebSearch {
		err = ChatWithWeb(&req)
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	}
	if req.FileContent != "" {
		err = ChatWithFile(&req)
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
	}
	headers, url, err := GetHeadersAndUrl(req, "chat")
	// log.Printf("url: %s", url)
	// log.Printf("headers: %v", headers)
	if err != nil {
		libs.ErrorMsg(w, "the chat request header or url errors:"+err.Error())
		return
	}
	ForwardHandler(w, r, req, url, headers, "POST")
}
func ChatWithFile(req *ChatRequest) error {
	fileContent, err := office.ProcessBase64File(req.FileContent, req.FileName)
	if err != nil {
		return err
	}
	lastMessage, err := GetLastMessage(*req)
	if err != nil {
		return err
	}
	userQuestion := fmt.Sprintf("请对\n%s\n的内容进行分析，给出对用户输入的回答: %s", fileContent, lastMessage)
	log.Printf("the search file is %v", userQuestion)
	req.Messages = append(req.Messages, Message{Role: "user", Content: userQuestion})
	return nil
}
func ChatWithWeb(req *ChatRequest) error {
	lastMessage, err := GetLastMessage(*req)
	if err != nil {
		return err
	}
	searchRequest := search.SearchWeb(lastMessage)
	if len(searchRequest) == 0 {
		return fmt.Errorf("the search web is empty")
	}
	var inputPrompt string
	for _, search := range searchRequest {
		inputPrompt += fmt.Sprintf("- 标题: %s\n- 内容: %s\n", search.Title, search.Content)
	}
	currentDate := time.Now().Format("2006-01-02")
	searchPrompt := fmt.Sprintf(`
# 以下是来自互联网的信息：
%s

# 当前日期: %s

# 要求：
根据最新发布的信息回答用户问题。

# 用户问题：%s

`, inputPrompt, currentDate, lastMessage)
	//log.Printf("the search web is %v", searchPrompt)
	// req.Messages = append([]Message{}, Message{Role: "assistant", Content: searchPrompt})
	req.Messages = append([]Message{}, Message{Role: "user", Content: searchPrompt})
	return nil
}
func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	headers, url, err := GetHeadersAndUrl(req, "embeddings")
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	ForwardHandler(w, r, req, url, headers, "POST")
}
func GetLastMessage(req ChatRequest) (string, error) {
	if len(req.Messages) == 0 {
		return "", fmt.Errorf("the messages is empty")
	}
	lastMessage := req.Messages[len(req.Messages)-1]
	if lastMessage.Role != "user" {
		return "", fmt.Errorf("the last message is not user")
	}
	return lastMessage.Content, nil
}
