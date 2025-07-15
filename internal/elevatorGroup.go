package internal

import (
	"slices"
	"sync"
)

type elevatorGroup struct {
	elevators []*elevator
	totFloors int
	reqsChnl  chan Request
}

func NewElevatorGroup(totalFloors int, maxReqs int, totalElevators int) *elevatorGroup {
	eg := &elevatorGroup{
		elevators: make([]*elevator, 0, totalElevators),
		totFloors: totalFloors,
		reqsChnl:  make(chan Request, maxReqs),
	}
	return eg
}

func (eg *elevatorGroup) AddElevator(ele *elevator) {
	eg.elevators = append(eg.elevators, ele)
}

func (eg *elevatorGroup) Serve(req Request) {
	eg.reqsChnl <- req
}

func (eg *elevatorGroup) selectBestAndAdd(req Request) {

	disAEleChnl := make(chan disAEle, eg.totFloors)
	var wg sync.WaitGroup
	for _, ele := range eg.elevators {
		wg.Add(1)
		go func() {
			ele.mu.Lock()
			defer ele.mu.Unlock()
			defer wg.Done()
			stops := ele.insertFromTo(req.From, req.To)
			dis := ele.findDistance(stops)
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
		isAdded = de.ele.addStops(req, de.dis)
		if isAdded {
			break
		}
	}

	if !isAdded {
		println("No elevator available to serve the request from", req.From, "to", req.To)
		eg.Serve(req)
	}

}

func (eg *elevatorGroup) processReqs() {
	for req := range eg.reqsChnl {
		go eg.selectBestAndAdd(req)
	}
}

func (eg *elevatorGroup) Start() {
	for _, ele := range eg.elevators {
		go ele.start()
	}
	go eg.processReqs()
}

func (eg *elevatorGroup) Stop() {
	close(eg.reqsChnl)
	for _, elevator := range eg.elevators {
		elevator.Stop()
	}
}
