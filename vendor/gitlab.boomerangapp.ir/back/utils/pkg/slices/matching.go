package slices

//SliceMatchSlice  check all needles exists in heyStack
func SliceMatchSlice(heyStack []string, needles []string) bool {
	for i := range heyStack {
		for j := range needles {
			if needles[j] == heyStack[i] {
				needles = RemoveStrItem(needles, j)
				break
			}
		}
	}
	return len(needles) == 0
}

//StringContains  check heyStack contains needle and return needle`s positions
func StringContains(heyStack []string, needle string) int {
	for i, item := range heyStack {
		if item == needle {
			return i
		}
	}
	return -1
}
