package main

import (
	"fmt"
	"slices"
	"sync"
)

type elevatorGroup struct {
	elevators []*elevator
	totFloors int
	reqsChnl  chan fromTo
}

func NewElevatorGroup(totalFloors int, elevators []*elevator) *elevatorGroup {
	eg := &elevatorGroup{
		elevators: elevators,
		totFloors: totalFloors,
		reqsChnl:  make(chan fromTo, totalFloors*len(elevators)*20),
	}
	go eg.Start()
	return eg
}

func (eg *elevatorGroup) findDistances(req fromTo, disAEleChnl chan disAEle) {
	defer close(disAEleChnl)
	var wg sync.WaitGroup
	for _, ele := range eg.elevators {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stopsAndCurrFloor := ele.GetStopsAndCurrFloor()
			nextStartInd:=0
			stops, nextStartInd := ele.AddStopAndGet(nextStartInd, stopsAndCurrFloor.currFloor, req.from, stopsAndCurrFloor.stops)
			stops, nextStartInd = ele.AddStopAndGet(nextStartInd, req.from, req.to, stops)
			dis := ele.FindDistance(stops, stopsAndCurrFloor.currFloor, nextStartInd)
			disAEleChnl <- disAEle{dis: dis, ele: ele}
		}()
	}
	wg.Wait()
}

func (eg *elevatorGroup) TryAddStops(req fromTo, disAEles *[]disAEle) {
	isAdded := false
	for _, de := range *disAEles {
		isAdded = de.ele.TryAddStops(req, de.dis)
		if isAdded {
			break
		} 
	}
	if !isAdded {
		fmt.Println("No elevator available to serve the request from", req.from, "to", req.to, "... retrying")
		eg.Serve(req.from, req.to)
	}
}

func (eg *elevatorGroup) selectBestAndAdd(req fromTo) {

	disAEleChnl := make(chan disAEle, len(eg.elevators))
	go eg.findDistances(req, disAEleChnl)

	var disAEles []disAEle
	for de := range disAEleChnl {
		disAEles = append(disAEles, de)
	}

	slices.SortFunc(disAEles, func(a, b disAEle) int { return a.dis - b.dis })
	eg.TryAddStops(req, &disAEles)
}

func (eg *elevatorGroup) Serve(from, to int) {
	eg.reqsChnl <- fromTo{from: from, to: to}
}

func (eg *elevatorGroup) Start() {
	defer close(eg.reqsChnl)
	for req := range eg.reqsChnl {
		go eg.selectBestAndAdd(req)
	}
}
