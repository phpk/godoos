package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetEmbeddings(engine string, model string, text []string) ([][]float32, error) {
	if engine == "ollama" {
		res, err := getOllamaEmbedding(model, text)
		if err != nil {
			log.Printf("getOllamaEmbedding error: %v", err)
			return nil, err
		}
		return res, nil
	}
	return nil, nil
}

func getOllamaEmbedding(model string, text []string) ([][]float32, error) {
	url := GetOllamaUrl() + "/api/embed"
	//log.Printf("url: %s", url)
	client := &http.Client{}

	// Prepare the request body.
	reqBody, err := json.Marshal(map[string]interface{}{
		"model": model,
		"input": text,
	})
	//log.Printf("reqBody: %s", reqBody)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from the embedding API: %s", resp.Status)
	}

	// Read and decode the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	var embeddingResponse struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
	}
	//log.Printf("Embedding: %v", embeddingResponse.Embeddings)

	// Return the embeddings directly.
	if len(embeddingResponse.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned from the API")
	}
	return embeddingResponse.Embeddings, nil
}
