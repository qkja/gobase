package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// testCfg 测试配置结构体
type testCfg struct {
	Name   string         `toml:"name"`
	Nested nestedCfg      `toml:"nested"`
	Map    map[string]int `toml:"map"`
}

type nestedCfg struct {
	QPS int    `toml:"qps"`
	TTL string `toml:"default_ttl"` // 连续大写键，验证显式 tag
}

func writeCfg(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "application.toml"), []byte(content), 0o644); err != nil {
		t.Fatalf("写配置失败: %v", err)
	}
}

const testToml = `
name = "test-svc"

[nested]
qps = 30
default_ttl = "5m"

[map]
a = 1
b = 2
`

// TestInit_FillPointer 断言 Init 赋值全局变量 + Get 返回快照
func TestInit_FillPointer(t *testing.T) {
	dir := t.TempDir()
	writeCfg(t, dir, testToml)

	var appCfg *Config[testCfg]
	if err := Init(&appCfg, dir, nil); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}
	defer appCfg.Stop()

	// 全局变量被 Init 赋值（非 nil）
	if appCfg == nil {
		t.Fatal("Init 后 appCfg 仍为 nil，应被赋值")
	}

	// Get 返回当前快照
	got := appCfg.Get()
	if got.Name != "test-svc" {
		t.Errorf("Name = %q, want test-svc", got.Name)
	}
	if got.Nested.QPS != 30 {
		t.Errorf("Nested.QPS = %d, want 30", got.Nested.QPS)
	}
	if got.Nested.TTL != "5m" {
		t.Errorf("Nested.TTL = %q, want 5m（连续大写键 tag 映射）", got.Nested.TTL)
	}
	if len(got.Map) != 2 || got.Map["a"] != 1 {
		t.Errorf("Map = %v, want {a:1 b:2}", got.Map)
	}
}

// TestInit_Defaults 断言 defaults 回调生效
func TestInit_Defaults(t *testing.T) {
	dir := t.TempDir()
	writeCfg(t, dir, "[nested]\nqps = 0\n") // qps 缺省

	var appCfg *Config[testCfg]
	if err := Init(&appCfg, dir, func(c *testCfg) {
		if c.Nested.QPS <= 0 {
			c.Nested.QPS = 10
		}
	}); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}
	defer appCfg.Stop()

	if got := appCfg.Get(); got.Nested.QPS != 10 {
		t.Errorf("defaults 后 QPS = %d, want 10", got.Nested.QPS)
	}
}

// TestWatcher_HotReload 集成测试：改文件后轮询检测并更新快照
func TestWatcher_HotReload(t *testing.T) {
	dir := t.TempDir()
	writeCfg(t, dir, testToml)

	var appCfg *Config[testCfg]
	if err := Init(&appCfg, dir, nil); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}
	defer appCfg.Stop()

	// 修改配置：qps 30 → 88
	content, _ := os.ReadFile(filepath.Join(dir, "application.toml"))
	newContent := replaceAll(string(content), "qps = 30", "qps = 88")
	os.WriteFile(filepath.Join(dir, "application.toml"), []byte(newContent), 0o644)

	// 轮询等待热加载生效
	deadline := time.Now().Add(2 * time.Second)
	for {
		if appCfg.Get().Nested.QPS == 88 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("热加载未生效: QPS 仍 = %d", appCfg.Get().Nested.QPS)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// TestWatcher_BadToml_KeepsOld 断言坏 toml 重载失败时保留旧快照
func TestWatcher_BadToml_KeepsOld(t *testing.T) {
	dir := t.TempDir()
	writeCfg(t, dir, testToml)

	var appCfg *Config[testCfg]
	if err := Init(&appCfg, dir, nil); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}
	defer appCfg.Stop()

	before := appCfg.Get()

	// 写坏 toml
	os.WriteFile(filepath.Join(dir, "application.toml"), []byte("[nested\ninvalid"), 0o644)

	// 轮询等待重载尝试（应失败保留旧快照）
	time.Sleep(300 * time.Millisecond)

	if got := appCfg.Get(); got.Nested.QPS != before.Nested.QPS {
		t.Errorf("坏 toml 后快照被改: QPS = %d, want %d", got.Nested.QPS, before.Nested.QPS)
	}
}

func replaceAll(s, old, new string) string {
	out := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			out += new
			i += len(old)
		} else {
			out += string(s[i])
			i++
		}
	}
	return out
}
