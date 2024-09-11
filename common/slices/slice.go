package slices

func Range[t any](sl []t, handle func(i int, v t) (_break bool)) {
	for i := range sl {
		b := handle(i, sl[i])
		if b {
			break
		}
	}
}

func RangeToNew[old, new any](sl []old, handle func(i int, v old) new) []new {
	ns := make([]new, len(sl))
	for i := range ns {
		ns[i] = handle(i, sl[i])
	}
	return ns
}
