package utils

func Abs[T int | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func DeepCopy[T any](src []T) []T {
	dest := make([]T, len(src))
	copy(dest, src)
	return dest
}