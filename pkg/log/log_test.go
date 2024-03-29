package log

import (
	"testing"
)

func TestInitUnsupportedLogLevel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for unsupported log level")
		}
	}()

	Init("INVALID")
}
