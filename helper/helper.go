package helper

import "strings"

func Zips(a1, a2 []string, s string) []string {
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
		z = append(z, a1[i]+s+a2[i])
	}
	return z
}

func ReturnTypesToString(r []string) string {
	switch len(r) {
	case 0:
		return ""
	case 1:
		return r[0]
	default:
		return "(" + strings.Join(r, ", ") + ")"
	}
}
