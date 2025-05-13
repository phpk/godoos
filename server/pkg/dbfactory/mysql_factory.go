package dbfactory

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// MySQL工厂实现
type MySQLFactory struct{}

func NewMySQLFactory() *MySQLFactory {
	return &MySQLFactory{}
}

func (f *MySQLFactory) CreateConnection() DatabaseConnection {
	return MySQLDB
}

func (f *MySQLFactory) CreateRepository() BaseRepository {
	return &mysqlRepository{
		db: MySQLDB,
	}
}

type mysqlRepository struct {
	db *gorm.DB
}

func (c *mysqlRepository) AutoMigrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}
func (r *mysqlRepository) Create(table string, entity interface{}) error {
	return r.db.Table(table).Create(entity).Error
}
func (r *mysqlRepository) Count(table string, conditions map[string]interface{}) (int64, error) {
	var count int64
	err := r.db.Table(table).Where(conditions).Count(&count).Error
	return count, err
}
func (r *mysqlRepository) GetByID(table string, id interface{}, result interface{}) error {
	return r.db.Table(table).First(result, id).Error
}

func (r *mysqlRepository) Update(table string, entity interface{}) error {
	return r.db.Table(table).Save(entity).Error
}

func (r *mysqlRepository) Delete(table string, id interface{}) error {
	return r.db.Table(table).Delete(nil, id).Error
}

func (r *mysqlRepository) BatchCreate(table string, entities []interface{}) error {
	return r.db.Table(table).CreateInBatches(entities, len(entities)).Error
}

func (r *mysqlRepository) BatchDelete(table string, ids []interface{}) error {
	return r.db.Table(table).Delete(nil, ids).Error
}

func (r *mysqlRepository) JoinQuery(mainTable string, joinTables []JoinTable, conditions map[string]interface{}, result interface{}) error {
	tx := r.db.Table(mainTable)
	for _, join := range joinTables {
		tx = tx.Joins(join.JoinType + " JOIN " + join.Table + " ON " + join.OnCondition)
	}
	return tx.Where(conditions).Find(result).Error
}

func (r *mysqlRepository) GetOne(table string, conditions map[string]interface{}, result interface{}) error {
	return r.db.Table(table).Where(conditions).First(result).Error
}

func (r *mysqlRepository) GetList(table string, conditions map[string]interface{}, result interface{}) error {
	return r.db.Table(table).Where(conditions).Find(result).Error
}

func (r *mysqlRepository) GetPage(table string, conditions map[string]interface{}, pageParams PageParams, result *PageResult) error {
	var total int64
	if err := r.db.Table(table).Where(conditions).Count(&total).Error; err != nil {
		return err
	}

	query := r.db.Table(table).Where(conditions)
	if pageParams.OrderBy != "" {
		query = query.Order(pageParams.OrderBy)
	}
	if err := query.Offset((pageParams.Page - 1) * pageParams.PageSize).
		Limit(pageParams.PageSize).
		Find(result.Data).Error; err != nil {
		return err
	}

	result.Total = total
	return nil
}
func (r *mysqlRepository) WithTransaction(ctx context.Context, opts *TxOptions, fn func(repo BaseRepository) error) error {
	txOpts := &sql.TxOptions{}
	switch opts.Isolation {
	case LevelReadUncommitted:
		txOpts.Isolation = sql.LevelReadUncommitted
	case LevelReadCommitted:
		txOpts.Isolation = sql.LevelReadCommitted
	case LevelRepeatableRead:
		txOpts.Isolation = sql.LevelRepeatableRead
	case LevelSerializable:
		txOpts.Isolation = sql.LevelSerializable
	}
	txOpts.ReadOnly = opts.ReadOnly

	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&mysqlRepository{db: tx.WithContext(ctx)})
	}, txOpts)
}
