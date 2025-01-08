package godovec

import (
	"fmt"

	_ "embed"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"

	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var VecDb *gorm.DB

type VectorList struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	FilePath       string `json:"file_path" gorm:"not null"`
	Engine         string `json:"engine" gorm:"not null"`
	EmbeddingModel string `json:"model" gorm:"not null"`
}

type VectorDoc struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Content  string `json:"content"`
	FilePath string `json:"file_path" gorm:"not null"`
	ListID   int    `json:"list_id"`
}

func main() {
	InitVector()
}
func InitVector() error {

	db, err := gorm.Open(gormlite.Open("./data.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open vector db: %w", err)
	}

	// Enable PRAGMAs
	// - busy_timeout (ms) to prevent db lockups as we're accessing the DB from multiple separate processes in otto8
	tx := db.Exec(`
PRAGMA busy_timeout = 5000;
`)
	if tx.Error != nil {
		return fmt.Errorf("failed to execute pragma busy_timeout: %w", tx.Error)
	}
	err = db.AutoMigrate(&VectorList{}, &VectorDoc{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}
	VecDb = db
	return nil
}
