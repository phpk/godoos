package vector

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"godo/model"
	"godo/office"
	"net/http"
	"path/filepath"
)

func HandlerCreateKnowledge(w http.ResponseWriter, r *http.Request) {
	var req model.VecList
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "the chat request error:"+err.Error())
		return
	}
	if req.FilePath == "" {
		libs.ErrorMsg(w, "file path is empty")
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.ErrorMsg(w, "get vector db path error:"+err.Error())
		return
	}
	req.FilePath = filepath.Join(basePath, req.FilePath)

	id, err := CreateVector(req)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, id, "create vector success")
}

// CreateVector 创建一个新的 VectorList 记录
func CreateVector(data model.VecList) (uint, error) {
	if data.FilePath == "" {
		return 0, fmt.Errorf("file path is empty")
	}
	if data.Engine == "" {
		return 0, fmt.Errorf("engine is empty")
	}

	if !libs.PathExists(data.FilePath) {
		return 0, fmt.Errorf("file path does not exist")
	}
	if data.EmbeddingModel == "" {
		return 0, fmt.Errorf("embedding model is empty")
	}
	if data.EmbedSize == 0 {
		data.EmbedSize = 768
	}
	// Create the new VectorList
	result := model.Db.Create(&data)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create vector list: %w", result.Error)
	}

	// Start background tasks
	go office.SetDocument(data.FilePath, uint(data.ID))

	return uint(data.ID), nil
}

// DeleteVector 删除指定id的 VectorList 记录
func DeleteVector(id int) error {

	return model.Db.Delete(&model.VecList{}, id).Error
}

// RenameVectorDb 更改指定名称的 VectorList 的数据库名称
func RenameVectorDb(oldName string, newName string) error {
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("failed to find old vector list: %w", err)
	}

	// 获取旧的 VectorList 记录
	var oldList model.VecList
	oldPath := filepath.Join(basePath, oldName)
	if err := model.Db.Where("file_path = ?", oldPath).First(&oldList).Error; err != nil {
		return fmt.Errorf("failed to find old vector list: %w", err)
	}

	// 更新 VectorList 记录中的 FilePath
	newPath := filepath.Join(basePath, newName)
	if err := model.Db.Model(&model.VecList{}).Where("id = ?", oldList.ID).Update("file_path", newPath).Error; err != nil {
		return fmt.Errorf("failed to update vector list: %w", err)
	}

	return nil
}

func GetVectorList() ([]model.VecList, error) {
	var vectorList []model.VecList
	if err := model.Db.Find(&vectorList).Error; err != nil {
		return nil, fmt.Errorf("failed to get vector list: %w", err)
	}
	return vectorList, nil
}

func GetVector(id uint) (model.VecList, error) {
	var vectorList model.VecList
	if err := model.Db.First(&vectorList, id).Error; err != nil {
		return vectorList, fmt.Errorf("failed to get vector: %w", err)
	}
	return vectorList, nil
}

// func SimilaritySearch(query string, numDocuments int, collection string, where map[string]string) ([]vs.Document, error) {
// 	ef := v.embeddingFunc
// 	if embeddingFunc != nil {
// 		ef = embeddingFunc
// 	}

// 	q, err := ef(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to compute embedding: %w", err)
// 	}

// 	qv, err := sqlitevec.SerializeFloat32(q)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to serialize query embedding: %w", err)
// 	}

// 	var docs []vs.Document
// 	err = v.db.Transaction(func(tx *gorm.DB) error {
// 		// Query matching document IDs and distances
// 		rows, err := tx.Raw(fmt.Sprintf(`
//             SELECT document_id, distance
//             FROM [%s_vec]
//             WHERE embedding MATCH ?
//             ORDER BY distance
//             LIMIT ?
//         `, collection), qv, numDocuments).Rows()
// 		if err != nil {
// 			return fmt.Errorf("failed to query vector table: %w", err)
// 		}
// 		defer rows.Close()

// 		for rows.Next() {
// 			var docID string
// 			var distance float32
// 			if err := rows.Scan(&docID, &distance); err != nil {
// 				return fmt.Errorf("failed to scan row: %w", err)
// 			}
// 			docs = append(docs, vs.Document{
// 				ID:              docID,
// 				SimilarityScore: 1 - distance, // Higher score means closer match
// 			})
// 		}

// 		// Fetch content and metadata for each document
// 		for i, doc := range docs {
// 			var content string
// 			var metadataJSON []byte
// 			err := tx.Raw(fmt.Sprintf(`
//                 SELECT content, metadata
//                 FROM [%s]
//                 WHERE id = ?
//             `, v.embeddingsTableName), doc.ID).Row().Scan(&content, &metadataJSON)
// 			if err != nil {
// 				return fmt.Errorf("failed to query embeddings table for document %s: %w", doc.ID, err)
// 			}

// 			var metadata map[string]interface{}
// 			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
// 				return fmt.Errorf("failed to parse metadata for document %s: %w", doc.ID, err)
// 			}

// 			docs[i].Content = content
// 			docs[i].Metadata = metadata
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return docs, nil
// }

// func AddDocuments(docs []VectorDoc, collection string) ([]string, error) {
// 	ids := make([]string, len(docs))

// 	err := VecDb.Transaction(func(tx *gorm.DB) error {
// 		if len(docs) > 0 {
// 			valuePlaceholders := make([]string, len(docs))
// 			args := make([]interface{}, 0, len(docs)*2) // 2 args per doc: document_id and embedding

// 			for i, doc := range docs {
// 				emb, err := v.embeddingFunc(ctx, doc.Content)
// 				if err != nil {
// 					return fmt.Errorf("failed to compute embedding for document %s: %w", doc.ID, err)
// 				}

// 				serializedEmb, err := sqlitevec.SerializeFloat32(emb)
// 				if err != nil {
// 					return fmt.Errorf("failed to serialize embedding for document %s: %w", doc.ID, err)
// 				}

// 				valuePlaceholders[i] = "(?, ?)"
// 				args = append(args, doc.ID, serializedEmb)

// 				ids[i] = doc.ID
// 			}

// 			// Raw query for *_vec as gorm doesn't support virtual tables
// 			query := fmt.Sprintf(`
// 				INSERT INTO [%s_vec] (document_id, embedding)
// 				VALUES %s
// 			`, collection, strings.Join(valuePlaceholders, ", "))

// 			if err := tx.Exec(query, args...).Error; err != nil {
// 				return fmt.Errorf("failed to batch insert into vector table: %w", err)
// 			}
// 		}

// 		embs := make([]map[string]interface{}, len(docs))
// 		for i, doc := range docs {
// 			metadataJson, err := json.Marshal(doc.Metadata)
// 			if err != nil {
// 				return fmt.Errorf("failed to marshal metadata for document %s: %w", doc.ID, err)
// 			}
// 			embs[i] = map[string]interface{}{
// 				"id":            doc.ID,
// 				"collection_id": collection,
// 				"content":       doc.Content,
// 				"metadata":      metadataJson,
// 			}
// 		}

// 		// Use GORM's Create for the embeddings table
// 		if err := tx.Table(v.embeddingsTableName).Create(embs).Error; err != nil {
// 			return fmt.Errorf("failed to batch insert into embeddings table: %w", err)
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

//		return ids, nil
//	}
//func init() {

// dbPath := libs.GetVectorDb()
// sqlite_vec.Auto()

// db, err := sqlx.Connect("sqlite3", dbPath)
// if err != nil {
// 	fmt.Println("Failed to open SQLite database:", err)
// 	return
// }
// defer db.Close()
// VecDb = db
// dsn := "file:" + dbPath
// db, err := sqlite3.Open(dsn)
// //db, err := sqlite3.Open(":memory:")
// if err != nil {
// 	fmt.Println("Failed to open SQLite database:", err)
// 	return
// }
// stmt, _, err := db.Prepare(`SELECT vec_version()`)
// if err != nil {
// 	log.Fatal(err)
// }

// stmt.Step()
// log.Printf("vec_version=%s\n", stmt.ColumnText(0))
// stmt.Close()
// _, err = db.Exec("CREATE TABLE IF NOT EXISTS vec_list (id INTEGER PRIMARY KEY AUTOINCREMENT,file_path TEXT NOT NULL,engine TEXT NOT NULL,embedding_model TEXT NOT NULL)")
// if err != nil {
// 	log.Fatal(err)
// }
// _, err = db.Exec("CREATE TABLE IF NOT EXISTS vec_doc (id INTEGER PRIMARY KEY AUTOINCREMENT,list_id INTEGER DEFAULT 0,file_path TEXT,content TEXT)")
// if err != nil {
// 	log.Fatal(err)
// }
// _, err = db.Exec("CREATE VIRTUAL TABLE vec_items USING vec0(embedding float[768])")
// if err != nil {
// 	log.Fatal(err)
// }
// VecDb = db

//InitMonitor()
//}

// func HandlerCreateKnowledge(w http.ResponseWriter, r *http.Request) {
// 	var req VectorList
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		libs.ErrorMsg(w, "the chat request error:"+err.Error())
// 		return
// 	}
// 	if req.FilePath == "" {
// 		libs.ErrorMsg(w, "file path is empty")
// 		return
// 	}
// 	basePath, err := libs.GetOsDir()
// 	if err != nil {
// 		libs.ErrorMsg(w, "get vector db path error:"+err.Error())
// 		return
// 	}
// 	req.FilePath = filepath.Join(basePath, req.FilePath)

// 	// id, err := CreateVector(req)
// 	// if err != nil {
// 	// 	libs.ErrorMsg(w, err.Error())
// 	// 	return
// 	// }
// 	// libs.SuccessMsg(w, id, "create vector success")
// }

// CreateVector 创建一个新的 VectorList 记录
// func CreateVector(data VectorList) (uint, error) {
// 	if data.FilePath == "" {
// 		return 0, fmt.Errorf("file path is empty")
// 	}
// 	if data.Engine == "" {
// 		return 0, fmt.Errorf("engine is empty")
// 	}

// 	if !libs.PathExists(data.FilePath) {
// 		return 0, fmt.Errorf("file path does not exist")
// 	}
// 	if data.EmbeddingModel == "" {
// 		return 0, fmt.Errorf("embedding model is empty")
// 	}

// 	// Create the new VectorList
// 	if err := tx.Table("vec_list").Create(data).Error; err != nil {
// 		return fmt.Errorf("failed to batch insert into embeddings table: %w", err)
// 	}
// 	// Get the ID of the newly created VectorList
// 	vectorID, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Start background tasks
// 	go office.SetDocument(data.FilePath, uint(vectorID))

// 	return uint(vectorID), nil
// }

// // DeleteVector 删除指定id的 VectorList 记录
// func DeleteVector(id int) error {
// 	tx, err := VecDb.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Delete from vec_doc first
// 	_, err = tx.Exec("DELETE FROM vec_doc WHERE list_id = ?)", id)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete from vec_list
// 	result, err := tx.Exec("DELETE FROM vec_list WHERE id = ?", id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("vector list not found")
// 	}

// 	return tx.Commit()
// }

// // RenameVectorDb 更改指定名称的 VectorList 的数据库名称
// func RenameVectorDb(oldName string, newName string) error {
// 	basePath, err := libs.GetOsDir()
// 	if err != nil {
// 		return fmt.Errorf("failed to find old vector list: %w", err)
// 	}

// 	// 2. 获取旧的 VectorList 记录
// 	var oldList VectorList
// 	oldPath := filepath.Join(basePath, oldName)
// 	err = VecDb.QueryRow("SELECT id FROM vec_list WHERE file_path = ?", oldPath).Scan(&oldList.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to find old vector list: %w", err)
// 	}
// 	MapFilePathMonitors[oldPath] = 0

// 	// 5. 更新 VectorList 记录中的 DbPath 和 Name
// 	newPath := filepath.Join(basePath, newName)
// 	_, err = VecDb.Exec("UPDATE vec_list SET file_path = ? WHERE id = ?", newPath, oldList.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update vector list: %w", err)
// 	}
// 	MapFilePathMonitors[newPath] = oldList.ID

// 	return nil
// }
// func InsertVectorDoc(data []VectorDoc, embedlist [][]float32) error {
// 	rowIds := map[int][]float32{}
// 	for i, v := range data {
// 		err := VecDb.Exec("INSERT INTO vec_doc (list_id, file_path, content) VALUES (?, ?, ?)", v.ListID, v.FilePath, v.Content)
// 		if err != nil {
// 			return err
// 		}
// 		rowID, err := result.LastInsertRowID()
// 		if err != nil {
// 			return err
// 		}
// 		rowid := int(rowID)
// 		rowIds[rowid] = embedlist[i]
// 	}
// 	stmt, err := VecDb.Prepare("INSERT INTO vec_items(rowid, embedding) VALUES (?, ?)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()

// 	for id, values := range rowIds {
// 		v, err := sqlite_vec.SerializeFloat32(values)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = stmt.BindInt64(1, int64(id))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = stmt.BindBlob(2, v)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = stmt.Exec()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		stmt.Reset()
// 	}

// 	return nil
// }
// func InitMonitor() {
// 	list := GetVectorList()
// 	for _, v := range list {
// 		MapFilePathMonitors[v.FilePath] = v.ID
// 	}
// 	FolderMonitor()
// }

// func GetVectorList() []VectorList {
// 	var vectorList []VectorList
// 	// stmt, _, err := VecDb.Prepare("SELECT id, file_path, engine, embedding_model FROM vec_list")
// 	// if err != nil {
// 	// 	fmt.Println("Failed to get vector list:", err)
// 	// 	return vectorList
// 	// }
// 	// stmt.Step()
// 	// log.Printf("vec_version=%s\n", stmt.ColumnText(0))
// 	// stmt.Close()
// 	// defer rows.Close()

// 	// for rows.Next() {
// 	// 	var v VectorList
// 	// 	err := rows.Scan(&v.ID, &v.FilePath, &v.Engine, &v.EmbeddingModel)
// 	// 	if err != nil {
// 	// 		fmt.Println("Failed to scan vector list row:", err)
// 	// 		continue
// 	// 	}
// 	// 	vectorList = append(vectorList, v)
// 	// }

// 	return vectorList
// }
// func GetVector(id uint) VectorList {
// 	var vectorList VectorList
// 	// sql := "SELECT id, file_path, engine, embedding_model FROM vec_list WHERE id = " + fmt.Sprintf("%d", id)
// 	// stmt, _, err := VecDb.Prepare(sql)
// 	// if err != nil {
// 	// 	fmt.Println("Failed to get vector list:", err)
// 	// 	return vectorList
// 	// }
// 	// stmt.Step()
// 	// log.Printf("vec_version=%s\n", stmt.ColumnText(0))
// 	// stmt.Close()
// 	// err := VecDb.QueryRow("SELECT id, file_path, engine, embedding_model FROM vec_list WHERE id = ?", id).Scan(&vectorList.ID, &vectorList.FilePath, &vectorList.Engine, &vectorList.EmbeddingModel)
// 	// if err != nil {
// 	// 	fmt.Println("Failed to get vector:", err)
// 	// }
// 	return vectorList
// }
