package dbfactory

import (
	"context"
	"reflect"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

// MongoDB工厂实现
type MongoDBFactory struct{}

func NewMongoDBFactory() *MongoDBFactory {
	return &MongoDBFactory{}
}

func (f *MongoDBFactory) CreateConnection() DatabaseConnection {
	return MongoDB
}

func (f *MongoDBFactory) CreateRepository() BaseRepository {
	return &mongoRepository{
		client: MongoDB,
	}
}

type mongoRepository struct {
	client *qmgo.Database
}

func (r *mongoRepository) Create(table string, entity interface{}) error {
	_, err := r.client.Collection(table).InsertOne(context.Background(), entity)
	return err
}
func (r *mongoRepository) Count(table string, conditions map[string]interface{}) (int64, error) {
	// 使用 CountDocuments 替代 Count
	count, err := r.client.Collection(table).Find(context.Background(), bson.M(conditions)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *mongoRepository) GetByID(table string, id interface{}, result interface{}) error {
	return r.client.Collection(table).Find(context.Background(), map[string]interface{}{"_id": id}).One(result)
}

func (r *mongoRepository) Update(table string, entity interface{}) error {
	return r.client.Collection(table).UpdateOne(context.Background(), map[string]interface{}{"_id": getIdFromEntity(entity)}, entity)
}

func (r *mongoRepository) Delete(table string, id interface{}) error {
	return r.client.Collection(table).Remove(context.Background(), map[string]interface{}{"_id": id})
}

func getIdFromEntity(entity interface{}) interface{} {
	// 获取entity的反射值
	val := reflect.ValueOf(entity)

	// 确保传入的是一个结构体指针
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil
	}

	// 获取结构体的反射类型
	elemType := val.Elem().Type()

	// 查找名为"_id"的字段
	field, found := elemType.FieldByName("_id")
	if !found {
		return nil
	}

	// 获取字段的值
	id := val.Elem().FieldByName(field.Name)

	// 返回字段的值
	return id.Interface()
}

func (r *mongoRepository) BatchCreate(table string, entities []interface{}) error {
	_, err := r.client.Collection(table).InsertMany(context.Background(), entities)
	return err
}

func (r *mongoRepository) BatchDelete(table string, ids []interface{}) error {
	_, err := r.client.Collection(table).RemoveAll(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	return err
}

func (r *mongoRepository) JoinQuery(mainTable string, joinTables []JoinTable, conditions map[string]interface{}, result interface{}) error {
	// MongoDB需要特殊处理联表查询
	// 这里实现$lookup聚合查询
	pipeline := []bson.M{
		{"$match": conditions},
	}
	for _, join := range joinTables {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         join.Table,
				"localField":   "_id",
				"foreignField": "_id",
				"as":           join.Table,
			},
		})
	}
	return r.client.Collection(mainTable).Aggregate(context.Background(), pipeline).All(result)
}

func (r *mongoRepository) GetOne(table string, conditions map[string]interface{}, result interface{}) error {
	return r.client.Collection(table).Find(context.Background(), conditions).One(result)
}

func (r *mongoRepository) GetList(table string, conditions map[string]interface{}, result interface{}) error {
	return r.client.Collection(table).Find(context.Background(), conditions).All(result)
}

func (r *mongoRepository) GetPage(table string, conditions map[string]interface{}, pageParams PageParams, result *PageResult) error {
	// 获取总数
	total, err := r.client.Collection(table).Find(context.Background(), conditions).Count()
	if err != nil {
		return err
	}

	// 设置分页选项（使用 qmgo 的链式调用）
	findOptions := r.client.Collection(table).Find(context.Background(), conditions).
		Skip(int64((pageParams.Page - 1) * pageParams.PageSize)).
		Limit(int64(pageParams.PageSize))

	if pageParams.OrderBy != "" {
		findOptions = findOptions.Sort(pageParams.OrderBy)
	}

	// 执行查询
	var data []interface{} // 假设 result.Data 是一个切片
	if err := findOptions.All(&data); err != nil {
		return err
	}

	// 将数据赋值给 result.Data
	result.Data = data
	result.Total = total
	return nil
}

func (r *mongoRepository) AutoMigrate(models ...interface{}) error {
	// MongoDB不需要显式迁移，返回nil即可
	return nil
}

func (r *mongoRepository) WithTransaction(ctx context.Context, opts *TxOptions, fn func(repo BaseRepository) error) error {
	// // 获取 qmgo.Client
	// client := r.client.Client()
	// if client == nil {
	// 	return fmt.Errorf("client is nil")
	// }

	// // 定义事务回调函数
	// callback := func(sessCtx context.Context) (interface{}, error) {
	// 	// 创建一个新的 mongoRepository 实例，使用 sessCtx 作为上下文
	// 	txRepo := &mongoRepository{client: r.client}
	// 	return nil, fn(txRepo)
	// }

	// // 设置事务选项
	// transactionOpts := options.Transaction()
	// switch opts.Isolation {
	// case LevelReadUncommitted:
	// 	transactionOpts = transactionOpts.SetReadConcern(options.ReadConcern{}.SetLevel("local"))
	// case LevelReadCommitted:
	// 	transactionOpts = transactionOpts.SetReadConcern(options.ReadConcern{}.SetLevel("majority"))
	// case LevelRepeatableRead:
	// 	transactionOpts = transactionOpts.SetReadConcern(options.ReadConcern{}.SetLevel("snapshot"))
	// case LevelSerializable:
	// 	transactionOpts = transactionOpts.SetReadConcern(options.ReadConcern{}.SetLevel("linearizable"))
	// }
	// if opts.ReadOnly {
	// 	transactionOpts = transactionOpts.SetReadOnly(true)
	// }

	// // 执行事务
	// result, err := client.DoTransaction(ctx, callback, transactionOpts)
	// if err == qmgo.ErrTransactionRetry {
	// 	return r.WithTransaction(ctx, opts, fn)
	// }

	// return err
	return nil
}
