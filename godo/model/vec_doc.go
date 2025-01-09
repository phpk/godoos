package model

import (
	"fmt"
	"sort"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"gorm.io/gorm"
)

type VecDoc struct {
	gorm.Model
	Content  string `json:"content"`
	FilePath string `json:"file_path" gorm:"not null"`
	ListID   uint   `json:"list_id"`
}

func (VecDoc) TableName() string {
	return "vec_doc"
}

func Adddocument(listId uint, docs []VecDoc, embeds [][]float32) error {
	// 批量删除具有相同 file_path 的 VecDoc 数据
	filePath := docs[0].FilePath
	// 删除向量表中的数据
	var existingDocs []VecDoc
	if err := Db.Where("file_path = ?", filePath).Find(&existingDocs).Error; err != nil {
		return fmt.Errorf("failed to find existing documents: %v", err)
	}

	for _, existingDoc := range existingDocs {
		documentID := fmt.Sprintf("%d", existingDoc.ID)
		result := Db.Exec(fmt.Sprintf("DELETE FROM [%d_vec] WHERE document_id = ?", listId), documentID)
		if result.Error != nil {
			return fmt.Errorf("failed to delete vector data: %v", result.Error)
		}
	}

	// 删除 vec_doc 中的数据（硬删除）
	if err := Db.Unscoped().Where("file_path = ?", filePath).Delete(&VecDoc{}).Error; err != nil {
		return fmt.Errorf("failed to delete vec_doc data: %v", err)
	}
	// 批量插入新的 vec_doc 数据
	if err := Db.CreateInBatches(docs, 100).Error; err != nil {
		return fmt.Errorf("failed to create vec_doc data: %v", err)
	}

	// 批量插入向量数据到虚拟表
	for i, doc := range docs {
		v, err := sqlite_vec.SerializeFloat32(embeds[i])
		if err != nil {
			return fmt.Errorf("failed to serialize vector: %v", err)
		}

		documentID := fmt.Sprintf("%d", doc.ID)
		result := Db.Exec(fmt.Sprintf("INSERT INTO [%d_vec] (document_id, embedding) VALUES (?, ?)", listId), documentID, v)
		if result.Error != nil {
			return fmt.Errorf("failed to insert vector data: %v", result.Error)
		}
	}

	return nil
}

func Deletedocument(listId uint, filePath string) error {
	var existingDocs []VecDoc
	if err := Db.Where("file_path = ?", filePath).Find(&existingDocs).Error; err != nil {
		return fmt.Errorf("failed to find existing documents: %v", err)
	}

	for _, existingDoc := range existingDocs {
		documentID := fmt.Sprintf("%d", existingDoc.ID)
		result := Db.Exec(fmt.Sprintf("DELETE FROM [%d_vec] WHERE document_id = ?", listId), documentID)
		if result.Error != nil {
			return fmt.Errorf("failed to delete vector data: %v", result.Error)
		}
	}

	// 删除 vec_doc 中的数据（硬删除）
	if err := Db.Unscoped().Where("file_path = ?", filePath).Delete(&VecDoc{}).Error; err != nil {
		return fmt.Errorf("failed to delete vec_doc data: %v", err)
	}

	return nil
}

type AskDocResponse struct {
	Content  string  `json:"content"`
	Score    float32 `json:"score"`
	FilePath string  `json:"file_path"`
}

func AskDocument(listId uint, query []float32) ([]AskDocResponse, error) {
	// 序列化查询向量
	queryVec, err := sqlite_vec.SerializeFloat32(query)
	if err != nil {
		return []AskDocResponse{}, fmt.Errorf("failed to serialize query vector: %v", err)
	}

	// 查询最相似的文档
	var results []struct {
		DocumentID uint    `gorm:"column:document_id"`
		Distance   float32 `gorm:"column:distance"`
	}
	result := Db.Raw(fmt.Sprintf(`
		SELECT
			document_id,
			distance
		FROM [%d_vec]
		WHERE embedding MATCH ?
		ORDER BY distance
		LIMIT 10
	`, listId), queryVec).Scan(&results)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to query vector data: %v", result.Error)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no matching documents found")
	}

	// 获取最相似的文档
	var docs []VecDoc
	var docIDs []uint
	for _, res := range results {
		docIDs = append(docIDs, res.DocumentID)
	}

	if err := Db.Where("id IN ?", docIDs).Find(&docs).Error; err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}

	// 构建响应
	var responses []AskDocResponse
	for i, doc := range docs {
		responses = append(responses, AskDocResponse{
			Content:  doc.Content,
			Score:    results[i].Distance,
			FilePath: doc.FilePath,
		})
	}
	// 按 Score 降序排序
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Score > responses[j].Score
	})
	return responses, nil
}
