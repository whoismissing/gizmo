package debug

import (
	"testing"
)

func TestDebugLogStatusFalse(t *testing.T) {
	t.Logf("TestDebugLogStatusFalse")

    LogBegin()
    LogEnd()

}

func TestDebugLogStatusTrue(t *testing.T) {
	t.Logf("TestDebugLogStatusTrue")

    Status = true
    LogBegin()
    LogEnd()

}
