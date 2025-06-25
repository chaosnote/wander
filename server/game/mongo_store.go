package game

import (
	"context"

	"github.com/chaosnote/wander/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type OutputMongoStore interface {
	RecordSave(uid string, model any) (e error)
	// RecordLoad
	//
	//	var model struct
	//	RecordLoad(uid, &model)
	RecordLoad(uid string, model any) (e error)
}

type MongoStore interface {
	OutputMongoStore
}

type mongo_store struct {
	logger *zap.Logger

	record *mongo.Collection
}

func (s *mongo_store) RecordSave(uid string, model any) (e error) {
	const msg = "RecordSave"

	// 定義查詢條件 (根據 _id 欄位)
	filter := bson.D{{Key: "_id", Value: uid}}

	// 定義要更新或插入的資料 ($set 用於更新，如果不存在則插入整個文檔)
	update := bson.D{{Key: "$set", Value: model}}

	// 當無資料時，則自動更新
	opts := options.Update().SetUpsert(true)
	_, e = s.record.UpdateOne(context.Background(), filter, update, opts)
	if e != nil {
		s.logger.Error(msg, zap.Error(e))
		return
	}
	return
}

func (s *mongo_store) RecordLoad(uid string, model any) (e error) {
	filter := bson.D{{Key: "_id", Value: uid}}
	e = s.record.FindOne(context.TODO(), filter).Decode(model)
	return
}

//-----------------------------------------------

func NewMongoStore() MongoStore {
	di := utils.GetDI()

	client := di.MustGet(SERVICE_MONGO).(*mongo.Client)
	store := &mongo_store{
		logger: di.MustGet(LOGGER_SYSTEM).(*zap.Logger),

		record: client.Database("GameDB").Collection("Record"),
	}

	return store
}
