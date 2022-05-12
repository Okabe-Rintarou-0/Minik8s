package logger

import "fmt"

// Log returns a closure
func Log(prefix string) func(format string, vs ...interface{}) {
	return func(format string, vs ...interface{}) {
		content := fmt.Sprintf(format, vs...)
		fmt.Printf("[%s] %s\n", prefix, content)
	}
}
