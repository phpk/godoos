package mongodm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
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

	// 注册 CURD 回调
	db.Callback().Create().Replace("gorm:create", d.createCallback)
	db.Callback().Query().Replace("gorm:query", d.queryCallback)
	db.Callback().Update().Replace("gorm:update", d.updateCallback)
	db.Callback().Delete().Replace("gorm:delete", d.deleteCallback)

	db.Dialector = d
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

// =======================
// CURD Callbacks
// =======================

func (d MongoDBDialector) createCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)

	// 获取模型字段信息
	stmt := db.Statement
	for _, field := range stmt.Schema.Fields {
		if field.AutoIncrementIncrement == 0 && field.TagSettings["AUTOINCREMENT"] == "true" {
			// 如果是 autoIncrement 主键且值为空，则生成新 ID
			nextID, err := GetNextID(d.Client, d.Database, collectionName)
			if err != nil {
				db.AddError(err)
				return
			}

			// 设置自增 ID 到模型中（使用新的 Set 接口）
			val := reflect.ValueOf(nextID).Convert(field.FieldType).Interface()
			if err := field.Set(context.Background(), stmt.ReflectValue, val); err != nil {
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

func (d MongoDBDialector) queryCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)
	// 构建聚合管道（传入 dest）
	dest := db.Statement.Dest
	pipeline := buildMongoPipeline(db.Statement.Clauses, dest)

	// 执行聚合查询
	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		db.AddError(err)
		return
	}

	// 将结果映射到目标结构体
	err = cursor.All(context.Background(), db.Statement.Dest)
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

	updateData := bson.M{"$set": db.Statement.Model}
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

// =======================
// 工具函数：将 Where/Join/GroupBy 等转换为聚合 Pipeline
// =======================
func buildMongoPipeline(clauses map[string]clause.Clause, dest interface{}) []bson.M {
	var pipeline []bson.M

	// 处理 WHERE 条件
	if whereClause, ok := clauses["WHERE"].Expression.(clause.Where); ok {
		filter := convertWhereToBSON(whereClause)
		if len(filter) > 0 {
			pipeline = append(pipeline, bson.M{"$match": filter})
		}
	}
	// 处理 JOIN
	if joinClause, ok := clauses["JOIN"].Expression.(clause.Join); ok {
		lookupStage := buildJoinPipeline(joinClause, dest)
		if lookupStage != nil {
			pipeline = append(pipeline, lookupStage)

			// INNER JOIN 处理
			if joinClause.Type == clause.InnerJoin {
				as := getNestedFieldName(dest, joinClause.Table.Name)
				if as == "" {
					as = strings.ToLower(joinClause.Table.Name) + "s"
				}
				unwindStage := bson.M{"$unwind": "$" + as}
				pipeline = append(pipeline, unwindStage)
			}

			// RIGHT JOIN 处理
			if joinClause.Type == clause.RightJoin {
				as := getNestedFieldName(dest, joinClause.Table.Name)
				if as == "" {
					as = strings.ToLower(joinClause.Table.Name) + "s"
				}

				// 展开嵌套字段
				unwindStage := bson.M{"$unwind": "$" + as}

				// 处理 NULL 值，确保右表记录始终存在
				addFieldsStage := bson.M{
					"$addFields": bson.M{
						as: bson.M{
							"$ifNull": []interface{}{
								"$" + as,
								bson.M{"$literal": []interface{}{}},
							},
						},
					},
				}

				pipeline = append(pipeline, unwindStage)
				pipeline = append(pipeline, addFieldsStage)
			}
		}
	}

	// 可选：处理 GroupBy、Select、Sort 等聚合操作
	// 示例：GroupBy + Count
	if groupByClause, ok := clauses["GROUP BY"].Expression.(clause.GroupBy); ok {
		groupFields := bson.M{"_id": nil}
		for _, item := range groupByClause.Columns {
			groupFields["_id"] = "$" + getColumnName(item)
		}
		groupFields["count"] = bson.M{"$sum": 1}
		pipeline = append(pipeline, bson.M{"$group": groupFields})
	}

	return pipeline
}
func buildJoinPipeline(join clause.Join, dest interface{}) bson.M {
	// 默认值
	localField := ""
	foreignField := ""
	//as := join.As
	from := join.Table.Name

	// 遍历 ON 条件
	for _, expr := range join.ON.Exprs {
		if eq, ok := expr.(clause.Eq); ok {
			// 处理左侧 Column
			if lc, lok := eq.Column.(clause.Column); lok {
				localField = lc.Name
			}

			// 处理右侧 Value（假设是 string 类型的字段名）
			if strVal, ok := eq.Value.(string); ok {
				foreignField = strVal
			}
		}
	}

	// 如果无法提取字段，则返回空
	if localField == "" || foreignField == "" {
		return nil
	}
	// 从 dest 结构体中获取嵌套字段名（如 Posts）
	as := getNestedFieldName(dest, from)
	if as == "" {
		as = strings.ToLower(from) + "s" // 默认命名规则
	}
	// 返回 $lookup stage
	return bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	}
}
func buildUnwindPipeline(as string) bson.M {
	return bson.M{"$unwind": "$" + as}
}
func getNestedFieldName(dest interface{}, collectionName string) string {
	destType := reflect.TypeOf(dest).Elem()

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tag := field.Tag.Get("gorm")

		if tag != "" {
			gormTag := parseGormTag(tag)
			if gormTag["foreignKey"] == collectionName {
				return field.Name
			}
		}
	}

	return ""
}

// =======================
// 工具函数：将 Where 转换为 BSON
// =======================
func convertWhereToBSON(where clause.Where) bson.M {
	filter := bson.M{}

	for _, cond := range where.Exprs {
		switch expr := cond.(type) {
		case clause.Eq:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = expr.Value
			}
		case clause.Neq:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$ne": expr.Value}
			}
		case clause.Gt:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$gt": expr.Value}
			}
		case clause.Gte:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$gte": expr.Value}
			}
		case clause.Lt:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$lt": expr.Value}
			}
		case clause.Lte:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$lte": expr.Value}
			}
		case clause.IN:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$in": expr.Values}
			}
		case clause.Like:
			if col, ok := expr.Column.(clause.Column); ok {
				key := getColumnName(col)
				filter[key] = bson.M{"$regex": expr.Value}
			}
		default:
			// 可选：处理未知条件或记录日志
		}
	}
	return filter
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
