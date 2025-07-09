package internal

func Abs[T int | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
