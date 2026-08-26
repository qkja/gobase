// Package discovery 提供基于 discovery.json 的服务发现：按服务名解析 gRPC/HTTP 后端地址。
//
// 两层解析（服务自己的文件优先）：
//  1. 服务自己的 <cfgDir>/discovery.json（k8s 中随 ConfigMap 挂载）——高优先级；
//  2. gobase 内置默认 discovery.default.json（go:embed）——兜底。默认地址中的 {namespace}
//     占位符在 Init 时替换为实际命名空间（POD_NAMESPACE 环境变量，或 WithNamespace 选项，兜底 default）。
//
// 两级都未配置的服务返回 ("", false)，由调用方决定如何处理（如直接报错）。
package discovery

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/qkja/gobase/logger"
)

const (
	// discoveryFileName 服务自己的服务发现文件（与 application.toml 同目录）。
	discoveryFileName = "discovery.json"
	// defaultDir 空配置目录时的默认值，与 config 包内部默认目录保持一致。
	defaultDir = "./config"
	// nsPlaceholder 默认文件地址中的命名空间占位符，Init 时替换为实际命名空间。
	nsPlaceholder = "{namespace}"
	// defaultNamespace 未配置命名空间时的兜底值。
	defaultNamespace = "default"
)

//go:embed discovery.default.json
var defaultJSON []byte

// localStore 服务自己的 discovery.json 映射（高优先级）。
var localStore atomic.Value

// defaultStore gobase 内置默认映射（兜底）。
var defaultStore atomic.Value

// Option 定义 Init 的可选参数。
type Option func(*options)

type options struct {
	namespace string
}

// WithNamespace 指定默认地址使用的 k8s 命名空间（默认取 POD_NAMESPACE 环境变量，再兜底 default）。
func WithNamespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

// Init 加载 gobase 内置默认（替换命名空间占位符）+ 服务自己的 cfgDir/discovery.json（可选，覆盖默认）。
// 文件缺失不报错（本地空 map，全部走默认）；文件存在但 JSON 损坏则返回错误，启动即暴露问题。
func Init(cfgDir string, opts ...Option) error {
	o := &options{}
	for _, fn := range opts {
		fn(o)
	}
	ns := o.namespace
	if ns == "" {
		ns = os.Getenv("POD_NAMESPACE")
	}
	if ns == "" {
		ns = defaultNamespace
	}

	// 1、gobase 内置默认（兜底）
	dflt := make(map[string]string)
	if err := json.Unmarshal(defaultJSON, &dflt); err != nil {
		return errors.New("discovery: failed to parse built-in default: " + err.Error())
	}
	for k, addr := range dflt {
		dflt[k] = strings.ReplaceAll(addr, nsPlaceholder, ns)
	}
	defaultStore.Store(dflt)

	// 2、服务自己的文件（高优先级，可选）
	if cfgDir == "" {
		cfgDir = defaultDir
	}
	path := filepath.Join(cfgDir, discoveryFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			localStore.Store(map[string]string{})
			logger.Info("discovery: %s not found, use gobase built-in default", path)
			return nil
		}
		return err
	}

	local := make(map[string]string)
	if err := json.Unmarshal(data, &local); err != nil {
		return errors.New("discovery: failed to parse " + path + ": " + err.Error())
	}
	localStore.Store(local)
	logger.Info("discovery: loaded %d local service addresses, unconfigured use gobase default", len(local))
	return nil
}

// GetAddress 优先返回服务自己 discovery.json 中的地址，未命中回退 gobase 内置默认。
// 两级都未配置返回 ("", false)。Init 前调用恒返回 ("", false)。
func GetAddress(svcName string) (string, bool) {
	if v := localStore.Load(); v != nil {
		if m, ok := v.(map[string]string); ok {
			if addr, found := m[svcName]; found && addr != "" {
				return addr, true
			}
		}
	}
	if v := defaultStore.Load(); v != nil {
		if m, ok := v.(map[string]string); ok {
			if addr, found := m[svcName]; found && addr != "" {
				return addr, true
			}
		}
	}
	return "", false
}
