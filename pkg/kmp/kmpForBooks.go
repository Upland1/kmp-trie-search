package kmp

func lps(p string) []int {
	m := len(p)

	v := make([]int, m)

	j, i := 0, 1

	for i < m {
		if p[i] == p[j] {
			v[i] = j + 1
			j++
			i++
		} else {
			if j == 0 {
				v[i] = 0
				i++
			} else {
				j = v[j-1]
			}
		}
	}

	return v
}

func Kmp(p string, t string) []int {
	positions := []int{}

	n := len(t)
	m := len(p)

	if m == 0 || n == 0 {
		return positions
	}

	v := lps(p)

	i, j := 0, 0

	for i < n {
		if t[i] == p[j] {
			i++
			j++

			if j == m {
				positions = append(positions, i-j)
				j = v[j-1]
			}
		} else {
			if j != 0 {
				j = v[j-1]
			} else {
				i++
			}
		}
	}

	return positions
}
