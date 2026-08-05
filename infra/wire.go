package infra

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/qkja/gobase/db"
	extendredis "github.com/qkja/gobase/extend/redis"
)

// NewInfra wire provider：组合已初始化的 MongoDB client/database 与 Redis client。
func NewInfra(mc *mongo.Client, mdb *mongo.Database, rds goredis.UniversalClient) *Infra {
	return &Infra{MongoClient: mc, MongoDB: mdb, Redis: rds}
}

// InfraSet wire provider 集合：复用 db.MongoSet（MongoDB）与 extend/redis.NewClient（Redis）。
var InfraSet = wire.NewSet(NewInfra, db.MongoSet, extendredis.NewClient)
