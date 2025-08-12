package elevatorgroup

import (
	"fmt"
	"slices"
	"sync"

	"github.com/prithivilaksh/elevator-system/elevator"
)

type MetricElevatorGroup struct {
	elevators []elevator.Elevator
}

func NewMetricElevatorGroup() *MetricElevatorGroup {
	return &MetricElevatorGroup{
		elevators: make([]elevator.Elevator, 0),
	}
}

func (eg *MetricElevatorGroup) AddElevator(ele elevator.Elevator) {
	eg.elevators = append(eg.elevators, ele)
}

func (eg *MetricElevatorGroup) GetElevatorID(from int, to int) (int, error) {
	type metricElevator struct {
		metric int
		ele    elevator.Elevator
	}
	metricElevators := make([]metricElevator, len(eg.elevators))
	wg := sync.WaitGroup{}
	for i, ele := range eg.elevators {
		wg.Add(1)
		go func(i int, ele elevator.Elevator) {
			defer wg.Done()
			metricElevators[i] = metricElevator{ele.GetMetric(from, to), ele}
		}(i, ele)
	}
	wg.Wait()
	slices.SortFunc(metricElevators, func(a, b metricElevator) int {
		return a.metric - b.metric
	})

	for _, metricElevator := range metricElevators {
		expMetric := metricElevator.metric
		ele := metricElevator.ele
		if ele.TryAddStops(from, to, expMetric, 4) {
			return ele.GetID(), nil
		}
	}

	return -1, fmt.Errorf("no elevator available to serve from %d to %d", from, to)
}
