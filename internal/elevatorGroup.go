package internal

import (
	"sync"
)

type request struct {
	from int
	to   int
}

type elevatorGroup struct {
	elevators    []elevator
	totalFloors  int
	requestsChnl chan request
}

type disAndElevator struct {
	dis int
	ele elevator
}

func newElevatorGroup(totalFloors int) *elevatorGroup {
	eg := &elevatorGroup{
		elevators:    make([]elevator, totalFloors),
		totalFloors:  totalFloors,
		requestsChnl: make(chan request, totalFloors),
	}
	return eg
}

func (eg *elevatorGroup) addElevator(e elevator) {
	eg.elevators = append(eg.elevators, e)
}

func (eg *elevatorGroup) serve(req request) {
	eg.requestsChnl <- req
}


/*

curr from to dest
curr to from dest
curr from dest to
curr to dest from

*/

func (eg *elevatorGroup) findDistance(e elevator, req request, res chan<- disAndElevator) {
	var from, to = req.from, req.to
	var curr, dir, dest = e.currentFloor, e.direction, e.destination
	dist := Abs(curr - dest)

	if curr <= from && from <= dest {

	}

}

func (eg *elevatorGroup) SelectBestAndAssign(req request) {

	res := make(chan disAndElevator, eg.totalFloors)
	var wg sync.WaitGroup
	for _, e := range eg.elevators {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eg.findDistance(e, req, res)
		}()
	}

	go func() {
		defer close(res)
		wg.Wait()
	}()

	var minDis int = eg.totalFloors * 10
	var bestElevator elevator = nil

	for de := range res {
		dis, ele := de.dis, de.ele
		if dis < minDis {
			minDis, bestElevator = dis, ele
		}
	}

	bestElevator.serve(req)

}

func (eg *elevatorGroup) start() {
	for req := range eg.requestsChnl {
		go eg.SelectBestAndAssign(req)
	}
}
