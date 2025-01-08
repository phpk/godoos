package vector

import (
	"context"
	"errors"
	"fmt"
	"godo/libs"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// EmbeddingFunc 是一个为给定文本创建嵌入的函数。
// 默认使用 OpenAI 的 "text-embedding-3-small" 模型。
// 该函数必须返回一个已归一化的向量。
type EmbeddingFunc func(ctx context.Context, text string) ([]float32, error)

// DB 包含多个集合，每个集合包含多个文档。
type DB struct {
	collections     map[string]*Collection
	collectionsLock sync.RWMutex

	persistDirectory string
	compress         bool
}

// NewDB 创建一个新的内存中的数据库。
func NewDB() *DB {
	return &DB{
		collections: make(map[string]*Collection),
	}
}

// NewPersistentDB 创建一个新的持久化的数据库。
// 如果路径为空，默认为 "./godoos/data/godoDB"。
// 如果 compress 为 true，则文件将使用 gzip 压缩。
func NewPersistentDB(path string, compress bool) (*DB, error) {
	homeDir, err := libs.GetAppDir()
	if err != nil {
		return nil, err
	}
	if path == "" {
		path = filepath.Join(homeDir, "data", "godoDB")
	} else {
		path = filepath.Clean(path)
	}

	ext := ".gob"
	if compress {
		ext += ".gz"
	}

	db := &DB{
		collections:      make(map[string]*Collection),
		persistDirectory: path,
		compress:         compress,
	}

	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := os.MkdirAll(path, 0o700)
			if err != nil {
				return nil, fmt.Errorf("无法创建持久化目录: %w", err)
			}
			return db, nil
		}
		return nil, fmt.Errorf("无法获取持久化目录信息: %w", err)
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("路径不是目录: %s", path)
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取持久化目录: %w", err)
	}
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		collectionPath := filepath.Join(path, dirEntry.Name())
		collectionDirEntries, err := os.ReadDir(collectionPath)
		if err != nil {
			return nil, fmt.Errorf("无法读取集合目录: %w", err)
		}
		c := &Collection{
			documents:        make(map[string]*Document),
			persistDirectory: collectionPath,
			compress:         compress,
		}
		for _, collectionDirEntry := range collectionDirEntries {
			if collectionDirEntry.IsDir() {
				continue
			}
			fPath := filepath.Join(collectionPath, collectionDirEntry.Name())
			if collectionDirEntry.Name() == metadataFileName+ext {
				pc := struct {
					Name     string
					Metadata map[string]string
				}{}
				err := readFromFile(fPath, &pc, "")
				if err != nil {
					return nil, fmt.Errorf("无法读取集合元数据: %w", err)
				}
				c.Name = pc.Name
				c.metadata = pc.Metadata
			} else if strings.HasSuffix(collectionDirEntry.Name(), ext) {
				d := &Document{}
				err := readFromFile(fPath, d, "")
				if err != nil {
					return nil, fmt.Errorf("无法读取文档: %w", err)
				}
				c.documents[d.ID] = d
			}
		}
		if c.Name == "" && len(c.documents) == 0 {
			continue
		}
		if c.Name == "" {
			return nil, fmt.Errorf("未找到集合元数据文件: %s", collectionPath)
		}
		db.collections[c.Name] = c
	}

	return db, nil
}

// ImportFromFile 从给定路径的文件导入数据库。
func (db *DB) ImportFromFile(filePath string, encryptionKey string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径为空")
	}
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须为 32 字节长")
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("文件不存在: %s", filePath)
		}
		return fmt.Errorf("无法获取文件信息: %w", err)
	} else if fi.IsDir() {
		return fmt.Errorf("路径是目录: %s", filePath)
	}

	type persistenceCollection struct {
		Name      string
		Metadata  map[string]string
		Documents map[string]*Document
	}
	persistenceDB := struct {
		Collections map[string]*persistenceCollection
	}{
		Collections: make(map[string]*persistenceCollection, len(db.collections)),
	}

	db.collectionsLock.Lock()
	defer db.collectionsLock.Unlock()

	err = readFromFile(filePath, &persistenceDB, encryptionKey)
	if err != nil {
		return fmt.Errorf("无法读取文件: %w", err)
	}

	for _, pc := range persistenceDB.Collections {
		c := &Collection{
			Name:      pc.Name,
			metadata:  pc.Metadata,
			documents: pc.Documents,
		}
		if db.persistDirectory != "" {
			c.persistDirectory = filepath.Join(db.persistDirectory, hash2hex(pc.Name))
			c.compress = db.compress
		}
		db.collections[c.Name] = c
	}

	return nil
}

// ImportFromReader 从 reader 导入数据库。
func (db *DB) ImportFromReader(reader io.ReadSeeker, encryptionKey string) error {
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须为 32 字节长")
	}

	type persistenceCollection struct {
		Name      string
		Metadata  map[string]string
		Documents map[string]*Document
	}
	persistenceDB := struct {
		Collections map[string]*persistenceCollection
	}{
		Collections: make(map[string]*persistenceCollection, len(db.collections)),
	}

	db.collectionsLock.Lock()
	defer db.collectionsLock.Unlock()

	err := readFromReader(reader, &persistenceDB, encryptionKey)
	if err != nil {
		return fmt.Errorf("无法读取流: %w", err)
	}

	for _, pc := range persistenceDB.Collections {
		c := &Collection{
			Name:      pc.Name,
			metadata:  pc.Metadata,
			documents: pc.Documents,
		}
		if db.persistDirectory != "" {
			c.persistDirectory = filepath.Join(db.persistDirectory, hash2hex(pc.Name))
			c.compress = db.compress
		}
		db.collections[c.Name] = c
	}

	return nil
}

// ExportToFile 将数据库导出到给定路径的文件。
func (db *DB) ExportToFile(filePath string, compress bool, encryptionKey string) error {
	if filePath == "" {
		filePath = "./gododb.gob"
		if compress {
			filePath += ".gz"
		}
		if encryptionKey != "" {
			filePath += ".enc"
		}
	}
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须为 32 字节长")
	}

	type persistenceCollection struct {
		Name      string
		Metadata  map[string]string
		Documents map[string]*Document
	}
	persistenceDB := struct {
		Collections map[string]*persistenceCollection
	}{
		Collections: make(map[string]*persistenceCollection, len(db.collections)),
	}

	db.collectionsLock.RLock()
	defer db.collectionsLock.RUnlock()

	for k, v := range db.collections {
		persistenceDB.Collections[k] = &persistenceCollection{
			Name:      v.Name,
			Metadata:  v.metadata,
			Documents: v.documents,
		}
	}

	err := persistToFile(filePath, persistenceDB, compress, encryptionKey)
	if err != nil {
		return fmt.Errorf("无法导出数据库: %w", err)
	}

	return nil
}

// ExportToWriter 将数据库导出到 writer。
func (db *DB) ExportToWriter(writer io.Writer, compress bool, encryptionKey string) error {
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须为 32 字节长")
	}

	type persistenceCollection struct {
		Name      string
		Metadata  map[string]string
		Documents map[string]*Document
	}
	persistenceDB := struct {
		Collections map[string]*persistenceCollection
	}{
		Collections: make(map[string]*persistenceCollection, len(db.collections)),
	}

	db.collectionsLock.RLock()
	defer db.collectionsLock.RUnlock()

	for k, v := range db.collections {
		persistenceDB.Collections[k] = &persistenceCollection{
			Name:      v.Name,
			Metadata:  v.metadata,
			Documents: v.documents,
		}
	}

	err := persistToWriter(writer, persistenceDB, compress, encryptionKey)
	if err != nil {
		return fmt.Errorf("无法导出数据库: %w", err)
	}

	return nil
}

// CreateCollection 创建具有给定名称和元数据的新集合。
func (db *DB) CreateCollection(name string, metadata map[string]string, embeddingFunc EmbeddingFunc) (*Collection, error) {
	if name == "" {
		return nil, errors.New("集合名称为空")
	}
	if embeddingFunc == nil {
		embeddingFunc = NewEmbeddingFuncDefault()
	}
	collection, err := newCollection(name, metadata, embeddingFunc, db.persistDirectory, db.compress)
	if err != nil {
		return nil, fmt.Errorf("无法创建集合: %w", err)
	}

	db.collectionsLock.Lock()
	defer db.collectionsLock.Unlock()
	db.collections[name] = collection
	return collection, nil
}

// ListCollections 返回数据库中的所有集合。
func (db *DB) ListCollections() map[string]*Collection {
	db.collectionsLock.RLock()
	defer db.collectionsLock.RUnlock()

	res := make(map[string]*Collection, len(db.collections))
	for k, v := range db.collections {
		res[k] = v
	}

	return res
}

// GetCollection 返回具有给定名称的集合。
func (db *DB) GetCollection(name string, embeddingFunc EmbeddingFunc) *Collection {
	db.collectionsLock.RLock()
	defer db.collectionsLock.RUnlock()

	c, ok := db.collections[name]
	if !ok {
		return nil
	}

	if c.embed == nil {
		if embeddingFunc == nil {
			c.embed = NewEmbeddingFuncDefault()
		} else {
			c.embed = embeddingFunc
		}
	}
	return c
}

// GetOrCreateCollection 返回数据库中已有的集合，或创建一个新的集合。
func (db *DB) GetOrCreateCollection(name string, metadata map[string]string, embeddingFunc EmbeddingFunc) (*Collection, error) {
	collection := db.GetCollection(name, embeddingFunc)
	if collection == nil {
		var err error
		collection, err = db.CreateCollection(name, metadata, embeddingFunc)
		if err != nil {
			return nil, fmt.Errorf("无法创建集合: %w", err)
		}
	}
	return collection, nil
}

// DeleteCollection 删除具有给定名称的集合。
func (db *DB) DeleteCollection(name string) error {
	db.collectionsLock.Lock()
	defer db.collectionsLock.Unlock()

	col, ok := db.collections[name]
	if !ok {
		return nil
	}

	if db.persistDirectory != "" {
		collectionPath := col.persistDirectory
		err := os.RemoveAll(collectionPath)
		if err != nil {
			return fmt.Errorf("无法删除集合目录: %w", err)
		}
	}

	delete(db.collections, name)
	return nil
}

// Reset 从数据库中移除所有集合。
func (db *DB) Reset() error {
	db.collectionsLock.Lock()
	defer db.collectionsLock.Unlock()

	if db.persistDirectory != "" {
		err := os.RemoveAll(db.persistDirectory)
		if err != nil {
			return fmt.Errorf("无法删除持久化目录: %w", err)
		}
		err = os.MkdirAll(db.persistDirectory, 0o700)
		if err != nil {
			return fmt.Errorf("无法重新创建持久化目录: %w", err)
		}
	}

	db.collections = make(map[string]*Collection)
	return nil
}
