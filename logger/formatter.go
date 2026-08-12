package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/qkja/gobase/config"
	"github.com/sirupsen/logrus"
)

// ANSI 颜色码。
const (
	black  = 30
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
	gray   = 37
)

// StandardFormatter 统一日志格式：
//
//	[2006-01-02 15:04:05] [app-name] LEVEL path:line#func message fields
//
// 颜色输出受 logger.color.enable 控制（默认关闭）。
type StandardFormatter struct{}

func (m *StandardFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := entry.Buffer
	if b == nil {
		b = &bytes.Buffer{}
	}

	// 额外字段（k=v 列表）
	var fields []string
	for k, v := range entry.Data {
		fields = append(fields, fmt.Sprintf("%v=%v", k, v))
	}
	var fieldsStr string
	if len(fields) > 0 {
		fieldsStr = fmt.Sprintf("[\x1b[%dm%s\x1b[0m]", blue, strings.Join(fields, " "))
	}

	// 调用位置
	funPath := entry.Message
	if entry.HasCaller() {
		frame := getCallerFrame()
		funPath = fmt.Sprintf("%s:%d#%s", shortLogPath(frame.File), frame.Line, functionName(frame))
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	appName := config.GetValueStringDefault("application.name", "gobase")
	level := strings.ToUpper(entry.Level.String())

	var line string
	if gColor {
		line = fmt.Sprintf("[%s] [%s] \x1b[%dm%s\x1b[0m \x1b[%dm%s\x1b[0m %s %s\n",
			timestamp, appName,
			levelColor(entry.Level), level,
			black, funPath,
			entry.Message, fieldsStr)
	} else {
		line = fmt.Sprintf("[%s] [%s] %s %s %s %s\n",
			timestamp, appName, level, funPath, entry.Message, fieldsStr)
	}

	b.WriteString(line)
	return b.Bytes(), nil
}

// levelColor 返回级别对应的 ANSI 颜色码。
func levelColor(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return blue
	case logrus.InfoLevel:
		return green
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	default:
		return gray
	}
}

const (
	maximumCallerDepth    = 25
	knownBaseLoggerFrames = 5
)

var (
	callerInitOnce     sync.Once
	minimumCallerDepth = 0
	baseLoggerPackage  string
)

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

// getCallerFrame 定位调用日志的真实代码位置（跳过 logger 与 logrus 栈帧）。
func getCallerFrame() *runtime.Frame {
	pcs := make([]uintptr, maximumCallerDepth)
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "logger.getCallerFrame") {
				baseLoggerPackage = getPackageName(funcName)
				break
			}
		}
		minimumCallerDepth = knownBaseLoggerFrames
	})

	pcs = make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if pkg != baseLoggerPackage && pkg != "github.com/sirupsen/logrus" {
			return &f
		}
	}
	return nil
}

func functionName(frame *runtime.Frame) string {
	pathMeta := strings.Split(frame.Function, ".")
	if len(pathMeta) > 1 {
		return pathMeta[len(pathMeta)-1]
	}
	return frame.Function
}

// shortLogPath 按 logger.path.type 缩短日志路径（short 默认取后三段）。
func shortLogPath(logPath string) string {
	loggerPath := config.GetValueStringDefault("logger.path.type", "short")
	switch loggerPath {
	case "short":
		pathMeta := strings.Split(logPath, string(os.PathSeparator))
		if len(pathMeta) > 3 {
			return strings.Join(pathMeta[len(pathMeta)-3:], string(os.PathSeparator))
		}
		return logPath
	case "full":
		pathMeta := strings.Split(logPath, "@2/project")
		if len(pathMeta) > 1 {
			pathMeta[0] = "../.."
			return strings.Join(pathMeta, "")
		}
		return logPath
	default:
		return logPath
	}
}
