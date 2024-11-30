package model

import "gorm.io/gorm"

type VectorList struct {
	gorm.Model
	Name            string `json:"name"`
	FilePath        string `json:"file_path"`
	DbPath          string `json:"db_path"`
	Engine          string `json:"engine"`
	EmbeddingUrl    string `json:"embedding_url"`
	EmbeddingApiKey string `json:"embedding_api_key"`
	EmbeddingModel  string `json:"embedding_model"`
}
