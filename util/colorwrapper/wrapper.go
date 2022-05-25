package colorwrapper

var (
	reset = string([]byte{27, 91, 48, 109})
	red   = string([]byte{27, 91, 57, 49, 109})
	green = string([]byte{27, 91, 51, 50, 109})
)

func Green(target string) string {
	return green + target + reset
}
