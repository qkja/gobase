package logger

import (
	"log"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/isc"
)

// rotateMap 按 "path-level" 缓存滚动器，避免重复创建。
var (
	rotateMap map[string]*rotatelogs.RotateLogs
	rotateMu  sync.RWMutex
)

// rotateLog 创建指定 path/level 的文件滚动器。
func rotateLog(path, level string) *rotatelogs.RotateLogs {
	rotateMu.Lock()
	defer rotateMu.Unlock()

	if rotateMap == nil {
		rotateMap = map[string]*rotatelogs.RotateLogs{}
	}
	if path == "" {
		path = "./logs/"
	}

	maxSizeStr := config.GetValueStringDefault("logger.rotate.max-size", "300MB")
	maxHistoryStr := config.GetValueStringDefault("logger.rotate.max-history", "60d")
	rotateTimeStr := config.GetValueStringDefault("logger.rotate.time", "1d")

	rotateOptions := []rotatelogs.Option{rotatelogs.WithLinkName(path + "app-" + level + ".log")}
	if maxSizeStr != "" {
		rotateOptions = append(rotateOptions, rotatelogs.WithRotationSize(isc.ParseByteSize(maxSizeStr)))
	}
	if maxHistory, err := time.ParseDuration(maxHistoryStr); err == nil {
		rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(maxHistory))
	}
	if rotateTime, err := time.ParseDuration(rotateTimeStr); err == nil {
		rotateOptions = append(rotateOptions, rotatelogs.WithRotationTime(rotateTime))
	}

	data, err := rotatelogs.New(path+"app-"+level+".%Y%m%d.log", rotateOptions...)
	if err != nil {
		log.Printf("[logger] rotateLog 创建滚动器失败: %v", err)
		return nil
	}
	rotateMap[path+"-"+level] = data
	return data
}

// rotateLogWithCache 从缓存取滚动器，无则创建。
func rotateLogWithCache(path, level string) *rotatelogs.RotateLogs {
	rotateMu.RLock()
	if v, ok := rotateMap[path+"-"+level]; ok {
		rotateMu.RUnlock()
		return v
	}
	rotateMu.RUnlock()
	return rotateLog(path, level)
}
