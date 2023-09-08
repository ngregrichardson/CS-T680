package utils

// Taken from https://stackoverflow.com/a/74149031
func ShiftEnd[T any](s []T, x int) []T {
	if x < 0 {
		return s
	}
	if x >= len(s)-1 {
		return s
	}
	tmp := s[x]

	s = append(s[:x], s[x+1:]...)

	s = append(s, tmp)
	return s
}
