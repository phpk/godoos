package dbfactory

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type PostgreSQLFactory struct{}

func NewPostgreSQLFactory() *PostgreSQLFactory {
	return &PostgreSQLFactory{}
}

func (f *PostgreSQLFactory) CreateConnection() DatabaseConnection {
	return &postgresqlConnection{
		db: PostgreSQLDB,
	}
}

func (f *PostgreSQLFactory) CreateRepository() BaseRepository {
	return &postgresqlRepository{
		db: PostgreSQLDB,
	}
}

type postgresqlConnection struct {
	db *gorm.DB
}

func (c *postgresqlConnection) AutoMigrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}

type postgresqlRepository struct {
	db *gorm.DB
}

func (r *postgresqlRepository) Create(table string, entity interface{}) error {
	return r.db.Table(table).Create(entity).Error
}
func (r *postgresqlRepository) Count(table string, conditions map[string]interface{}) (int64, error) {
	var count int64
	err := r.db.Table(table).Where(conditions).Count(&count).Error
	return count, err
}
func (r *postgresqlRepository) GetByID(table string, id interface{}, result interface{}) error {
	return r.db.Table(table).First(result, id).Error
}

func (r *postgresqlRepository) Update(table string, entity interface{}) error {
	return r.db.Table(table).Save(entity).Error
}

func (r *postgresqlRepository) Delete(table string, id interface{}) error {
	return r.db.Table(table).Delete(nil, id).Error
}

func (r *postgresqlRepository) GetOne(table string, conditions map[string]interface{}, result interface{}) error {
	return r.db.Table(table).Where(conditions).First(result).Error
}

func (r *postgresqlRepository) GetList(table string, conditions map[string]interface{}, result interface{}) error {
	return r.db.Table(table).Where(conditions).Find(result).Error
}

func (r *postgresqlRepository) GetPage(table string, conditions map[string]interface{}, pageParams PageParams, result *PageResult) error {
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

func (r *postgresqlRepository) BatchCreate(table string, entities []interface{}) error {
	return r.db.Table(table).CreateInBatches(entities, len(entities)).Error
}

func (r *postgresqlRepository) BatchDelete(table string, ids []interface{}) error {
	return r.db.Table(table).Delete(nil, ids).Error
}

func (r *postgresqlRepository) JoinQuery(mainTable string, joinTables []JoinTable, conditions map[string]interface{}, result interface{}) error {
	tx := r.db.Table(mainTable)
	for _, join := range joinTables {
		tx = tx.Joins(join.JoinType + " JOIN " + join.Table + " ON " + join.OnCondition)
	}
	return tx.Where(conditions).Find(result).Error
}

// PostgreSQL特有功能：JSONB操作
func (r *postgresqlRepository) JSONBContains(table, column, path string, value interface{}, result interface{}) error {
	return r.db.Table(table).Where(column+"::jsonb @> ?", `{"`+path+`":`+value.(string)+`}`).Find(result).Error
}

func (r *postgresqlRepository) AutoMigrate(models ...interface{}) error {
	return r.db.AutoMigrate(models...)
}

// 事务支持
func (r *postgresqlRepository) WithTransaction(ctx context.Context, opts *TxOptions, fn func(repo BaseRepository) error) error {
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
		return fn(&postgresqlRepository{db: tx.WithContext(ctx)})
	}, txOpts)
}
