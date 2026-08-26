package test

import (
	"fmt"
	"testing"
	t0 "time"

	"github.com/qkja/gobase/debug"
	"github.com/qkja/gobase/time"
)

func TestWatcher(t *testing.T) {
	debug.Init()
	debug.AddWatcher("debug.test", func(key string, value string) {
		fmt.Println("最新的值1 key=", key, ", value=", value)
	})
	debug.AddWatcher("debug.test", func(key string, value string) {
		fmt.Println("最新的值2 key=", key, ", value=", value)
	})
	debug.StartWatch()

	t0.Sleep(1000000000000)
}

func TestPush(t *testing.T) {
	debug.Init()
	var tim = time.TimeToStringYmdHmsS(time.Now())
	fmt.Println(tim)
	debug.Update("debug.test", tim)
}
