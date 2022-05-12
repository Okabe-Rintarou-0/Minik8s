package logger

import (
	"io"
	"testing"
)

func TestWarn(t *testing.T) {
	Warn(io.EOF.Error())
}
