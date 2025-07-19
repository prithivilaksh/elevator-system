package elevatorGroup

import (
	"fmt"
	"slices"
	"sync"

	. "github.com/prithivilaksh/elevator-system/types"
)

type DisAEle struct {
	dis int
	ele Elevator
}

type elevatorGroup struct {
	elevators []Elevator
	totFloors int
	reqsChnl  chan FromTo
}

func NewElevatorGroup(totalFloors int, elevators []Elevator) *elevatorGroup {
	eg := &elevatorGroup{
		elevators: elevators,
		totFloors: totalFloors,
		reqsChnl:  make(chan FromTo, totalFloors*len(elevators)*20),
	}
	go eg.Start()
	return eg
}

func (eg *elevatorGroup) findDistances(req FromTo, disAEleChnl chan DisAEle) {
	defer close(disAEleChnl)
	var wg sync.WaitGroup
	for _, ele := range eg.elevators {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stopsAndCurrFloor := ele.GetStopsAndCurrFloor()
			nextStartInd:=0
			stops, nextStartInd := ele.AddStopAndGetNextInd(nextStartInd, stopsAndCurrFloor.CurrFloor, req.From, stopsAndCurrFloor.Stops)
			stops, nextStartInd = ele.AddStopAndGetNextInd(nextStartInd, req.From, req.To, stops)
			dis := ele.FindDistance(stops, stopsAndCurrFloor.CurrFloor, nextStartInd)
			disAEleChnl <- DisAEle{dis: dis, ele: ele}
		}()
	}
	wg.Wait()
}

func (eg *elevatorGroup) TryAddStops(req FromTo, disAEles *[]DisAEle) {
	isAdded := false
	for _, de := range *disAEles {
		isAdded = de.ele.TryAddStops(req, de.dis)
		if isAdded {
			break
		} 
	}
	if !isAdded {
		fmt.Println("No elevator available to serve the request from", req.From, "to", req.To, "... retrying")
		eg.Serve(req.From, req.To)
	}
}

func (eg *elevatorGroup) selectBestAndAdd(req FromTo) {

	disAEleChnl := make(chan DisAEle, len(eg.elevators))
	go eg.findDistances(req, disAEleChnl)

	var disAEles []DisAEle
	for de := range disAEleChnl {
		disAEles = append(disAEles, de)
	}

	slices.SortFunc(disAEles, func(a, b DisAEle) int { return a.dis - b.dis })
	eg.TryAddStops(req, &disAEles)
}

func (eg *elevatorGroup) Serve(from, to int) {
	eg.reqsChnl <- FromTo{From: from, To: to}
}

func (eg *elevatorGroup) Start() {
	defer close(eg.reqsChnl)
	for req := range eg.reqsChnl {
		go eg.selectBestAndAdd(req)
	}
}
