package vector

import (
	"encoding/json"
	"godo/libs"
	"net/http"
	"path/filepath"

	_ "embed"
	// sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	// "github.com/ncruces/go-sqlite3"
)

//var VecDb *sqlx.DB

type VectorList struct {
	ID             int    `json:"id"`
	FilePath       string `json:"file_path"`
	Engine         string `json:"engine"`
	EmbeddingModel string `json:"model"`
}
type VectorDoc struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
	ListID   int    `json:"list_id"`
}
type VectorItem struct {
	DocID     int       `json:"rowid"`
	Embedding []float32 `json:"embedding"`
}

func init() {

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
}

func HandlerCreateKnowledge(w http.ResponseWriter, r *http.Request) {
	var req VectorList
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

	// id, err := CreateVector(req)
	// if err != nil {
	// 	libs.ErrorMsg(w, err.Error())
	// 	return
	// }
	// libs.SuccessMsg(w, id, "create vector success")
}

// // CreateVector 创建一个新的 VectorList 记录
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

// 	// Check if a VectorList with the same path already exists

// 	stmt, _, err := VecDb.Prepare(`SELECT id FROM vec_list WHERE file_path =`+ data.FilePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	for stmt.Step() {
// 		fmt.Println(stmt.ColumnInt(0), stmt.ColumnText(1))
// 	}
// 	if err := stmt.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	err = stmt.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Create the new VectorList
// 	err = VecDb.Exec("INSERT INTO vec_list (file_path, engine, embedding_model) VALUES (?, ?, ?)", data.FilePath, data.Engine, data.EmbeddingModel)
// 	if err != nil {
// 		return 0, err
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

func GetVectorList() []VectorList {
	var vectorList []VectorList
	// stmt, _, err := VecDb.Prepare("SELECT id, file_path, engine, embedding_model FROM vec_list")
	// if err != nil {
	// 	fmt.Println("Failed to get vector list:", err)
	// 	return vectorList
	// }
	// stmt.Step()
	// log.Printf("vec_version=%s\n", stmt.ColumnText(0))
	// stmt.Close()
	// defer rows.Close()

	// for rows.Next() {
	// 	var v VectorList
	// 	err := rows.Scan(&v.ID, &v.FilePath, &v.Engine, &v.EmbeddingModel)
	// 	if err != nil {
	// 		fmt.Println("Failed to scan vector list row:", err)
	// 		continue
	// 	}
	// 	vectorList = append(vectorList, v)
	// }

	return vectorList
}
func GetVector(id uint) VectorList {
	var vectorList VectorList
	// sql := "SELECT id, file_path, engine, embedding_model FROM vec_list WHERE id = " + fmt.Sprintf("%d", id)
	// stmt, _, err := VecDb.Prepare(sql)
	// if err != nil {
	// 	fmt.Println("Failed to get vector list:", err)
	// 	return vectorList
	// }
	// stmt.Step()
	// log.Printf("vec_version=%s\n", stmt.ColumnText(0))
	// stmt.Close()
	// err := VecDb.QueryRow("SELECT id, file_path, engine, embedding_model FROM vec_list WHERE id = ?", id).Scan(&vectorList.ID, &vectorList.FilePath, &vectorList.Engine, &vectorList.EmbeddingModel)
	// if err != nil {
	// 	fmt.Println("Failed to get vector:", err)
	// }
	return vectorList
}
