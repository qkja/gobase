package goid

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	// ulidMu 串行化 ULID 生成：同一毫秒内熵单调递增，保证 code 严格有序
	ulidMu sync.Mutex
	// monoEnt oklog Monotonic 熵源：同毫秒 +1 递增，跨毫秒重新随机
	monoEnt = ulid.Monotonic(rand.Reader, 0)
	// lastULIDMs 进程内最近一次毫秒时间戳：时钟回拨时复用，避免生成"更早"的 code
	lastULIDMs uint64
)

// GenerateCode 生成带前缀、按时间有序的业务 code：prefix + "_" + 26 位标准 ULID。
//
// ULID 前 48 bit 为毫秒时间戳（时间戳作为基数）→ 字典序即生成时间序；
// 后 80 bit 熵保证多实例间不冲突；同一进程内同毫秒通过熵单调递增保证严格有序。
// 相比随机 UUID（GenerateUUID），code 天然可按创建时间排序，便于日志与 DB 索引定位。
//
// 用法：code := goid.GenerateCode("AUD") → "AUD_01HV..."；前缀由调用方传入，
// 前缀与生成值之间自动加下划线分隔（prefix 为空串时仅返回纯 ULID，无下划线）。
func GenerateCode(prefix string) string {
	ulidMu.Lock()
	defer ulidMu.Unlock()

	ms := ulid.Timestamp(time.Now())
	if ms < lastULIDMs {
		ms = lastULIDMs // 时钟回拨保护：不生成早于上一次的 code
	}
	id, err := ulid.New(ms, monoEnt)
	if err != nil {
		// 同毫秒熵递增溢出等极端情况（约 2^80 分之一）：顺延 1ms 重新生成
		ms++
		id, _ = ulid.New(ms, monoEnt)
	}
	lastULIDMs = ms

	// 前缀与生成值之间自动拼下划线；空前缀仅返回纯 ULID
	if prefix == "" {
		return id.String()
	}
	return prefix + "_" + id.String()
}
