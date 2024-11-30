package llms

func GetOllamaChatUrl(url string) string {
	return url + "/v1/chat/completions"
}
func GetOllamaEmbeddingUrl(url string) string {
	return url + "/api/embeddings"
}
