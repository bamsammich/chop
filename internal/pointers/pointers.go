package pointers

func To[T any](v T) *T {
	return &v
}
