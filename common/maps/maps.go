package maps

func Merge[k comparable, v any](m1 map[k]v, m2 ...map[k]v) map[k]v {
	for _, m2v := range m2 {
		for k2, v2 := range m2v {
			m1[k2] = v2
		}
	}
	return m1
}
