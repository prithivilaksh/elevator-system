package main

type fromTo struct {
	from int
	to   int
}

type fromToDisIsAdded struct {
	fromTo
	dis     int
	isAddedChnl chan bool
}

type disAEle struct {
	dis int
	ele *elevator
}

type stopsACurrFloor struct {
	stops     []int
	currFloor int
}

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