package slice

// Map applies the given function to each element in the input list and returns a new list
// containing the results. Returns an empty list if the input list is empty or if the provided
// function is nil.
func Map[T, R any](list []T, f func(t T) R) []R {
	if f == nil {
		return make([]R, 0)
	}

	output := make([]R, 0, len(list))

	for idx := range list {
		output = append(output, f(list[idx]))
	}

	return output
}

// Filter returns a new slice containing the elements of the input slice that pass the provided
// filter function. If the input slice is empty or if the filter function is nil, an empty slice
// is returned.
func Filter[T comparable](arr []T, filterFn func(v T) bool) []T {
	output := make([]T, 0, len(arr))
	for _, v := range arr {
		if filterFn == nil || filterFn(v) {
			output = append(output, v)
		}
	}

	return output
}
