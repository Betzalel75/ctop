package utils

// Capitalize the first letter of a string
func Capitalize(s string) string {
	n := []rune(s)
	first := true

	for i := range n {
		if isValid(n[i]) && first {
			if n[i] >= 'a' && n[i] <= 'z' {
				n[i] = n[i] - 32
			}
			first = false
		} else if n[i] >= 'A' && n[i] <= 'Z' {
			n[i] = n[i] + 32
		} else if !isValid(n[i]) {
			first = true
		}
	}

	return string(n)
}

func isValid(f rune) bool {
	if (f >= 'a' && f <= 'z') || (f >= 'A' && f <= 'Z') || (f >= '0' && f <= '9') {
		return true
	}
	return false
}
