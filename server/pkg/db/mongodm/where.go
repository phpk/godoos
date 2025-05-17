package mongodm

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm/clause"
)

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
			if nativeCond, ok := cond.(clause.Expr); ok {
				sqlStr := nativeCond.SQL
				vals := nativeCond.Vars

				// 简单处理：支持形如 "column = ?" 或带 OR 的多个条件
				orClauses := parseNativeSQLConditions(sqlStr, vals)

				if len(orClauses) > 0 {
					if _, exists := filter["$or"]; exists {
						// 如果已有 $or，合并进去（用于 JOIN、GROUP BY 等复杂场景）
						existingOr := filter["$or"].([]bson.M)
						filter["$or"] = append(existingOr, orClauses...)
					} else {
						filter["$or"] = orClauses
					}
				}
			}
		}
	}
	return filter
}

// convertSelectToBSON 将 Select clause 转换为 MongoDB 的 projection
func convertSelectToBSON(selectClause clause.Select) bson.M {
	projection := bson.M{}
	for _, col := range selectClause.Columns {
		fieldName := getColumnName(col)
		projection[fieldName] = 1 // 1 表示包含该字段
	}
	return projection
}

// 解析 SQL 片段，如 "username = ? OR email = ?"，返回对应的 bson.M 列表
// 支持更复杂的 SQL 表达式解析，包括：=, !=, >, <, >=, <=, IN, LIKE, IS NULL 等
func parseNativeSQLConditions(sql string, values []interface{}) []bson.M {
	var conditions []bson.M
	sql = strings.TrimSpace(sql)

	// 分割 OR 条件
	orParts := splitSQLCondition(sql, "OR")
	for _, orPart := range orParts {
		orPart = strings.TrimSpace(orPart)

		// 处理单个 OR 子句中的 AND 条件（可嵌套）
		andParts := splitSQLCondition(orPart, "AND")
		var andConditions []bson.M

		for _, andPart := range andParts {
			andPart = strings.TrimSpace(andPart)
			if andPart == "" {
				continue
			}

			cond := parseSingleCondition(andPart, values, len(conditions)+len(andConditions))
			if cond != nil {
				andConditions = append(andConditions, cond)
			}
		}

		if len(andConditions) > 1 {
			// 如果有多个 AND 条件，合并为一个 bson.M
			andFilter := bson.M{}
			for _, c := range andConditions {
				for k, v := range c {
					andFilter[k] = v
				}
			}
			conditions = append(conditions, andFilter)
		} else if len(andConditions) == 1 {
			conditions = append(conditions, andConditions[0])
		}
	}

	return conditions
}

// 分割 SQL 条件（支持 AND / OR）
func splitSQLCondition(sql string, separator string) []string {
	// 简化处理，避免破坏带引号的字符串中可能含有的 AND/OR
	// 实际生产环境可用 SQL 解析器或正则表达式加强匹配
	return strings.Split(strings.ToLower(sql), strings.ToLower(separator))
}

// 解析单个 SQL 条件，返回对应的 bson.M
func parseSingleCondition(part string, values []interface{}, index int) bson.M {
	part = strings.TrimSpace(part)

	// 跳过空条件
	if part == "" {
		return nil
	}

	// 匹配 IS NULL / IS NOT NULL
	if strings.Contains(strings.ToUpper(part), "IS NULL") {
		fieldName := strings.TrimSpace(strings.Split(strings.ToUpper(part), "IS NULL")[0])
		return bson.M{fieldName: nil}
	}
	if strings.Contains(strings.ToUpper(part), "IS NOT NULL") {
		fieldName := strings.TrimSpace(strings.Split(strings.ToUpper(part), "IS NOT NULL")[0])
		return bson.M{fieldName: bson.M{"$ne": nil}}
	}

	// 匹配 LIKE
	if strings.Contains(strings.ToUpper(part), "LIKE") {
		parts := strings.SplitN(part, "LIKE", 2)
		if len(parts) == 2 && index < len(values) {
			fieldName := strings.TrimSpace(parts[0])
			value := values[index]
			if strVal, ok := value.(string); ok {
				// 支持 %xxx% 的模糊匹配
				regexStr := strings.ReplaceAll(strVal, "%", ".*")
				return bson.M{fieldName: bson.M{"$regex": regexStr, "$options": "i"}}
			}
		}
	}

	// 匹配 IN
	if strings.Contains(strings.ToUpper(part), "IN") {
		parts := strings.SplitN(part, "IN", 2)
		if len(parts) == 2 && index < len(values) {
			fieldName := strings.TrimSpace(parts[0])
			if reflect.TypeOf(values[index]).Kind() == reflect.Slice {
				slice := reflect.ValueOf(values[index])
				inValues := make([]interface{}, slice.Len())
				for i := 0; i < slice.Len(); i++ {
					inValues[i] = slice.Index(i).Interface()
				}
				return bson.M{fieldName: bson.M{"$in": inValues}}
			}
		}
	}

	// 匹配比较运算符：>, >=, <, <=, !=
	comparisonOp := map[string]string{
		"!=": "$ne",
		">":  "$gt",
		">=": "$gte",
		"<":  "$lt",
		"<=": "$lte",
	}

	for opStr, mongoOp := range comparisonOp {
		if strings.Contains(part, opStr) {
			parts := strings.Split(part, opStr)
			if len(parts) >= 2 && index < len(values) {
				fieldName := strings.TrimSpace(parts[0])
				return bson.M{fieldName: bson.M{mongoOp: values[index]}}
			}
		}
	}

	// 默认：支持等值查询 =
	if strings.Contains(part, "=") {
		parts := strings.Split(part, "=")
		if len(parts) >= 2 && index < len(values) {
			fieldName := strings.TrimSpace(parts[0])
			return bson.M{fieldName: values[index]}
		}
	}

	return nil
}
