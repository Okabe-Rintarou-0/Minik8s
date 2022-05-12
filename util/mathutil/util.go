package mathutil

func Clamp(x, min, max int) int {
	if x >= min && x <= max {
		return x
	} else if x > max {
		return max
	} else {
		return min
	}
}
