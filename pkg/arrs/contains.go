package arrs

func Contains[T comparable](value T, arr []T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
