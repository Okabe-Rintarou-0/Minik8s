package imageutil

func isUppercase(char rune) bool {
	return char >= 'A' && char <= 'Z'
}

func letterToLower(char rune) rune {
	return char - 'A' + 'a'
}

func FormatImageName(imageName string) string {
	var formatted string
	for _, char := range imageName {
		if isUppercase(char) {
			formatted += "-" + string(letterToLower(char))
		} else {
			formatted += string(char)
		}
	}
	return formatted
}
