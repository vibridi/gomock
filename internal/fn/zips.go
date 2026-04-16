package fn

// Zips two string slices together. Each new entry is separated by sep. sep can be an empty string.
// If the slices aren't of equal length, the shortest one is padded with empty strings.
func Zips(a1, a2 []string, sep string) []string {
	l1 := len(a1)
	l2 := len(a2)
	lmax := 0

	if l1 <= l2 {
		for i := 0; i < l2-l1; i++ {
			a1 = append(a1, "")
		}
		lmax = l2

	} else {
		for i := 0; i < l1-l2; i++ {
			a2 = append(a2, "")
		}
		lmax = l1
	}

	z := make([]string, 0, lmax)
	for i := 0; i < lmax; i++ {
		z = append(z, a1[i]+sep+a2[i])
	}
	return z
}
