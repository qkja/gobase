// Package infra 提供基础设施（MongoDB / Redis）的统一初始化入口。
// 消费服务在启动阶段调用 Init 一次性初始化并持有，或通过 wire 装配 InfraSet 注入。
package infra

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/db"
	extendredis "github.com/qkja/gobase/extend/redis"
)

// Infra 统一持有已初始化的基础设施；未启用/未配置的组件为 nil。
type Infra struct {
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
	Redis       goredis.UniversalClient
}

// Init 读取 gobase 配置，初始化已启用的组件：
//   - MongoDB：配置存在 database.mongodb.uri 时连接（复用 db.NewClient/NewDatabase）；
//   - Redis：redis.enable=true 时连接（复用 extend/redis.NewClient，自动识别 standalone/sentinel/cluster）。
//
// 未启用/未配置的组件保持 nil，不建立默认地址连接。调用前需先加载配置（config.LoadConfig / gobasecfg.Init）。
func Init() (*Infra, error) {
	inf := &Infra{}

	if config.GetValueString("database.mongodb.uri") != "" {
		mc, err := db.NewClient()
		if err != nil {
			return nil, err
		}
		mdb, err := db.NewDatabase(mc)
		if err != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_ = mc.Disconnect(ctx)
			cancel()
			return nil, err
		}
		inf.MongoClient = mc
		inf.MongoDB = mdb
	}

	if config.GetValueBoolDefault("redis.enable", false) {
		rd, err := extendredis.NewClient()
		if err != nil {
			return nil, err
		}
		inf.Redis = rd
	}

	return inf, nil
}

// Close 关闭底层连接（幂等；nil 组件跳过）。返回首个非 nil 错误。
func (i *Infra) Close() error {
	if i == nil {
		return nil
	}
	var firstErr error
	if i.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := i.MongoClient.Disconnect(ctx)
		cancel()
		if err != nil {
			firstErr = err
		}
	}
	if i.Redis != nil {
		if err := i.Redis.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
