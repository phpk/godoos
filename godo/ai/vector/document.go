package vector

import (
	"context"
	"errors"
	"fmt"
)

// Document 表示单个文档。
type Document struct {
	ID        string            // 文档的唯一标识符
	Metadata  map[string]string // 文档的元数据
	Embedding []float32         // 文档的嵌入向量
	Content   string            // 文档的内容

	// ⚠️ 当在此处添加未导出字段时，请考虑在 [DB.Export] 和 [DB.Import] 中添加一个持久化结构版本。
}

// NewDocument 创建一个新的文档，包括其嵌入向量。
// 元数据是可选的。
// 如果未提供嵌入向量，则使用嵌入函数创建。
// 如果内容为空但需要存储嵌入向量，可以仅提供嵌入向量。
// 如果 embeddingFunc 为 nil，则使用默认的嵌入函数。
//
// 如果你想创建没有嵌入向量的文档，例如让 [Collection.AddDocuments] 并发创建它们，
// 可以使用 `chromem.Document{...}` 而不是这个构造函数。
func NewDocument(ctx context.Context, id string, metadata map[string]string, embedding []float32, content string, embeddingFunc EmbeddingFunc) (Document, error) {
	if id == "" {
		return Document{}, errors.New("ID 不能为空")
	}
	if len(embedding) == 0 && content == "" {
		return Document{}, errors.New("嵌入向量或内容必须至少有一个非空")
	}
	if embeddingFunc == nil {
		embeddingFunc = NewEmbeddingFuncDefault()
	}

	if len(embedding) == 0 {
		var err error
		embedding, err = embeddingFunc(ctx, content)
		if err != nil {
			return Document{}, fmt.Errorf("无法生成嵌入向量: %w", err)
		}
	}

	return Document{
		ID:        id,
		Metadata:  metadata,
		Embedding: embedding,
		Content:   content,
	}, nil
}
