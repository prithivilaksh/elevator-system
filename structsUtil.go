package main

type fromTo struct {
	from int
	to   int
}

type fromToIsAdded struct {
	fromTo
	dis int
	isAdded chan bool
}

type disAEle struct {
	dis int
	ele *elevator
}

type stopsAndCurrFloor struct {
	stops   []int
	currFloor int
}

func Abs[T int | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
