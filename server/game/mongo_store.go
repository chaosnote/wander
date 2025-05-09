package game

import (
	"context"

	"github.com/chaosnote/wander/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type temp struct {
	Name string
	Age  int
}

type MongoStore interface{}

type mongo_store struct {
	utils.LogStore

	record *mongo.Collection
}

func (s *mongo_store) RecordSave(uid string, model any) (e error) {
	// 定義查詢條件 (根據 _id 欄位)
	filter := bson.D{{Key: "_id", Value: uid}}

	// 定義要更新或插入的資料 ($set 用於更新，如果不存在則插入整個文檔)
	update := bson.D{{Key: "$set", Value: model}}

	// 當無資料時，則自動更新
	opts := options.Update().SetUpsert(true)
	_, e = s.record.UpdateOne(context.Background(), filter, update, opts)
	if e != nil {
		s.Error(e)
		return
	}
	return
}

// RecordLoad
//
//	var model struct
//	RecordLoad(uid, &model)
func (s *mongo_store) RecordLoad(uid string, model any) (e error) {
	filter := bson.D{{Key: "_id", Value: uid}}
	e = s.record.FindOne(context.TODO(), filter).Decode(model)
	return
}

//-----------------------------------------------

func NewMongoStore() MongoStore {
	di := utils.GetDI()

	client := di.MustGet(utils.SERVICE_MONGO).(*mongo.Client)
	store := &mongo_store{
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),

		record: client.Database("GameDB").Collection("Record"),
	}

	store.RecordSave("123", temp{Name: "chris", Age: 20})
	var m temp
	e := store.RecordLoad("123", &m)
	if e != nil {
		store.Error(e)
	}
	store.Debug(utils.LogFields{"name": m.Name, "age": m.Age})

	return store
}
