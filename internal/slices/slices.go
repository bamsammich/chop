package slices

func contains[T comparable](sl []T, e T) bool {
	for _, i := range sl {
		if i == e {
			return true
		}
	}
	return false
}
