// Package logger 提供基于 logrus 的统一日志封装：
// 分组、文件滚动、级别热更新、彩色输出、调用位置、ctx trace 前缀。
//
// 基本用法：
//
//	logger.Info("hello %s", "world")
//	logger.ErrorCtx(ctx, "call failed: %v", err)
//	logger.Group("orm").Debugf("[SQL] %s", sql)
package logger

import (
	"strings"

	"github.com/qkja/gobase/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 包级全局状态。
var (
	loggerMap  map[string]*logrus.Logger
	rootLogger *logrus.Logger
	gColor     bool
)

func init() {
	loggerMap = map[string]*logrus.Logger{}
	rootLogger = getOrCreate("root")
	// 租户级 debug 专用 logger：恒为 DebugLevel，由 DebugTenant 的租户检查控制输出
	tenantDebugLogger = newLogger("tenant_debug")
	tenantDebugLogger.SetLevel(logrus.DebugLevel)
	gColor = config.GetValueBoolDefault("logger.color.enable", false)
}

// getOrCreate 统一获取或创建分组 logger（合并了旧的 Group/doGroup 重复逻辑）。
func getOrCreate(groupName string) *logrus.Logger {
	if l, ok := loggerMap[groupName]; ok {
		return l
	}
	l := newLogger(groupName)
	loggerMap[groupName] = l
	return l
}

// newLogger 创建带 formatter + 文件滚动 hook 的 logger，并读取分组级别。
func newLogger(groupName string) *logrus.Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &StandardFormatter{}

	loggerDir := config.GetValueStringDefault("logger.home", "./logs/")
	l.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotateLogWithCache(loggerDir, "debug"),
		logrus.InfoLevel:  rotateLogWithCache(loggerDir, "info"),
		logrus.WarnLevel:  rotateLogWithCache(loggerDir, "warn"),
		logrus.ErrorLevel: rotateLogWithCache(loggerDir, "error"),
		logrus.PanicLevel: rotateLogWithCache(loggerDir, "panic"),
		logrus.FatalLevel: rotateLogWithCache(loggerDir, "fatal"),
	}, l.Formatter))

	l.SetLevel(levelFor(groupName))
	return l
}

// levelFor 读取分组日志级别，未配置回退全局 logger.level。
func levelFor(groupName string) logrus.Level {
	rootLevel := config.GetValueStringDefault("logger.level", "0")
	lvl := config.GetValueString("logger.group." + groupName + ".level")
	if lvl == "" {
		lvl = rootLevel
	}
	return parseLevel(lvl)
}

// parseLevel 解析日志级别，支持数字和字符串两种形式。
// 数字映射：-2=trace -1=debug 0=info 1=warn 2=error 3=fatal 4=panic。
// 字符串兼容旧的 "debug"/"info" 等写法。
func parseLevel(v string) logrus.Level {
	switch v {
	case "-2", "trace":
		return logrus.TraceLevel
	case "-1", "debug":
		return logrus.DebugLevel
	case "0", "info":
		return logrus.InfoLevel
	case "1", "warn", "warning":
		return logrus.WarnLevel
	case "2", "error":
		return logrus.ErrorLevel
	case "3", "fatal":
		return logrus.FatalLevel
	case "4", "panic":
		return logrus.PanicLevel
	default:
		if le, err := logrus.ParseLevel(v); err == nil {
			return le
		}
		return logrus.InfoLevel
	}
}

// Group 返回指定分组 logger。支持传入多个别名：返回第一个已存在的分组；
// 若都未创建，则创建第一个名字的分组。
func Group(groupNames ...string) *logrus.Logger {
	if len(groupNames) == 0 {
		return rootLogger
	}
	for _, name := range groupNames {
		if l, ok := loggerMap[name]; ok {
			return l
		}
	}
	return getOrCreate(groupNames[0])
}

// ---- 包级便捷函数（root logger）----

func Info(format string, v ...any)  { rootLogger.Infof(format, v...) }
func Warn(format string, v ...any)  { rootLogger.Warnf(format, v...) }
func Error(format string, v ...any) { rootLogger.Errorf(format, v...) }
func Debug(format string, v ...any) { rootLogger.Debugf(format, v...) }
func Fatal(format string, v ...any) { rootLogger.Fatalf(format, v...) }

// ---- 不定参版本（供需要直接传 args 的调用方）----

func InfoDirect(v ...any)  { rootLogger.Info(v...) }
func WarnDirect(v ...any)  { rootLogger.Warn(v...) }
func ErrorDirect(v ...any) { rootLogger.Error(v...) }
func DebugDirect(v ...any) { rootLogger.Debug(v...) }

// Record 按字符串级别动态输出（默认 debug）。
func Record(level, format string, v ...any) {
	switch strings.ToLower(level) {
	case "info":
		Info(format, v...)
	case "warn":
		Warn(format, v...)
	case "error":
		Error(format, v...)
	case "fatal":
		Fatal(format, v...)
	default:
		Debug(format, v...)
	}
}
