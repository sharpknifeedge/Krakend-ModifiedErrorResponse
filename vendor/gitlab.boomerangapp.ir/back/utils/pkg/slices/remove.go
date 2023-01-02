package slices

//RemoveStrItem remove item by key
// for []string
func RemoveStrItem(s []string, i int) []string {
	if i < 0 || i > len(s)-1 {
		return s
	}

	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

//RemoveIntItem remove item by key
// for []int
func RemoveIntItem(s []int, i int) []int {
	if i < 0 || i > len(s)-1 {
		return s
	}

	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
