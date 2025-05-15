// join.go
package mongodm

import (
	"gorm.io/gorm/clause"
)

// MongoJoin 定义 MongoDB $lookup 所需的信息
type MongoJoin struct {
	Table     string        // 要连接的集合名（from）
	On        clause.Column // 当前集合中的字段（localField）
	RefColumn string        // 参考集合中的字段（foreignField）
	As        string        // 存储结果的字段名（as）
	IsUnwind  bool          // 是否执行 $unwind 展开
}

// Name 实现 clause.Expression 接口
func (j MongoJoin) Name() string {
	return "MongoJoin"
}

// Build 构建对应的 BSON 表达式（用于调试或 Explain）
func (j MongoJoin) Build(builder clause.Builder) {
	builder.WriteString("JOIN ")
	builder.WriteString(j.Table)
	builder.WriteString(" ON ")
	builder.WriteQuoted(j.On)
	builder.WriteString(" = ")
	builder.WriteString(j.RefColumn)
}
