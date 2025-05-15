package mongodm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type MongoMigrator struct {
	db        *gorm.DB
	dialector MongoDBDialector
}

// HasTable 检查集合是否存在
func (m *MongoMigrator) HasTable(value interface{}) bool {
	name := m.db.NamingStrategy.TableName(reflect.TypeOf(value).Name())
	coll := m.dialector.Client.Database(m.dialector.Database).Collection(name)
	ctx := context.Background()
	colls, err := coll.Database().ListCollectionNames(ctx, bson.M{"name": name})
	return err == nil && len(colls) > 0
}

// CreateTable 创建集合（MongoDB 中相当于创建集合）
func (m *MongoMigrator) CreateTable(values ...interface{}) error {
	for _, value := range values {
		modelType := reflect.TypeOf(value)
		name := m.db.NamingStrategy.TableName(modelType.Name())
		coll := m.dialector.Client.Database(m.dialector.Database).Collection(name)

		ctx := context.Background()
		if err := coll.Database().CreateCollection(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

// DropTable 删除集合
func (m *MongoMigrator) DropTable(values ...interface{}) error {
	for _, value := range values {
		modelType := reflect.TypeOf(value)
		name := m.db.NamingStrategy.TableName(modelType.Name())
		coll := m.dialector.Client.Database(m.dialector.Database).Collection(name)

		ctx := context.Background()
		err := coll.Drop(ctx)
		if err != nil && !strings.Contains(err.Error(), "ns not found") {
			return err
		}
	}
	return nil
}

// GetTables 返回当前数据库中的所有表名（MongoDB 中对应集合名称）
func (m *MongoMigrator) GetTables() ([]string, error) {
	//var collectionNames []string
	// 使用 dialector 中的 Client 和 Database 名称获取集合列表
	collections, err := m.dialector.Client.Database(m.dialector.Database).ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	return collections, nil
}

// 实现 GetTypeAliases 方法以满足 gorm.Migrator 接口
func (m *MongoMigrator) GetTypeAliases(s string) []string {
	return []string{} // MongoDB 不需要类型别名，返回空切片即可
}

// AddIndex 添加索引到指定字段
func (m *MongoMigrator) AddIndex(value interface{}, indexName string, columnNames ...string) error {
	name := m.db.NamingStrategy.TableName(reflect.TypeOf(value).Name())
	coll := m.dialector.Client.Database(m.dialector.Database).Collection(name)

	keys := bson.D{}
	for _, col := range columnNames {
		keys = append(keys, bson.E{Key: col, Value: 1}) // 默认升序
	}

	indexOptions := options.Index()

	// 设置唯一索引（通过 tag 判断）
	if isUnique, ok := m.getIndexOptionFromTag(value, "unique", columnNames...); ok && isUnique {
		indexOptions.SetUnique(true)
	}

	// 可选：扩展其他索引类型（如稀疏、TTL 等）
	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions.SetName(indexName),
	}

	ctx := context.Background()
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	return err
}

// getIndexOptionFromTag 检查字段是否有特定索引选项（如 unique）
func (m *MongoMigrator) getIndexOptionFromTag(value interface{}, option string, columnNames ...string) (bool, bool) {
	v := reflect.ValueOf(value).Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("gorm")
		if tag == "" {
			continue
		}

		gormTag := parseGormTag(tag)
		for _, col := range columnNames {
			if gormTag["column"] == col || field.Name == col {
				switch option {
				case "unique":
					return true, gormTag["unique"] == "true"
				}
			}
		}
	}
	return false, false
}

// RenameTable 实现接口方法。MongoDB 不支持直接重命名集合，可以返回错误或忽略。
func (m *MongoMigrator) RenameTable(oldName interface{}, newName interface{}) error {
	// 可选：模拟重命名操作
	return nil // 或者返回 errors.New("rename table not supported in MongoDB")
}
func (m *MongoMigrator) TableType(name interface{}) (gorm.TableType, error) {
	return nil, nil
}

// parseGormTag 将 gorm tag 解析为 map
func parseGormTag(tag string) map[string]string {
	parts := strings.Split(tag, ";")
	result := make(map[string]string)
	for _, part := range parts {
		kv := strings.Split(part, ":")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else {
			result[part] = "true"
		}
	}
	return result
}

// AddColumn 不支持
func (m *MongoMigrator) AddColumn(value interface{}, field string) error {
	return nil // MongoDB 无需显式添加列
}

// HasConstraint 检查是否存在指定名称的约束（MongoDB 不支持，返回 false）
func (m *MongoMigrator) HasConstraint(s interface{}, name string) bool {
	return false
}

// 实现 MigrateColumnUnique 方法（MongoDB 不支持迁移，可以留空或返回 nil）
// interface{}, *schema.Field, gorm.ColumnType
func (m *MongoMigrator) MigrateColumnUnique(tableName interface{}, field *schema.Field, unique gorm.ColumnType) error {
	// MongoDB 不支持结构化迁移，直接返回 nil
	return nil
}

// interface{}, *schema.Field, gorm.ColumnType
func (m *MongoMigrator) MigrateColumn(schema interface{}, field *schema.Field, autoCreate gorm.ColumnType) error {
	return nil
}

// 实现 gorm.Migrator 接口的最简必要方法
func (m *MongoMigrator) HasColumn(table interface{}, column string) bool {
	// MongoDB 不依赖 schema，所以直接返回 true 或者不支持该功能
	return true // 或者根据实际集合结构判断是否存在字段
}

// DropColumn 不支持
func (m *MongoMigrator) DropColumn(value interface{}, field string) error {
	return nil
}

// ModifyColumn 不支持
func (m *MongoMigrator) ModifyColumn(value interface{}, field string) error {
	return nil
}

// RenameColumn 不支持
func (m *MongoMigrator) RenameColumn(value interface{}, oldName, newName string) error {
	return nil
}

// AlterColumn 修改字段（MongoDB 不支持，可返回 nil 或错误）
func (m *MongoMigrator) AlterColumn(value interface{}, field string) error {
	// MongoDB 不支持修改字段类型等操作，可以选择返回 nil 或提示错误
	return nil // 或 return fmt.Errorf("alter column not supported in MongoDB")
}

// FullDataTypeOf 返回字段完整数据类型（MongoDB 不适用，返回空字符串）
func (m *MongoMigrator) FullDataTypeOf(field *schema.Field) clause.Expr {
	return clause.Expr{}
}

// GetIndexes 返回模型的索引信息（MongoDB 不支持，返回空切片）
func (m *MongoMigrator) GetIndexes(s interface{}) ([]gorm.Index, error) {
	return []gorm.Index{}, nil
}

// CreateIndex 不支持
func (m *MongoMigrator) CreateIndex(value interface{}, name string) error {
	return nil
}

// DropIndex 不支持
func (m *MongoMigrator) DropIndex(value interface{}, name string) error {
	return nil
}

// HasIndex 检查索引是否存在
func (m *MongoMigrator) HasIndex(value interface{}, name string) bool {
	return true
}

// RenameIndex 不支持
func (m *MongoMigrator) RenameIndex(value interface{}, oldName, newName string) error {
	return nil
}
func (m *MongoMigrator) CreateConstraint(value interface{}, name string) error {
	// 类型断言获取 *gorm.DB
	_, ok := value.(*gorm.DB)
	if !ok {
		return fmt.Errorf("expected *gorm.DB, got %T", value)
	}

	// 原有逻辑继续使用 db
	// 示例：
	// constraintSQL, err := m.db.Migrator().(gorm.Migrator).BuildCreateConstraintSQL(name)
	// if err != nil {
	//     return err
	// }
	// return m.db.Exec(constraintSQL).Error

	return nil
}
func (m *MongoMigrator) CreateView(name string, option gorm.ViewOption) error {
	return nil // MongoDB 不支持视图或按需实现
}

func (m *MongoMigrator) DropView(name string) error {
	return nil // MongoDB 不支持视图
}

func (m *MongoMigrator) HasView(name string) bool {
	return false // MongoDB 不支持视图
}

// CurrentDatabase 返回当前数据库名（实现 gorm.Migrator 接口）
func (m *MongoMigrator) CurrentDatabase() string {
	// 假设 MongoDBDialector 中的 Database 字段保存了当前数据库名
	dialector, ok := m.db.Dialector.(*MongoDBDialector)
	if !ok {
		return ""
	}
	return dialector.Database
}

// DropConstraint 实现接口方法，MongoDB 不支持约束，所以不做任何操作
func (m *MongoMigrator) DropConstraint(tableName interface{}, constraintName string) error {
	return nil
}

// MongoDB 不需要自动建表，因此这里可留空或输出日志提示
func (m *MongoMigrator) AutoMigrate(values ...interface{}) error {
	return nil
}

// ColumnTypes 返回模型字段的列类型信息（模拟实现）
func (m *MongoMigrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	// 获取模型 schema
	stmt := &gorm.Statement{DB: m.db}
	if err := stmt.Parse(value); err != nil {
		return nil, err
	}

	var columnTypes []gorm.ColumnType
	for _, field := range stmt.Schema.Fields {
		columnTypes = append(columnTypes, &mongoColumnType{
			name:    field.Name,
			dbName:  field.DBName,
			sqlType: string(field.DataType),
		})
	}
	return columnTypes, nil
}

// 自定义 gorm.ColumnType 实现
type mongoColumnType struct {
	name    string
	dbName  string
	sqlType string
}

// 在 ColumnTypes 函数外定义 mongoColumnType 结构体（已存在）

// 补充实现 gorm.ColumnType 接口所需的 ColumnType 方法
func (c *mongoColumnType) ColumnType() (name string, ok bool) {
	return c.sqlType, true
}
func (c *mongoColumnType) Name() string             { return c.name }
func (c *mongoColumnType) OriginName() string       { return c.name }
func (c *mongoColumnType) DatabaseTypeName() string { return c.sqlType }
func (c *mongoColumnType) DatabaseTypeNameAndSize() (string, string) {
	return c.sqlType, ""
}
func (c *mongoColumnType) Length() (length int64, ok bool) { return 0, false }
func (c *mongoColumnType) PrecisionScale() (precision int64, scale int64, ok bool) {
	return 0, 0, false
}
func (c *mongoColumnType) Nullable() (nullable bool, ok bool) { return true, true }
func (c *mongoColumnType) Unique() (unique bool, ok bool)     { return false, true }
func (c *mongoColumnType) ScanType() reflect.Type             { return nil }
func (c *mongoColumnType) DefaultValue() (string, bool) {
	// 如果没有默认值，返回空字符串和 false
	return "", false
}
func (c *mongoColumnType) AutoIncrement() (auto bool, ok bool) { return false, true }
func (c *mongoColumnType) Comment() (string, bool) {
	return "", false // MongoDB 不支持字段注释，返回空值和 false
}
func (c *mongoColumnType) DecimalSize() (precision int64, scale int64, ok bool) {
	return 0, 0, false
}
func (c *mongoColumnType) PrimaryKey() (bool, bool) {
	return false, true
}
