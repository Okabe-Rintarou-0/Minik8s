package logger

import "fmt"

var (
	redFont = string([]byte{27, 91, 57, 49, 109})
)

// Log returns a closure
func Log(prefix string) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		content := fmt.Sprintf(format, v...)
		fmt.Printf("[%s] %s\n", prefix, content)
	}
}

func Warn(format string, v ...interface{}) {
	fmt.Println(redFont + "[Warn] " + fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	fmt.Println(redFont + "[Error] " + fmt.Sprintf(format, v...))
}
