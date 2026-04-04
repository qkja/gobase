package test

import (
	"testing"

	"github.com/qkja/gobase/goid"
)

func TestUUID(t *testing.T) {
	id := goid.GenerateUUID()
	t.Logf("UUID: %s\n", id)
}
