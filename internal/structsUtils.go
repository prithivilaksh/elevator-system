package internal

type elevatorGroup struct {
	elevators    []*elevator
	totFloors  int
	reqsChnl chan request
}

type request struct {
	from int
	to   int
}

type disAEle struct {
	dis int
	ele *elevator
}

type floorAStop struct {
	floor int
	isAStop bool
}

func Abs[T int | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
