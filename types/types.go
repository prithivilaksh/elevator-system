package types

type FromTo struct {
	From int
	To   int
}

type StopsACurrFloor struct {
	Stops     []int
	CurrFloor int
}

type Elevator interface {
	GetStopsAndCurrFloor() StopsACurrFloor
	AddStopAndGetNextInd(startInd, afterFloor, stop int, stops []int) ([]int, int)
	FindDistance(stops []int, currFloor int, lastInd int) int
	TryAddStops(req FromTo, dis int) bool
}
