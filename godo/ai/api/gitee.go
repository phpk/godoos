package api

import (
	"fmt"
	"godo/libs"
)

func GetGiteeChatUrl(model string) string {
	return "https://ai.gitee.com/api/serverless/" + model + "/chat/completions"
}
func GetGiteeEmbeddingUrl(model string) string {
	return "https://ai.gitee.com/api/serverless/" + model + "/embeddings"
}
func GetGiteeText2ImgUrl(model string) string {
	return "https://ai.gitee.com/api/serverless/" + model + "/text-to-image"
}
func GetGiteeSecret() (string, error) {
	secret, has := libs.GetConfig("giteeSecret")
	if !has {
		return "", fmt.Errorf("the gitee secret is not set")
	}
	return secret.(string), nil
}
