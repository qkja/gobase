package config

import (
	"sync/atomic"
)

// Config 泛型配置管理器：加载 + 热加载 + 原子快照。
//
// 用法：服务声明一个包级全局变量 `var AppCfg *config.Config[T]`，启动时
//
//	config.Init(&AppCfg, "./config", applyDefaults)
//	AppCfg.Get().Xxx   // 读当前快照
//
// Init 会：
//  1. 从 dir/application.toml 加载配置并启动文件轮询热加载（k8s ConfigMap 友好）；
//  2. 把管理器直接赋值给调用方传入的全局变量 out（Init 内部完成赋值，无需再手动赋值）；
//  3. 热加载后请用 AppCfg.Get() 读最新值。
type Config[T any] struct {
	current  atomic.Pointer[T]
	dir      string
	defaults func(*T)
	watcher  *watcher
}

// Init 加载配置并自动启动热加载，将结果赋值给 out 指向的全局变量。
// out 为调用方的 *gobase.Config[T] 全局变量地址（如 &AppCfg），Init 内部完成赋值；
// defaults 为业务默认值回调（nil 可省略）；dir 为空时用 "./config"。
func Init[T any](out **Config[T], dir string, defaults func(*T)) error {
	if dir == "" {
		dir = defaultDir
	}

	m := &Config[T]{dir: dir, defaults: defaults}

	// 初始加载（本 goroutine 触碰 gobase 全局）
	snapshot, err := m.snapshot()
	if err != nil {
		return err
	}
	m.current.Store(snapshot)

	// 启动热加载轮询
	m.watcher = newWatcher(dir, defaultWatchInterval, m.reload)
	m.watcher.start()

	// Init 内部直接赋值全局变量
	if out != nil {
		*out = m
	}
	return nil
}

// Get 返回当前只读快照（原子加载，Init 后永不为 nil）。
// 返回的指针为不可变快照，调用方不得修改其字段。
func (m *Config[T]) Get() *T {
	return m.current.Load()
}

// Stop 停止热加载轮询。
func (m *Config[T]) Stop() {
	if m.watcher != nil {
		m.watcher.stopWatcher()
	}
}

// snapshot 生成一个全新的配置快照：gobase 重读 → 填充 T → 应用默认值 → yaml 校验。
// 只在 Init / reload goroutine 中调用（gobase 全局 map 无锁）。
func (m *Config[T]) snapshot() (*T, error) {
	if err := validateTomlFile(m.dir); err != nil {
		return nil, err
	}

	LoadConfigFromAbsPath(m.dir)
	val := new(T)
	if err := GetValueObject("", val); err != nil {
		return nil, err
	}
	if m.defaults != nil {
		m.defaults(val)
	}
	return val, nil
}

// reload 供 watcher 调用：生成新快照并原子替换。
func (m *Config[T]) reload() error {
	snapshot, err := m.snapshot()
	if err != nil {
		return err
	}
	m.current.Store(snapshot)
	return nil
}
