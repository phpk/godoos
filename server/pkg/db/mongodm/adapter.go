// adapter.go
package mongodm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
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
func GetNextID(client *mongo.Client, dbName, collectionName string) (int64, error) {
	counterColl := client.Database(dbName).Collection("counters")

	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result struct {
		ID  string `bson:"_id"`
		Seq int64  `bson:"seq"`
	}

	err := counterColl.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}
