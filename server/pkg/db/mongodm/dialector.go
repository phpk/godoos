package mongodm

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// MongoDBDialector 实现 GORM 的 Dialector 接口
type MongoDBDialector struct {
	URI      string        // MongoDB 连接 URI
	Database string        // 默认数据库
	Client   *mongo.Client // MongoDB 客户端
}

// Name 返回该驱动的名称
func (d MongoDBDialector) Name() string {
	return "mongodb"
}

// Initialize 实际连接 MongoDB
func (d *MongoDBDialector) Initialize(db *gorm.DB) error {
	clientOptions := options.Client().ApplyURI(d.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// 设置连接健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	// 保存client到结构体
	d.Client = client

	// 将 client 存入 GORM 的 ConnPool
	db.ConnPool = &mongoConnPool{db: client.Database(d.Database)}
	// 注册自定义 Clause Builders
	db.ClauseBuilders["WHERE"] = d.buildWhereClause
	db.ClauseBuilders["LIMIT"] = d.buildLimitClause
	db.ClauseBuilders["ORDER BY"] = d.buildOrderByClause
	db.ClauseBuilders["SELECT"] = d.buildSelectClause
	db.ClauseBuilders["COUNT"] = d.buildCountClause
	// 注册 CURD 回调
	db.Callback().Create().Replace("gorm:create", d.createCallback)
	db.Callback().Query().Replace("gorm:query", d.queryCallback)
	db.Callback().Update().Replace("gorm:update", d.updateCallback)
	db.Callback().Delete().Replace("gorm:delete", d.deleteCallback)

	db.Dialector = d
	// 注册查询完成后的回调函数
	// db.Callback().Query().After("gorm:query").Register("track_count", func(db *gorm.DB) {
	// 	// 检查是否为 count 查询
	// 	if isCountQuery(db.Statement.Clauses["SELECT"].Expression) {
	// 		// 获取 count 值
	// 		if dest, ok := db.Statement.Dest.(*int64); ok {
	// 			log.Printf("===Count result: %d", *dest)
	// 		} else {
	// 			log.Printf("db.Statement.Dest is not a pointer to int64")
	// 		}
	// 	}
	// })
	return nil
}

// Migrator 返回一个迁移器（MongoDB 不需要迁移）
func (d MongoDBDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return &MongoMigrator{db: db}
}

// DataTypeOf 返回字段在数据库中的数据类型（MongoDB 不适用，返回空字符串）
func (d MongoDBDialector) DataTypeOf(field *schema.Field) string {
	return ""
}

// DefaultValueOf 返回字段的默认值表达式（MongoDB 不适用，返回 nil）
func (d MongoDBDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return nil
}

// BindVarTo 处理变量绑定（模拟为 ? 占位符）
func (d MongoDBDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteString("?")
}

// QuoteTo 引用标识符（如字段名），用于防止关键字冲突
func (d MongoDBDialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString("\"" + str + "\"")
}

// Explain 格式化输出 SQL 语句（此处直接返回原样）
func (d MongoDBDialector) Explain(sql string, vars ...interface{}) string {
	return sql
}

// =======================
// 自定义 Clause Builders
// =======================

func (d MongoDBDialector) buildWhereClause(c clause.Clause, builder clause.Builder) {
	if where, ok := c.Expression.(clause.Where); ok {
		builder.WriteString(" WHERE ")
		for i, cond := range where.Exprs {
			if i > 0 {
				builder.WriteString(" AND ")
			}
			cond.Build(builder)
		}
	}
}

func (d MongoDBDialector) buildLimitClause(c clause.Clause, builder clause.Builder) {
	if limit, ok := c.Expression.(clause.Limit); ok && limit.Limit != nil {
		builder.WriteString(fmt.Sprintf(" LIMIT %d", *limit.Limit))
		if limit.Offset > 0 {
			builder.WriteString(fmt.Sprintf(" OFFSET %d", limit.Offset))
		}
	}
}

func (d MongoDBDialector) buildOrderByClause(c clause.Clause, builder clause.Builder) {
	if orderBy, ok := c.Expression.(clause.OrderBy); ok {
		builder.WriteString(" ORDER BY ")
		for i, item := range orderBy.Columns {
			if i > 0 {
				builder.WriteString(", ")
			}
			d.QuoteTo(builder, item.Column.Name)
			if !item.Desc {
				builder.WriteString(" ASC")
			} else {
				builder.WriteString(" DESC")
			}
		}
	}
}
func (d MongoDBDialector) buildCountClause(c clause.Clause, builder clause.Builder) {
	builder.WriteString(" COUNT(*) ")
}

// =======================
// CURD Callbacks
// =======================

func (d MongoDBDialector) createCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)

	// 获取模型字段信息
	stmt := db.Statement
	for _, field := range stmt.Schema.Fields {
		// log.Printf("TagSettings: %+v\n", field.TagSettings)
		// log.Printf("FieldType: %+v\n", field.AutoIncrementIncrement)
		if field.TagSettings["AUTOINCREMENT"] == "true" {
			// 如果是 autoIncrement 主键且值为空，则生成新 ID
			nextID, err := GetNextID(context.Background(), d.Client, d.Database, collectionName)
			if err != nil {
				db.AddError(err)
				return
			}
			//log.Printf("Got next ID: %d for collection: %s", nextID, collectionName)
			// 设置自增 ID 到模型中（使用新的 Set 接口）
			val := reflect.ValueOf(nextID).Convert(field.FieldType).Interface()
			if err := field.Set(context.Background(), stmt.ReflectValue, val); err != nil {
				db.AddError(err)
				return
			}
		}
	}
	// 第二步：处理 default tag
	// 处理 default tag（包括 CURRENT_TIMESTAMP）
	now := time.Now()
	for _, field := range stmt.Schema.Fields {
		defaultValueStr := field.TagSettings["DEFAULT"]
		if defaultValueStr == "" {
			continue
		}

		fieldValue := field.ReflectValueOf(context.Background(), stmt.ReflectValue)
		if !fieldValue.IsValid() || isZeroValue(fieldValue) {
			var defaultVal interface{}
			if defaultValueStr == "CURRENT_TIMESTAMP" {
				if field.FieldType == reflect.TypeOf(time.Time{}) {
					defaultVal = now
				} else {
					db.AddError(fmt.Errorf("CURRENT_TIMESTAMP can only be used with time.Time fields"))
					return
				}
			} else {
				var err error
				defaultVal, err = parseDefaultValue(field.FieldType, defaultValueStr)
				if err != nil {
					db.AddError(fmt.Errorf("failed to parse default value for %s: %w", field.Name, err))
					return
				}
			}

			if err := field.Set(context.Background(), stmt.ReflectValue, defaultVal); err != nil {
				db.AddError(err)
				return
			}
		}
	}

	_, err := coll.InsertOne(context.Background(), stmt.Model)
	if err != nil {
		db.AddError(err)
	}
}

func (d MongoDBDialector) updateCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)

	var filter bson.M
	if db.Statement.Clauses != nil {
		whereClause, ok := db.Statement.Clauses["WHERE"].Expression.(clause.Where)
		if ok {
			filter = convertWhereToBSON(whereClause)
		}
	}

	// 构建更新数据
	updateData := bson.M{"$set": db.Statement.Model}

	// 检查是否有字段标记为 CURRENT_TIMESTAMP
	updateMap := updateData["$set"].(map[string]interface{})
	for _, field := range db.Statement.Schema.Fields {
		tagDefault := field.TagSettings["DEFAULT"]
		if tagDefault == "CURRENT_TIMESTAMP" {
			if fieldValue := field.ReflectValueOf(context.Background(), db.Statement.ReflectValue); fieldValue.IsValid() {
				if fieldValue.IsZero() || isZeroValue(fieldValue) {
					updateMap[field.DBName] = time.Now()
				}
			} else {
				updateMap[field.DBName] = time.Now()
			}
		}
	}

	_, err := coll.UpdateMany(context.Background(), filter, updateData)
	if err != nil {
		db.AddError(err)
	}
}

func (d MongoDBDialector) deleteCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)

	var filter bson.M
	if db.Statement.Clauses != nil {
		whereClause, ok := db.Statement.Clauses["WHERE"].Expression.(clause.Where)
		if ok {
			filter = convertWhereToBSON(whereClause)
		}
	}

	_, err := coll.DeleteMany(context.Background(), filter)
	if err != nil {
		db.AddError(err)
	}
}

// 辅助函数获取字段名（带表前缀）
func getColumnName(column clause.Column) string {
	if column.Table == "" {
		return column.Name
	}
	return fmt.Sprintf("%s.%s", column.Table, column.Name)
}
func (d *MongoDBDialector) Begin(db *gorm.DB) (err error) {
	session, err := d.Client.StartSession()
	if err != nil {
		return err
	}
	db.InstanceSet("mongo:session", session)
	session.StartTransaction()
	return nil
}

func (d *MongoDBDialector) Commit(db *gorm.DB) (err error) {
	session, ok := db.InstanceGet("mongo:session")
	if !ok || session == nil {
		return nil
	}
	sess := session.(mongo.Session)
	err = sess.CommitTransaction(context.Background())
	sess.EndSession(context.Background())
	return err
}

func (d *MongoDBDialector) Rollback(db *gorm.DB) (err error) {
	session, ok := db.InstanceGet("mongo:session")
	if !ok || session == nil {
		return nil
	}
	sess := session.(mongo.Session)
	err = sess.AbortTransaction(context.Background())
	sess.EndSession(context.Background())
	return err
}

// buildSelectClause 构建 SELECT 投影表达式
func (d MongoDBDialector) buildSelectClause(c clause.Clause, builder clause.Builder) {
	if selectClause, ok := c.Expression.(clause.Select); ok {
		builder.WriteString("SELECT ")
		for i, column := range selectClause.Columns {
			if i > 0 {
				builder.WriteString(", ")
			}
			d.QuoteTo(builder, column.Name)
		}
	}
}
