package server

import (
	"encoding/json"
	"fmt"
	"godo/ai/search"
	"godo/ai/types"
	"godo/libs"
	"godo/model"
	"godo/office"
	"log"
	"net/http"
	"time"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var url string
	var req types.ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	headers, url, err := GetHeadersAndUrl(req, "chat")
	// log.Printf("url: %s", url)
	// log.Printf("headers: %v", headers)
	if err != nil {
		libs.ErrorMsg(w, "the chat request header or url errors:"+err.Error())
		return
	}
	if req.WebSearch {
		searchRes, err := ChatWithWeb(&req)
		if err != nil {
			log.Printf("the chat with web error:%v", err)
			libs.ErrorMsg(w, err.Error())
			return
		}
		res, err := SendChat(w, r, req, url, headers)
		if err != nil {
			log.Printf("the chat with web error:%v", err)
			libs.ErrorMsg(w, err.Error())
			return
		}
		for _, s := range searchRes {
			res.WebSearch = append(res.WebSearch, types.WebSearchResult{Title: s.Title, Content: s.Content, Link: s.Url})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.FileContent != "" {
		err = ChatWithFile(&req)
		if err != nil {
			log.Printf("the chat with file error:%v", err)
		}
		res, err := SendChat(w, r, req, url, headers)
		if err != nil {
			log.Printf("the chat with web error:%v", err)
			libs.ErrorMsg(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.KnowledgeId != 0 {
		resk, err := ChatWithKnowledge(&req)
		if err != nil {
			log.Printf("the chat with knowledge error:%v", err)
		}
		res, err := SendChat(w, r, req, url, headers)
		if err != nil {
			log.Printf("the chat with web error:%v", err)
			libs.ErrorMsg(w, err.Error())
			return
		}

		for _, s := range resk {
			// s.FilePath = strings.TrimPrefix(s.FilePath, basePath)
			res.Documents = append(res.Documents, types.AskDocResponse{Content: s.Content, Score: s.Score, FilePath: s.FilePath, FileName: s.FileName})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	ForwardHandler(w, r, req, url, headers, "POST")
}
func ChatWithFile(req *types.ChatRequest) error {
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
	req.Messages = append([]types.Message{}, types.Message{Role: "user", Content: userQuestion})
	return nil
}
func ChatWithKnowledge(req *types.ChatRequest) ([]types.AskDocResponse, error) {
	var res []types.AskDocResponse
	lastMessage, err := GetLastMessage(*req)
	if err != nil {
		return res, err
	}
	askrequest := model.AskRequest{
		ID:    req.KnowledgeId,
		Input: lastMessage,
	}
	var knowData model.VecList
	if err := model.Db.First(&knowData, askrequest.ID).Error; err != nil {
		return res, fmt.Errorf("the knowledge id is not exist")
	}
	//var filterDocs
	filterDocs := []string{askrequest.Input}
	// 获取嵌入向量
	resList, err := GetEmbeddings(knowData.Engine, knowData.EmbeddingModel, filterDocs)
	if err != nil {
		return res, fmt.Errorf("the embeddings get error:%v", err)
	}
	resk, err := model.AskDocument(askrequest.ID, resList[0])
	if err != nil {
		return res, fmt.Errorf("the ask document error:%v", err)
	}
	msg := ""
	for _, item := range resk {
		msg += fmt.Sprintf("- %s\n", item.Content)
		res = append(res, types.AskDocResponse{Content: item.Content, Score: item.Score, FilePath: item.FilePath, FileName: item.FileName})
	}
	prompt := fmt.Sprintf(`从文档\n\"\"\"\n%s\n\"\"\"\n中找问题\n\"\"\"\n%s\n\"\"\"\n的答案，找到答案就使用文档语句回答问题，找不到答案就用自身知识回答并且告诉用户该信息不是来自文档。\n不要复述问题，直接开始回答。`, msg, lastMessage)
	req.Messages = append([]types.Message{}, types.Message{Role: "user", Content: prompt})
	return res, nil
}
func ChatWithWeb(req *types.ChatRequest) ([]search.Entity, error) {
	lastMessage, err := GetLastMessage(*req)
	var searchRequest []search.Entity
	if err != nil {
		return searchRequest, err
	}
	searchRequest = search.SearchWeb(lastMessage)
	if len(searchRequest) == 0 {
		return searchRequest, fmt.Errorf("the search web is empty")
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
	req.Messages = append([]types.Message{}, types.Message{Role: "user", Content: searchPrompt})
	return searchRequest, nil
}
func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
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
func GetLastMessage(req types.ChatRequest) (string, error) {
	if len(req.Messages) == 0 {
		return "", fmt.Errorf("the messages is empty")
	}
	lastMessage := req.Messages[len(req.Messages)-1]
	if lastMessage.Role != "user" {
		return "", fmt.Errorf("the last message is not user")
	}
	return lastMessage.Content, nil
}
