//go:build darwin && !cgo

package host

import (
	"context"

	"github.com/qkja/gobase/system/common"
)

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}
