package internal



type Request struct {
	From int
	To   int
}

type disAEle struct {
	dis int
	ele *elevator
}

type floorAStop struct {
	floor   int
	isAStop bool
}

func Abs[T int | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
