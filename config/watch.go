package config

import (
	"crypto/sha256"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/qkja/gobase/listener"
)

// 默认配置目录与轮询间隔。
// 间隔 1s：k8s ConfigMap 传播近实时，1s 在延迟与开销间均衡。
const (
	defaultDir           = "./config"
	defaultWatchInterval = time.Second
	// configFileName 项目配置文件（TOML 格式）
	configFileName = "application.toml"
)

// watcher 配置文件热加载轮询器（通用，不绑定具体结构体）。
// k8s 中 ConfigMap 挂载通过 symlink 原子替换（inotify 会监听旧 inode 而失效），
// 因此采用「size+ModTime 快速跳过 + 内容 sha256 确认」的轮询方案。
type watcher struct {
	dir        string // 配置目录
	file       string // application.yaml 绝对路径
	interval   time.Duration
	stop       chan struct{}
	done       chan struct{}
	reload     func() error // 由 Config[T] 注入：生成新快照并原子替换
	mu         sync.Mutex
	lastHash   [sha256.Size]byte
	lastStat   os.FileInfo // 上次 stat 摘要，用于快速跳过
	hasBaseline bool       // 是否已记录初始哈希基线
}

// newWatcher 创建轮询器。dir 为空用默认目录；interval<=0 用默认 1s。
func newWatcher(dir string, interval time.Duration, reload func() error) *watcher {
	if dir == "" {
		dir = defaultDir
	}
	if interval <= 0 {
		interval = defaultWatchInterval
	}
	return &watcher{
		dir:      dir,
		file:     filepath.Join(dir, configFileName),
		interval: interval,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
		reload:   reload,
	}
}

// start 启动轮询 goroutine（非阻塞）。
func (w *watcher) start() {
	go w.loop()
}

// stopWatcher 通知轮询退出并等待结束。
func (w *watcher) stopWatcher() {
	select {
	case <-w.done:
		return
	default:
	}
	close(w.stop)
	<-w.done
}

// loop 轮询主循环，直到 stop 被关闭。
func (w *watcher) loop() {
	defer close(w.done)

	w.mu.Lock()
	w.initHash() // 记录基线，避免首轮误触发
	w.mu.Unlock()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.poll()
		}
	}
}

// initHash 记录当前文件哈希基线（幂等）。调用方需持有 w.mu。
func (w *watcher) initHash() {
	if w.hasBaseline {
		return
	}
	if h, err := fileHash(w.file); err == nil {
		w.lastHash = h
		w.hasBaseline = true
	}
}

// poll 检查文件是否变化，变化则重载并原子替换快照。自持锁。
func (w *watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.initHash()

	// 快速路径：size + ModTime 未变则跳过
	info, err := os.Stat(w.file)
	if err != nil {
		return // 文件暂时不可见（ConfigMap 切换间隙），保留旧快照
	}
	if w.lastStat != nil &&
		w.lastStat.Size() == info.Size() && w.lastStat.ModTime().Equal(info.ModTime()) {
		return
	}
	w.lastStat = info

	// 慢路径：读内容 + sha256 确认
	h, err := fileHash(w.file)
	if err != nil {
		log.Printf("[gobase/config] 配置热加载: 读取文件失败 %s: %v", w.file, err)
		return
	}
	if h == w.lastHash {
		return
	}

	// 重载：全部成功才替换；失败保留旧快照，绝不 panic
	if err := w.reload(); err != nil {
		log.Printf("[gobase/config] 配置热加载: 重载失败，保留旧配置: %v", err)
		return
	}
	w.lastHash = h
	log.Printf("[gobase/config] 配置热加载: %s 已更新", w.file)

	// 补发合成事件，供未来模块订阅；核心机制不依赖它
	listener.PublishEvent(listener.ConfigChangeEvent{Key: "appconfig.reload", Value: "1"})
}

// fileHash 计算文件内容 sha256。
func fileHash(path string) ([sha256.Size]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	return sha256.Sum256(data), nil
}

// validateTomlFile 独立读 application.toml 并校验 toml 语法。
// gobase 的 LoadConfigFromAbsPath 内部会吞掉坏 toml 的解析错误（只打日志），
// 因此须在此预校验，保证坏配置不会被静默替换。
func validateTomlFile(dir string) error {
	data, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		return err
	}
	var m map[string]any
	if err := toml.Unmarshal(data, &m); err != nil {
		return err
	}
	return nil
}
