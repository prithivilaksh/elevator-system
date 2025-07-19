package main

import (
	"fmt"
	"slices"
	"sync"
)

type elevatorGroup struct {
	elevators  []*elevator
	totFloors  int
	fromToChnl chan fromTo
}

func NewElevatorGroup(totalFloors int, totalElevators int) *elevatorGroup {
	eg := &elevatorGroup{
		elevators:  make([]*elevator, 0, totalElevators),
		totFloors:  totalFloors,
		fromToChnl: make(chan fromTo, totalFloors*totalElevators*20),
	}
	return eg
}

func (eg *elevatorGroup) AddElevator(ele *elevator) {
	eg.elevators = append(eg.elevators, ele)
}

func (eg *elevatorGroup) selectBestAndAdd(req fromTo) {

	disAEleChnl := make(chan disAEle, eg.totFloors)
	var wg sync.WaitGroup
	for _, ele := range eg.elevators {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chnl := make(chan stopsAndCurrFloor)
			ele.getStopsAndCurrFloorChnl <- chnl
			stopsAndCurrFloor := <-chnl
			stops := ele.AddStopsAndGet(req.from, req.to, stopsAndCurrFloor.currFloor, stopsAndCurrFloor.stops)
			dis := ele.findDistance(stops, stopsAndCurrFloor.currFloor)
			disAEleChnl <- disAEle{dis: dis, ele: ele}
		}()
	}

	go func() {
		defer close(disAEleChnl)
		wg.Wait()
	}()

	var disAEles []disAEle

	for de := range disAEleChnl {
		disAEles = append(disAEles, de)
	}

	slices.SortFunc(disAEles, func(a, b disAEle) int {
		if a.dis < b.dis {
			return -1
		} else if a.dis > b.dis {
			return 1
		}
		return 0
	})

	isAdded := false
	for _, de := range disAEles {
		isAddedChnl := make(chan bool)
		defer close(isAddedChnl)
		de.ele.tryAddStopsChnl <- fromToIsAdded{fromTo: fromTo{from: req.from, to: req.to}, dis: de.dis, isAdded: isAddedChnl}
		isAdded = <-isAddedChnl
		if isAdded {
			break
		}
	}

	if !isAdded {
		fmt.Println("No elevator available to serve the request from", req.from, "to", req.to, "... retrying")
		eg.Serve(req.from, req.to)
	}

}

func (eg *elevatorGroup) Serve(from, to int) {
	eg.fromToChnl <- fromTo{from: from, to: to}
}

func (eg *elevatorGroup) Start() {
	defer close(eg.fromToChnl)
	for req := range eg.fromToChnl {
		go eg.selectBestAndAdd(req)
	}
}
