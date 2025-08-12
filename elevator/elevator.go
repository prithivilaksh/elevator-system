package elevator

type Elevator interface {
	GetID() int
	TryAddStops(from, to, expMetric, threshold int) bool
	GetMetric(from, to int) int
}