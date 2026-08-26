package i18n

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"sync"

	f0 "github.com/qkja/gobase/file"
	"github.com/qkja/gobase/isc"
)

// 语言标记，全仓库统一使用，其他服务可直接引用。
const (
	// LangZh 简体中文
	LangZh = "zh-CN"
	// LangEn 英文
	LangEn = "en-US"
)

//go:embed default/*.po
var embeddedPOs embed.FS

// 默认语言：Lookup 在当前语言缺失 key 时的回退语言。
// InitI18N 会通过 innerMap.DefaultLanguage 记录初始化语言，未初始化时回退中文。
var defaultLanguage = "zh-CN"

// langCache 按语言缓存解析结果（嵌入默认 + 服务文件叠加覆盖），并发安全。
// 服务文件仅首次探测缓存，之后纯内存查找，热路径无文件 IO。
var langCache sync.Map // lang -> map[string]string

type I18NMap struct {
	Language        string            // 当前的语言
	DefaultLanguage string            // 默认语言
	Data            map[string]string // 当前的字符串映射
	DefaultData     map[string]string // 默认的字符串映射，当从Data内找不到key时，从此处找，再找不到就报错
}

var innerMap *I18NMap

func NewI18NMap(language string, filePath string) *I18NMap {
	m := loadPo(filePath)
	return &I18NMap{
		Language:        language,
		DefaultLanguage: language,
		Data:            m,
		DefaultData:     m,
	}
}

// Lookup 按语言查找 key 的翻译（探测语义，缺失不打日志）。
// 查找链：<lang> 语言（嵌入默认 + 服务 i18n/<lang>.po 叠加）→ 默认语言 → 未命中。
func Lookup(lang, key string) (string, bool) {
	if v, ok := langMap(lang)[key]; ok && v != "" {
		return v, true
	}
	if dl := defaultLang(); dl != "" && dl != lang {
		if v, ok := langMap(dl)[key]; ok && v != "" {
			return v, true
		}
	}
	return "", false
}

// defaultLang 返回当前默认语言：优先 InitI18N 设置的语言，未初始化回退默认中文。
func defaultLang() string {
	if innerMap != nil && innerMap.DefaultLanguage != "" {
		return innerMap.DefaultLanguage
	}
	return defaultLanguage
}

// langMap 返回指定语言合并后的 map：嵌入默认 为底表，服务 <CWD>/i18n/<lang>.po 叠加覆盖（服务优先）。
func langMap(lang string) map[string]string {
	if v, ok := langCache.Load(lang); ok {
		return v.(map[string]string)
	}
	m := map[string]string{}
	if b, err := embeddedPOs.ReadFile("default/" + lang + ".po"); err == nil {
		for k, v := range parsePo(string(b)) {
			m[k] = v
		}
	}
	if p := poFilePath(lang); f0.FileExists(p) {
		for k, v := range parsePoLines(f0.ReadFileLines(p)) {
			m[k] = v
		}
	}
	langCache.Store(lang, m)
	return m
}

// poFilePath 服务语言文件路径：<CWD>/i18n/<lang>.po
func poFilePath(lang string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "i18n", lang+".po")
}

// loadPo 读取语言文件（兼容既有 InitI18N / LoadI18NLanguage 的全局语义）
func loadPo(filePath string) map[string]string {
	return parsePoLines(f0.ReadFileLines(filePath))
}

// parsePo 解析 po 内容字符串
func parsePo(content string) map[string]string {
	return parsePoLines(strings.Split(content, "\n"))
}

// parsePoLines 解析 po 行序列：key 为首个空格前内容，value 为之后内容。
// po 文件格式（每行一条）：key value，例如：
//
//	msgid "hello"
//	1001 系统内部错误
func parsePoLines(lines []string) map[string]string {
	m := make(map[string]string)
	for _, s := range lines {
		ss := isc.ISCString(s)
		key := ss.SubStringBefore(" ")
		if string(key) == "" {
			continue
		}
		val := ss.SubStringAfter(" ").TrimSpace().Trim("\"").ReplaceAll("\\n", "\n").ReplaceAll("\\r", "\r").ReplaceAll("\\t", "\t").ReplaceAll("\\\"", "\"").ReplaceAll("\\\\", "\\")
		m[string(key)] = string(val)
	}
	return m
}
