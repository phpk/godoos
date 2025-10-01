package vector

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
)

// Collection 表示一个文档集合。
// 它还包含一个配置好的嵌入函数，当添加没有嵌入的文档时会使用该函数。
type Collection struct {
	Name string

	metadata      map[string]string
	documents     map[string]*Document
	documentsLock sync.RWMutex
	embed         EmbeddingFunc

	persistDirectory string
	compress         bool

	// ⚠️ 当添加字段时，请考虑在 [DB.Export] 和 [DB.Import] 的持久化结构中添加相应的字段
}

// 我们不导出这个函数，以保持 API 表面最小。
// 用户通过 [Client.CreateCollection] 创建集合。
func newCollection(name string, metadata map[string]string, embed EmbeddingFunc, dbDir string, compress bool) (*Collection, error) {
	// 复制元数据以避免在创建集合后调用者修改元数据时发生数据竞争。
	m := make(map[string]string, len(metadata))
	for k, v := range metadata {
		m[k] = v
	}

	c := &Collection{
		Name: name,

		metadata:  m,
		documents: make(map[string]*Document),
		embed:     embed,
	}

	// 持久化
	if dbDir != "" {
		safeName := hash2hex(name)
		c.persistDirectory = filepath.Join(dbDir, safeName)
		c.compress = compress
		// 持久化名称和元数据
		metadataPath := filepath.Join(c.persistDirectory, metadataFileName)
		metadataPath += ".gob"
		if c.compress {
			metadataPath += ".gz"
		}
		pc := struct {
			Name     string
			Metadata map[string]string
		}{
			Name:     name,
			Metadata: m,
		}
		err := persistToFile(metadataPath, pc, compress, "")
		if err != nil {
			return nil, fmt.Errorf("无法持久化集合元数据: %w", err)
		}
	}

	return c, nil
}

// 添加嵌入到数据存储中。
//
//   - ids: 要添加的嵌入的 ID
//   - embeddings: 要添加的嵌入。如果为 nil，则基于内容使用集合的嵌入函数计算嵌入。可选。
//   - metadatas: 与嵌入关联的元数据。查询时可以过滤这些元数据。可选。
//   - contents: 与嵌入关联的内容。
//
// 这是一个类似于 Chroma 的方法。对于更符合 Go 风格的方法，请参见 [AddDocuments]。
func (c *Collection) Add(ctx context.Context, ids []string, embeddings [][]float32, metadatas []map[string]string, contents []string) error {
	return c.AddConcurrently(ctx, ids, embeddings, metadatas, contents, 1)
}

// AddConcurrently 类似于 Add，但并发地添加嵌入。
// 这在没有传递任何嵌入时特别有用，因为需要创建嵌入。
// 出现错误时，取消所有并发操作并返回错误。
//
// 这是一个类似于 Chroma 的方法。对于更符合 Go 风格的方法，请参见 [AddDocuments]。
func (c *Collection) AddConcurrently(ctx context.Context, ids []string, embeddings [][]float32, metadatas []map[string]string, contents []string, concurrency int) error {
	if len(ids) == 0 {
		return errors.New("ids 为空")
	}
	if len(embeddings) == 0 && len(contents) == 0 {
		return errors.New("必须填写 embeddings 或 contents")
	}
	if len(embeddings) != 0 {
		if len(embeddings) != len(ids) {
			return errors.New("ids 和 embeddings 的长度必须相同")
		}
	} else {
		// 分配空切片以便稍后通过索引访问
		embeddings = make([][]float32, len(ids))
	}
	if len(metadatas) != 0 {
		if len(ids) != len(metadatas) {
			return errors.New("当 metadatas 不为空时，其长度必须与 ids 相同")
		}
	} else {
		// 分配空切片以便稍后通过索引访问
		metadatas = make([]map[string]string, len(ids))
	}
	if len(contents) != 0 {
		if len(contents) != len(ids) {
			return errors.New("ids 和 contents 的长度必须相同")
		}
	} else {
		// 分配空切片以便稍后通过索引访问
		contents = make([]string, len(ids))
	}
	if concurrency < 1 {
		return errors.New("并发数必须至少为 1")
	}

	// 将 Chroma 风格的参数转换为文档切片
	docs := make([]Document, 0, len(ids))
	for i, id := range ids {
		docs = append(docs, Document{
			ID:        id,
			Metadata:  metadatas[i],
			Embedding: embeddings[i],
			Content:   contents[i],
		})
	}

	return c.AddDocuments(ctx, docs, concurrency)
}

// AddDocuments 使用指定的并发数将文档添加到集合中。
// 如果文档没有嵌入，则使用集合的嵌入函数创建嵌入。
// 出现错误时，取消所有并发操作并返回错误。
func (c *Collection) AddDocuments(ctx context.Context, documents []Document, concurrency int) error {
	if len(documents) == 0 {
		// TODO: 这是否应为无操作（no-op）？
		return errors.New("documents 切片为空")
	}
	if concurrency < 1 {
		return errors.New("并发数必须至少为 1")
	}
	// 对于其他验证，我们依赖于 AddDocument。

	var sharedErr error
	sharedErrLock := sync.Mutex{}
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	setSharedErr := func(err error) {
		sharedErrLock.Lock()
		defer sharedErrLock.Unlock()
		// 另一个 goroutine 可能已经设置了错误。
		if sharedErr == nil {
			sharedErr = err
			// 取消所有其他 goroutine 的操作。
			cancel(sharedErr)
		}
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)
	for _, doc := range documents {
		wg.Add(1)
		go func(doc Document) {
			defer wg.Done()

			// 如果另一个 goroutine 已经失败，则不开始。
			if ctx.Err() != nil {
				return
			}

			// 等待直到 $concurrency 个其他 goroutine 正在创建文档。
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := c.AddDocument(ctx, doc)
			if err != nil {
				setSharedErr(fmt.Errorf("无法添加文档 '%s': %w", doc.ID, err))
				return
			}
		}(doc)
	}

	wg.Wait()

	return sharedErr
}

// AddDocument 将文档添加到集合中。
// 如果文档没有嵌入，则使用集合的嵌入函数创建嵌入。
func (c *Collection) AddDocument(ctx context.Context, doc Document) error {
	if doc.ID == "" {
		return errors.New("文档 ID 为空")
	}
	if len(doc.Embedding) == 0 && doc.Content == "" {
		return errors.New("必须填写文档的 embedding 或 content")
	}

	// 复制元数据以避免在创建文档后调用者修改元数据时发生数据竞争。
	m := make(map[string]string, len(doc.Metadata))
	for k, v := range doc.Metadata {
		m[k] = v
	}

	// 如果嵌入不存在，则创建嵌入，否则如果需要则规范化
	if len(doc.Embedding) == 0 {
		embedding, err := c.embed(ctx, doc.Content)
		if err != nil {
			return fmt.Errorf("无法创建文档的嵌入: %w", err)
		}
		doc.Embedding = embedding
	} else {
		if !isNormalized(doc.Embedding) {
			doc.Embedding = normalizeVector(doc.Embedding)
		}
	}

	c.documentsLock.Lock()
	// 我们不使用 defer 解锁，因为我们希望尽早解锁。
	c.documents[doc.ID] = &doc
	c.documentsLock.Unlock()

	// 持久化文档
	if c.persistDirectory != "" {
		docPath := c.getDocPath(doc.ID)
		err := persistToFile(docPath, doc, c.compress, "")
		if err != nil {
			return fmt.Errorf("无法将文档持久化到 %q: %w", docPath, err)
		}
	}

	return nil
}

// Delete 从集合中删除文档。
//
//   - where: 元数据的条件过滤。可选。
//   - whereDocument: 文档的条件过滤。可选。
//   - ids: 要删除的文档的 ID。如果为空，则删除所有文档。
func (c *Collection) Delete(_ context.Context, where, whereDocument map[string]string, ids ...string) error {
	// 必须至少有一个 where、whereDocument 或 ids
	if len(where) == 0 && len(whereDocument) == 0 && len(ids) == 0 {
		return fmt.Errorf("必须至少有一个 where、whereDocument 或 ids")
	}

	if len(c.documents) == 0 {
		return nil
	}

	for k := range whereDocument {
		// 替换 slices.Contains 为手动实现
		if !containsString(supportedFilters, k) {
			return errors.New("不支持的 whereDocument 操作符")
		}
	}

	var docIDs []string

	c.documentsLock.Lock()
	defer c.documentsLock.Unlock()

	if where != nil || whereDocument != nil {
		// 元数据 + 内容过滤
		filteredDocs := filterDocs(c.documents, where, whereDocument)
		for _, doc := range filteredDocs {
			docIDs = append(docIDs, doc.ID)
		}
	} else {
		docIDs = ids
	}

	// 如果没有剩余的文档，则不执行操作
	if len(docIDs) == 0 {
		return nil
	}

	for _, docID := range docIDs {
		delete(c.documents, docID)

		// 从磁盘删除文档
		if c.persistDirectory != "" {
			docPath := c.getDocPath(docID)
			err := removeFile(docPath)
			if err != nil {
				return fmt.Errorf("无法删除文档 %q: %w", docPath, err)
			}
		}
	}

	return nil
}

// Count 返回集合中的文档数量。
func (c *Collection) Count() int {
	c.documentsLock.RLock()
	defer c.documentsLock.RUnlock()
	return len(c.documents)
}

// Result 表示查询结果中的单个结果。
type Result struct {
	ID        string
	Metadata  map[string]string
	Embedding []float32
	Content   string

	// 查询与文档之间的余弦相似度。
	// 值越高，文档与查询越相似。
	// 值的范围是 [-1, 1]。
	Similarity float32
}

// 在集合上执行详尽的最近邻搜索。
//
//   - queryText: 要搜索的文本。其嵌入将使用集合的嵌入函数创建。
//   - nResults: 要返回的结果数量。必须大于 0。
//   - where: 元数据的条件过滤。可选。
//   - whereDocument: 文档的条件过滤。可选。
func (c *Collection) Query(ctx context.Context, queryText string, nResults int, where, whereDocument map[string]string) ([]Result, error) {
	if queryText == "" {
		return nil, errors.New("queryText 为空")
	}

	queryVectors, err := c.embed(ctx, queryText)
	if err != nil {
		return nil, fmt.Errorf("无法创建查询的嵌入: %w", err)
	}

	return c.QueryEmbedding(ctx, queryVectors, nResults, where, whereDocument)
}

// 在集合上执行详尽的最近邻搜索。
//
//   - queryEmbedding: 要搜索的查询的嵌入。必须使用与集合中文档嵌入相同的嵌入模型创建。
//   - nResults: 要返回的结果数量。必须大于 0。
//   - where: 元数据的条件过滤。可选。
//   - whereDocument: 文档的条件过滤。可选。
func (c *Collection) QueryEmbedding(ctx context.Context, queryEmbedding []float32, nResults int, where, whereDocument map[string]string) ([]Result, error) {
	if len(queryEmbedding) == 0 {
		return nil, errors.New("queryEmbedding 为空")
	}
	if nResults <= 0 {
		return nil, errors.New("nResults 必须大于 0")
	}
	c.documentsLock.RLock()
	defer c.documentsLock.RUnlock()
	// if nResults > len(c.documents) {
	// 	return nil, errors.New("nResults 必须小于或等于集合中的文档数量")
	// }

	if len(c.documents) == 0 {
		return nil, nil
	}

	// 验证 whereDocument 操作符
	for k := range whereDocument {
		// 替换 slices.Contains 为手动实现
		if !containsString(supportedFilters, k) {
			return nil, errors.New("不支持的操作符")
		}
	}

	// 根据元数据和内容过滤文档
	filteredDocs := filterDocs(c.documents, where, whereDocument)

	// 如果过滤器删除了所有文档，则不继续
	if len(filteredDocs) == 0 {
		return nil, nil
	}

	// 对于剩余的文档，获取最相似的文档。
	nMaxDocs, err := getMostSimilarDocs(ctx, queryEmbedding, filteredDocs, nResults)
	if err != nil {
		return nil, fmt.Errorf("无法获取最相似的文档: %w", err)
	}
	length := len(nMaxDocs)
	if length > nResults {
		length = nResults
	}
	res := make([]Result, 0, length)
	for i := 0; i < length; i++ {
		doc := c.documents[nMaxDocs[i].docID]
		res = append(res, Result{
			ID:         nMaxDocs[i].docID,
			Metadata:   doc.Metadata,
			Embedding:  doc.Embedding,
			Content:    doc.Content,
			Similarity: nMaxDocs[i].similarity,
		})
	}

	// 返回前 nResults 个结果
	return res, nil
}

// getDocPath 生成文档文件的路径。
func (c *Collection) getDocPath(docID string) string {
	safeID := hash2hex(docID)
	docPath := filepath.Join(c.persistDirectory, safeID)
	docPath += ".gob"
	if c.compress {
		docPath += ".gz"
	}
	return docPath
}

// containsString 检查字符串切片中是否包含指定的字符串
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
