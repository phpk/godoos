package dbfactory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/qiniu/qmgo"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MySQL struct {
	Enable   bool   `json:"enable"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type MongosDB struct {
	Enable   bool   `json:"enable"`
	URI      string `json:"uri"`
	Database string `json:"database"`
}

type SQLite struct {
	Enable bool   `json:"enable"`
	Path   string `json:"path"` // 数据库文件路径
}

type PostgreSQL struct {
	Enable   bool   `json:"enable"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"` // disable/require
}
type DatabaseConfig struct {
	DBType     string     `json:"dbType"` // mysql/mongodb/sqlite/postgresql
	MySQL      MySQL      `json:"mysql"`
	MongoDB    MongosDB   `json:"mongodb"`
	SQLite     SQLite     `json:"sqlite"`
	PostgreSQL PostgreSQL `json:"postgresql"`
}

var (
	MySQLDB      *gorm.DB
	PostgreSQLDB *gorm.DB
	SQLiteDB     *gorm.DB
	MongoDB      *qmgo.Database
	AppConfig    DatabaseConfig
)

var Dbs = make(map[string]BaseRepository)

var Db BaseRepository

func LoadConfig() error {
	data, err := os.ReadFile("config/database.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &AppConfig)
}
func InitDatabase() error {
	// 加载配置并初始化数据库
	if err := LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
		return err
	}

	switch AppConfig.DBType {
	case "mongodb":
		if err := InitMongoDB(); err != nil {
			log.Fatal("Failed to connect to MongoDB:", err)
			return err
		}
		Db = NewMongoDBFactory().CreateRepository()
		Dbs["mongodb"] = Db
	case "sqlite":
		if err := InitSQLite(); err != nil {
			log.Fatal("Failed to connect to SQLite:", err)
			return err
		}
		Db = NewSQLiteFactory().CreateRepository()
		Dbs["sqlite"] = Db
	case "postgresql":
		if err := InitPostgreSQL(); err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
			return err
		}
		Db = NewPostgreSQLFactory().CreateRepository()
		Dbs["postgresql"] = Db
	default:
		if err := InitMySQL(); err != nil {
			log.Fatal("Failed to connect to MySQL:", err)
			return err
		}
		Db = NewMySQLFactory().CreateRepository()
		Dbs["mysql"] = Db
	}
	if AppConfig.MySQL.Enable && AppConfig.DBType != "mysql" {
		if err := InitMySQL(); err != nil {
			log.Println("Failed to connect to MySQL:", err)
		} else {
			Dbs["mysql"] = NewMySQLFactory().CreateRepository()
		}
	}
	if AppConfig.PostgreSQL.Enable && AppConfig.DBType != "postgresql" {
		if err := InitPostgreSQL(); err != nil {
			log.Println("Failed to connect to PostgreSQL:", err)
		} else {
			Dbs["postgresql"] = NewPostgreSQLFactory().CreateRepository()
		}
	}
	if AppConfig.SQLite.Enable && AppConfig.DBType != "sqlite" {
		if err := InitSQLite(); err != nil {
			log.Println("Failed to connect to SQLite:", err)
		} else {
			Dbs["sqlite"] = NewSQLiteFactory().CreateRepository()
		}
	}
	if AppConfig.MongoDB.Enable && AppConfig.DBType != "mongodb" {
		if err := InitMongoDB(); err != nil {
			log.Println("Failed to connect to MongoDB:", err)
		} else {
			Dbs["mongodb"] = NewMongoDBFactory().CreateRepository()
		}
	}
	return nil
}
func InitMySQL() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MySQL.User,
		AppConfig.MySQL.Password,
		AppConfig.MySQL.Host,
		AppConfig.MySQL.Port,
		AppConfig.MySQL.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	MySQLDB = db
	log.Println("MySQL connected successfully")
	return nil
}

func InitPostgreSQL() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		AppConfig.PostgreSQL.Host,
		AppConfig.PostgreSQL.Port,
		AppConfig.PostgreSQL.User,
		AppConfig.PostgreSQL.Password,
		AppConfig.PostgreSQL.DBName,
		AppConfig.PostgreSQL.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	PostgreSQLDB = db
	log.Println("PostgreSQL connected successfully")
	return nil
}

func InitSQLite() error {
	db, err := gorm.Open(sqlite.Open(AppConfig.SQLite.Path), &gorm.Config{})
	if err != nil {
		return err
	}

	SQLiteDB = db
	log.Println("SQLite connected successfully")
	return nil
}

func InitMongoDB() error {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: AppConfig.MongoDB.URI})
	if err != nil {
		return err
	}

	db := client.Database(AppConfig.MongoDB.Database)
	MongoDB = db
	log.Println("MongoDB connected successfully")
	return nil
}
