package mongodm

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d MongoDBDialector) queryCallback(db *gorm.DB) {
	collectionName := db.Statement.Table
	coll := db.ConnPool.(*mongoConnPool).Collection(collectionName)
	dest := db.Statement.Dest
	stmt := db.Statement
	// log.Printf("stmt.Clauses: %+v\n", stmt.Clauses)
	//log.Printf("stmt.Clauses: %+v\n", stmt.Clauses["SELECT"].Expression)
	// 检查是否有 GROUP BY 或 JOIN
	hasGroupBy := false
	if _, ok := stmt.Clauses["GROUP BY"]; ok {
		hasGroupBy = true
	}
	if _, ok := stmt.Clauses["JOIN"].Expression.(clause.Join); ok {
		// 所有类型的 JOIN（包括 LEFT JOIN）都触发聚合查询
		hasGroupBy = true
	}

	var limit int
	if limitClause, ok := stmt.Clauses["LIMIT"].Expression.(clause.Limit); ok {
		if limitClause.Limit != nil {
			limit = *limitClause.Limit
		} else {
			// First() 等隐式调用 LIMIT 1
			limit = 1
		}
	}

	// 构建查询条件
	whereClause, hasWhere := stmt.Clauses["WHERE"].Expression.(clause.Where)
	//  log.Printf("Where: %+v\n", whereClause)
	filter := bson.M{}
	if hasWhere {
		filter = convertWhereToBSON(whereClause)

	}
	// 构建投影字段（SELECT）
	var projection bson.M
	if isCountQuery(stmt.Clauses["SELECT"].Expression) {
		count, err := coll.CountDocuments(context.Background(), filter)
		if err != nil {
			db.AddError(fmt.Errorf("failed to count documents: %w", err))
			return
		}
		err = setValue(db.Statement.Dest, count)
		if err != nil {
			db.AddError(err)
			return
		}
		db.RowsAffected = count
		return
	}
	//log.Printf("stmt.Clauses: %+v\n", stmt.Clauses["SELECT"].Expression)
	if selectClause, ok := stmt.Clauses["SELECT"].Expression.(clause.Select); ok {
		// 判断是否是 count 查询
		exprStr := fmt.Sprintf("%v", selectClause)
		log.Printf("exprStr: %s\n", exprStr)

		projection = convertSelectToBSON(selectClause)
	}

	if hasGroupBy {
		pipeline := buildMongoPipeline(stmt.Clauses, dest)

		cursor, err := coll.Aggregate(context.Background(), pipeline)
		if err != nil {
			db.AddError(err)
			return
		}
		defer cursor.Close(context.Background())
		if limit == 1 {
			// 聚合结果中取第一条
			if cursor.Next(context.Background()) {
				if err := cursor.Decode(dest); err != nil {
					db.AddError(err)
				}
			} else {
				db.AddError(gorm.ErrRecordNotFound)
			}
		} else {
			// 获取所有聚合结果
			if err := cursor.All(context.Background(), dest); err != nil {
				db.AddError(err)
			}
		}
		return
	}

	if limit == 1 {
		// 处理 .First()
		opts := options.FindOne().SetProjection(projection)
		var result bson.M
		err := coll.FindOne(context.Background(), filter, opts).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				db.AddError(gorm.ErrRecordNotFound)
			} else {
				db.AddError(err)
			}
			return
		}

		// 将结果映射到原始 dest（例如 *User）
		//log.Printf("first dest: %+v\n", dest)
		if err := mapBSONToStruct(result, dest); err != nil {
			db.AddError(err)
		}
	} else {
		// 处理 .Find()
		opts := options.Find().SetProjection(projection)
		cursor, err := coll.Find(context.Background(), filter, opts)
		if err != nil {
			db.AddError(err)
			return
		}
		defer cursor.Close(context.Background())

		err = cursor.All(context.Background(), dest)
		if err != nil {
			db.AddError(err)
		}
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
	from := join.Table.Name

	for _, expr := range join.ON.Exprs {
		if eq, ok := expr.(clause.Eq); ok {
			if lc, lok := eq.Column.(clause.Column); lok {
				localField = lc.Name
			}
			if strVal, ok := eq.Value.(string); ok {
				foreignField = strVal
			}
		}
	}

	if localField == "" || foreignField == "" {
		return nil
	}

	as := getNestedFieldName(dest, from)
	if as == "" {
		as = strings.ToLower(from) + "s"
	}

	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	}

	// 区分不同 JOIN 类型
	switch join.Type {
	case clause.InnerJoin:
		unwindStage := bson.M{"$unwind": "$" + as}
		return bson.M{
			"$and": []bson.M{lookupStage, unwindStage},
		}
	case clause.RightJoin:
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
		unwindStage := bson.M{"$unwind": "$" + as}
		return bson.M{
			"$and": []bson.M{lookupStage, addFieldsStage, unwindStage},
		}
	case clause.LeftJoin:
		// 左连接只需保留空数组作为默认值，不强制展开
		return lookupStage
	default:
		return lookupStage
	}
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
func isCountQuery(expr clause.Expression) bool {
	if exprStmt, ok := expr.(clause.Expr); ok {
		exprStr := strings.ToLower(exprStmt.SQL)
		//log.Printf("exprStr: %s\n", exprStr)

		if strings.Contains(exprStr, "count(") {
			//log.Println("Intercepted COUNT(*) query")
			return true
		}
	}
	return false
}

// 支持 *int64、*int、**int64、***int64 等各种嵌套指针
func setValue(dest interface{}, count int64) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer, got %s", destVal.Kind())
	}

	elem := destVal.Elem()
	for elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		elem = elem.Elem()
	}

	//log.Printf("Before SetInt: %v (type: %s)", elem.Interface(), elem.Kind())

	switch elem.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		elem.SetInt(count)
	default:
		return fmt.Errorf("unsupported destination type: %s", elem.Kind())
	}

	//log.Printf("After SetInt: %v", elem.Interface())
	return nil
}
