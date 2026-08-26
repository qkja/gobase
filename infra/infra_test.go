package infra

import (
	"testing"
)

// TestInitDisabled 验证未配置 Mongo/未启用 Redis 时，Init 返回全部为 nil 的 Infra，不报错。
func TestInitDisabled(t *testing.T) {
	// 测试 CWD 无 application.yaml，config 空态；GetValueString/GetValueBoolDefault 返回默认值。
	inf, err := Init()
	if err != nil {
		t.Fatalf("Init() error = %v, want nil", err)
	}
	if inf == nil {
		t.Fatal("Init() 返回 nil Infra")
	}
	if inf.MongoClient != nil || inf.MongoDB != nil {
		t.Errorf("未配置 Mongo uri，应跳过 Mongo 连接：client=%v db=%v", inf.MongoClient, inf.MongoDB)
	}
	if inf.Redis != nil {
		t.Errorf("未启用 redis，应跳过 Redis 连接：redis=%v", inf.Redis)
	}
	_ = inf.Close()
}
