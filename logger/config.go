package logger

import (
	"strings"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/listener"
)

// InitLog 应用配置到 root logger，并注册配置变更监听（支持级别热更新）。
func InitLog() {
	rootLogger.SetLevel(parseLevel(config.GetValueStringDefault("logger.level", "0")))
	gColor = config.GetValueBoolDefault("logger.color.enable", false)
	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)
}

// ConfigChangeListener 配置变更监听：逐 key 事件 + 整文件重载合成事件。
func ConfigChangeListener(event listener.BaseEvent) {
	ev := event.(listener.ConfigChangeEvent)
	switch {
	case ev.Key == "logger.level":
		SetGlobalLevel(ev.Value)
	case ev.Key == "logger.debug_tenant_ids":
		reloadDebugTenants()
	case strings.HasPrefix(ev.Key, "logger.group"):
		words := strings.Split(ev.Key, ".")
		if len(words) != 5 {
			return
		}
		Group(words[3]).SetLevel(parseLevel(ev.Value))
	case ev.Key == "appconfig.reload":
		// 泛型 Config[T] 热加载：整文件重载时只补发合成事件，逐 key 不产生。
		// 这里重读配置并应用到已创建的 root/group logger。
		applyLevelFromConfig()
	}
}

// applyLevelFromConfig 热加载后重读 logger.* 配置，应用到已创建的 logger。
// 只更新已存在的分组，不创建新分组（新分组首次 Group() 时自行读配置）。
func applyLevelFromConfig() {
	rootLevel := config.GetValueStringDefault("logger.level", "0")
	rootLogger.SetLevel(parseLevel(rootLevel))
	reloadDebugTenants()
	for groupName := range loggerMap {
		if groupName == "root" {
			continue
		}
		lvl := config.GetValueString("logger.group." + groupName + ".level")
		if lvl == "" {
			lvl = rootLevel
		}
		loggerMap[groupName].SetLevel(parseLevel(lvl))
	}
}

// SetGlobalLevel 设置 root logger 级别（热更新入口）。
func SetGlobalLevel(strLevel string) {
	rootLogger.SetLevel(parseLevel(strLevel))
}
