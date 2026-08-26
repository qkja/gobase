package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDiscovery(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, discoveryFileName), []byte(content), 0o644); err != nil {
		t.Fatalf("写 discovery.json 失败: %v", err)
	}
}

func TestGetAddressBeforeInit(t *testing.T) {
	addr, ok := GetAddress("identityhubsvr")
	if ok {
		t.Fatalf("Init 前 GetAddress 不应命中: addr=%q", addr)
	}
	if addr != "" {
		t.Fatalf("Init 前 GetAddress 应返回空串，实际 %q", addr)
	}
}

func TestInitNoLocalUsesDefault(t *testing.T) {
	dir := t.TempDir() // 目录下没有 discovery.json
	if err := Init(dir, WithNamespace("qa")); err != nil {
		t.Fatalf("Init 应成功: %v", err)
	}
	addr, ok := GetAddress("identityhubsvr")
	if !ok {
		t.Fatal("无本地文件时应命中 gobase 默认")
	}
	want := "identityhubsvr.qa.svc.cluster.local:9090"
	if addr != want {
		t.Fatalf("默认地址不符: %q, want %q", addr, want)
	}
	// 默认里没有的服务仍不可命中
	if addr, ok := GetAddress("nobody"); ok {
		t.Fatalf("默认中没有的不应命中: addr=%q", addr)
	}
}

func TestInitEnvNamespace(t *testing.T) {
	t.Setenv("POD_NAMESPACE", "prod")
	dir := t.TempDir()
	if err := Init(dir); err != nil {
		t.Fatalf("Init 应成功: %v", err)
	}
	addr, ok := GetAddress("identityhubsvr")
	if !ok {
		t.Fatal("应命中 gobase 默认")
	}
	if addr != "identityhubsvr.prod.svc.cluster.local:9090" {
		t.Fatalf("应使用 POD_NAMESPACE 命名空间: %q", addr)
	}
}

func TestInitLocalOverridesDefault(t *testing.T) {
	dir := t.TempDir()
	// 本地只配 foosvr；identityhubsvr 未配，应回退 gobase 默认
	writeDiscovery(t, dir, `{"foosvr":"10.0.0.1:9090"}`)
	if err := Init(dir, WithNamespace("qa")); err != nil {
		t.Fatalf("Init 应成功: %v", err)
	}

	addr, ok := GetAddress("foosvr")
	if !ok || addr != "10.0.0.1:9090" {
		t.Fatalf("本地配置应优先: addr=%q ok=%v", addr, ok)
	}

	addr, ok = GetAddress("identityhubsvr")
	if !ok || addr != "identityhubsvr.qa.svc.cluster.local:9090" {
		t.Fatalf("本地未配的服务应回退默认: addr=%q ok=%v", addr, ok)
	}

	if addr, ok := GetAddress("nobody"); ok {
		t.Fatalf("两级都没有的不应命中: addr=%q", addr)
	}
}

func TestInitBadJSON(t *testing.T) {
	dir := t.TempDir()
	writeDiscovery(t, dir, `{invalid`)
	if err := Init(dir); err == nil {
		t.Fatal("坏 JSON 时 Init 应返回错误")
	}
}
