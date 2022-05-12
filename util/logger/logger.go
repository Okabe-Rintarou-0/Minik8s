package logger

import "fmt"

var (
	reset = string([]byte{27, 91, 48, 109})
	red   = string([]byte{27, 91, 57, 49, 109})
)

// Log returns a closure
func Log(prefix string) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		content := fmt.Sprintf(format, v...)
		fmt.Printf("[%s] %s\n", prefix, content)
	}
}

func Warn(format string, v ...interface{}) {
	fmt.Println(red + "[Warn] " + fmt.Sprintf(format, v...) + reset)
}

func Error(format string, v ...interface{}) {
	fmt.Println(red + "[Error] " + fmt.Sprintf(format, v...) + reset)
}
