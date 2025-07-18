package internal



type fromToIsAdded struct {
	from int
	to   int
	dis int
	isAdded chan bool
}

type disAEle struct {
	dis int
	ele *elevator
}

type floorAStop struct {
	floor   int
	isAStop bool
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
