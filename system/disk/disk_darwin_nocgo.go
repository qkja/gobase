//go:build darwin && !cgo

package disk

import (
	"context"

	"github.com/qkja/gobase/system/common"
)

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}
