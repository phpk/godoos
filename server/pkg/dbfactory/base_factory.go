package dbfactory

import "context"

// 事务隔离级别
type IsolationLevel int

const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelRepeatableRead
	LevelSerializable
)

// 事务选项
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

// 分页参数
type PageParams struct {
	Page     int
	PageSize int
	OrderBy  string
}

// 分页结果
type PageResult struct {
	Total int64
	Data  interface{}
}

// 基础CRUD接口
type BaseRepository interface {
	// 单条操作
	Create(table string, entity interface{}) error
	Count(table string, conditions map[string]interface{}) (int64, error)
	GetByID(table string, id interface{}, result interface{}) error
	Update(table string, entity interface{}) error
	Delete(table string, id interface{}) error
	AutoMigrate(models ...interface{}) error
	// 查询操作
	GetOne(table string, conditions map[string]interface{}, result interface{}) error
	GetList(table string, conditions map[string]interface{}, result interface{}) error
	GetPage(table string, conditions map[string]interface{}, pageParams PageParams, result *PageResult) error

	// 批量操作
	BatchCreate(table string, entities []interface{}) error
	BatchDelete(table string, ids []interface{}) error

	// 联表查询
	JoinQuery(mainTable string, joinTables []JoinTable, conditions map[string]interface{}, result interface{}) error

	// 事务支持
	WithTransaction(ctx context.Context, opts *TxOptions, fn func(repo BaseRepository) error) error
}

// 联表定义
type JoinTable struct {
	Table       string
	JoinType    string // "INNER", "LEFT", "RIGHT"
	OnCondition string
}

// 数据库连接接口
type DatabaseConnection interface {
	//AutoMigrate(models ...interface{}) error
}

// 工厂接口
type RepositoryFactory interface {
	CreateConnection() DatabaseConnection
	CreateRepository() BaseRepository
}
