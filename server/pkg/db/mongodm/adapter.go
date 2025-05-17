// adapter.go
package mongodm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConnPool struct {
	db *mongo.Database
}
type ConnPoolStats struct {
	TotalConn   int
	InUse       int
	Idle        int
	WaitCount   int64
	WaitTime    time.Duration
	MaxIdleTime time.Duration
}
type Counter struct {
	ID  string `bson:"_id"` // collection name
	Seq int64  `bson:"seq"` // 当前序列值
}

func (p *mongoConnPool) GetDB() *mongo.Database {
	return p.db
}
func (p *mongoConnPool) Collection(name string) *mongo.Collection {
	return p.db.Collection(name)
}
func (p *mongoConnPool) Conn(ctx context.Context) (driver.Conn, error) {
	return &mongoConn{}, nil
}
func (p *mongoConnPool) PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error) {
	return nil, nil // MongoDB 不使用 Prepare/Stmt，直接返回 nil
}

// 在 mongoConn 结构体后添加以下方法即可修复编译错误

func (c *mongoConn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("Prepare not implemented")
}

func (p *mongoConnPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// MongoDB 不支持原生 SQL，这里可以返回 nil 和 nil 表示忽略执行
	// 或者根据业务需求做具体处理（如日志记录、错误提示等）
	return nil, nil
}

func (p *mongoConnPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return &sql.Rows{}, nil
}

func (p *mongoConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

type mongoTx struct {
	session mongo.Session // 现在 session 是 *mongo.Session 类型
}

func (p *mongoConnPool) BeginTx(ctx context.Context, opts interface{}) (interface{}, error) {
	session, err := p.db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	session.StartTransaction()
	// 修改此处：取 session 的地址，使其成为 *mongo.Session 类型
	return &mongoTx{session: session}, nil
}
func (t *mongoTx) Commit() error {
	return t.session.CommitTransaction(context.Background())
}

func (t *mongoTx) Rollback() error {
	return t.session.AbortTransaction(context.Background())
}
func (p *mongoConnPool) Ping(ctx context.Context) error {
	return p.db.Client().Ping(ctx, nil)
}

func (p *mongoConnPool) Stats() ConnPoolStats {
	return ConnPoolStats{}
}
func (p *mongoConnPool) Close() {}

type mongoConn struct{}

func (c *mongoConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("PrepareContext not implemented")
}

func (c *mongoConn) Close() error {
	return nil
}

func (c *mongoConn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("Begin not implemented")
}

// CommitTx 提交事务
func (p *mongoConnPool) CommitTx(ctx context.Context) error {
	session, ok := contextFromDB(ctx)
	if !ok {
		return nil
	}
	defer session.EndSession(ctx)
	return session.CommitTransaction(ctx)
}

// RollbackTx 回滚事务
func (p *mongoConnPool) RollbackTx(ctx context.Context) error {
	session, ok := contextFromDB(ctx)
	if !ok {
		return nil
	}
	defer session.EndSession(ctx)
	return session.AbortTransaction(ctx)
}

// 从上下文中获取 session（需配合 GORM 上下文管理）
func contextFromDB(ctx context.Context) (mongo.Session, bool) {
	value := ctx.Value("mongo:session")
	if value == nil {
		return nil, false
	}
	session, ok := value.(mongo.Session)
	return session, ok
}
func InitCounter(client *mongo.Client, dbName, collectionName string) error {
	counterColl := client.Database(dbName).Collection("counters")
	_, err := counterColl.InsertOne(context.Background(), bson.M{
		"_id": collectionName,
		"seq": int64(0),
	})
	return err
}

// CounterDocument 用于表示 counters 集合中的文档
type CounterDocument struct {
	ID  string `bson:"_id"` // 对应 collection name
	Seq int64  `bson:"seq"` // 序列号
}

// GetNextID 获取指定 collection 的下一个自增 ID
func GetNextID(ctx context.Context, client *mongo.Client, dbName, collectionName string) (int64, error) {
	counterColl := client.Database(dbName).Collection("counters")

	// 创建 filter 和 update
	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	// 设置 FindOneAndUpdate 选项：返回更新后的文档
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result CounterDocument

	err := counterColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果计数器不存在，自动创建初始值
			_, insertErr := counterColl.InsertOne(ctx, bson.M{
				"_id": collectionName,
				"seq": int64(1),
			})
			if insertErr != nil {
				return 0, fmt.Errorf("failed to initialize counter for %s: %w", collectionName, insertErr)
			}
			return 1, nil
		}
		// 其他错误
		return 0, fmt.Errorf("failed to get next id for %s: %w", collectionName, err)
	}

	return result.Seq, nil
}

// 判断是否是零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// 将字符串解析为对应类型
func parseDefaultValue(typ reflect.Type, valueStr string) (interface{}, error) {
	// 新增 CURRENT_TIMESTAMP 支持
	if valueStr == "CURRENT_TIMESTAMP" {
		if typ.Kind() == reflect.Struct && typ == reflect.TypeOf(time.Time{}) {
			return time.Now(), nil
		}
		return nil, fmt.Errorf("CURRENT_TIMESTAMP only supported for time.Time fields")
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		_, err := fmt.Sscanf(valueStr, "%d", &val)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(val).Convert(typ).Interface(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		_, err := fmt.Sscanf(valueStr, "%d", &val)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(val).Convert(typ).Interface(), nil

	case reflect.String:
		return valueStr, nil

	case reflect.Bool:
		switch valueStr {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid boolean value: %s", valueStr)
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", typ.Kind())
	}
}
func mapBSONToStruct(src bson.M, dst interface{}) error {
	data, err := bson.Marshal(src)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, dst)
}
