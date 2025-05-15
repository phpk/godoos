package db

import (
	"encoding/json"
	"fmt"
	"godocms/pkg/db/mongodm"
	"log"
	"os"

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
	MongoDB      *gorm.DB
	AppConfig    DatabaseConfig
)
var DB *gorm.DB
var DBS = make(map[string]*gorm.DB)

func InitDatabase() error {
	if err := LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
		return err
	}

	switch AppConfig.DBType {
	case "postgresql":
		if err := InitPostgreSQL(); err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
			return err
		}
		DB = PostgreSQLDB
	case "sqlite":
		if err := InitSQLite(); err != nil {
			log.Fatal("Failed to connect to SQLite:", err)
			return err
		}
		DB = SQLiteDB

	case "mongodb":
		if err := InitMongoDB(); err != nil {
			log.Fatal("Failed to connect to MongoDB:", err)
			return err
		}
		DB = MongoDB
	default:
		if err := InitMySQL(); err != nil {
			log.Fatal("Failed to connect to MySQL:", err)
			return err
		}
		DB = MySQLDB
	}

	// 支持多数据库共存
	if AppConfig.MySQL.Enable && AppConfig.DBType != "mysql" {
		InitMySQL()
		DBS["mysql"] = MySQLDB
	}
	if AppConfig.MongoDB.Enable && AppConfig.DBType != "mongodb" {
		InitMongoDB()
		DBS["mongodb"] = MongoDB
	}
	if AppConfig.SQLite.Enable && AppConfig.DBType != "sqlite" {
		InitSQLite()
		DBS["sqlite"] = SQLiteDB
	}
	if AppConfig.PostgreSQL.Enable && AppConfig.DBType != "postgresql" {
		InitPostgreSQL()
		DBS["postgresql"] = PostgreSQLDB
	}
	return nil
}
func LoadConfig() error {
	data, err := os.ReadFile("config/database.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &AppConfig)
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
	dialector := &mongodm.MongoDBDialector{
		URI:      AppConfig.MongoDB.URI, // MongoDB 连接 URI
		Database: AppConfig.MongoDB.Database,
	}
	var err error
	MongoDB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("MongoDB connected successfully")
	return nil
}
