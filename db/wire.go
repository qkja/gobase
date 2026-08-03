// Package db 提供 MongoDB 连接的 wire provider（公共基础设施）。
// 各服务在 wire 装配时直接复用 MongoSet，无需重复实现连接逻辑。
package db

import (
	"context"
	"time"

	"github.com/google/wire"
	gobasecfg "github.com/qkja/gobase/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongodbCfg MongoDB 连接配置（从 gobase 配置读取 database.mongodb 段）
type mongodbCfg struct {
	URI            string `yaml:"uri"`
	Database       string `yaml:"database"`
	ConnectTimeout string `yaml:"connect_timeout"` // "10s"
}

// connectTimeout 解析连接超时，默认 10s
func (c mongodbCfg) connectTimeout() time.Duration {
	if c.ConnectTimeout == "" {
		return 10 * time.Second
	}
	if d, err := time.ParseDuration(c.ConnectTimeout); err == nil {
		return d
	}
	return 10 * time.Second
}

// loadMongoCfg 从 gobase 配置读取 database.mongodb 段
func loadMongoCfg() mongodbCfg {
	cfg := mongodbCfg{URI: "mongodb://localhost:27017", Database: "identityhub"}
	_ = gobasecfg.GetValueObject("database.mongodb", &cfg)
	return cfg
}

// NewClient wire provider：连接 MongoDB（读 gobase 配置），返回连接好的 client。
func NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mc := loadMongoCfg()
	opts := options.Client().
		ApplyURI(mc.URI).
		SetConnectTimeout(mc.connectTimeout()).
		SetServerSelectionTimeout(5 * time.Second)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return c, nil
}

// NewDatabase wire provider：根据 client 返回库实例。
func NewDatabase(client *mongo.Client) (*mongo.Database, error) {
	return client.Database(loadMongoCfg().Database), nil
}

// MongoSet 公共 provider 集合，各服务 wire 装配时引入。
var MongoSet = wire.NewSet(NewClient, NewDatabase)
