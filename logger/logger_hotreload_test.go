package logger

import (
	"testing"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/listener"
	"github.com/sirupsen/logrus"
)

// TestHotReloadAppliesLevel 复刻泛型 Config[T] 热加载路径（config/watch.go）：
//  1. AppendValue 静默更新配置 store（不发逐 key 事件）；
//  2. watcher 补发 appconfig.reload 合成事件。
//
// 回归验证：配置热加载后 rootLogger 级别随之更新。
// 修复前 ConfigChangeListener 只认 base.logger.level 等逐 key 事件，appconfig.reload 不处理 → 级别不更新。
func TestHotReloadAppliesLevel(t *testing.T) {
	// 修复后空 store 也能直接 SetValue 写入（config.SetValue 不再静默失败）
	config.SetValue("base.logger.level", "debug")
	InitLog()
	if got := rootLogger.GetLevel(); got != logrus.DebugLevel {
		t.Fatalf("baseline: rootLogger level = %v (store=%q), want debug",
			got, config.GetValueStringDefault("base.logger.level", "<none>"))
	}

	// 模拟热加载：store 更新为 error（静默，无逐 key 事件），再补发合成事件
	config.AppendValue("base.logger.level=error")
	listener.PublishEvent(listener.ConfigChangeEvent{Key: "appconfig.reload", Value: "1"})

	if got := rootLogger.GetLevel(); got != logrus.ErrorLevel {
		t.Fatalf("after hot reload: rootLogger level = %v (store=%q), want error",
			got, config.GetValueStringDefault("base.logger.level", "<none>"))
	}
}
