package internal

import (
	"slices"
	"sync"
)

func newElevatorGroup(totalFloors int) *elevatorGroup {
	eg := &elevatorGroup{
		elevators: make([]*elevator, totalFloors),
		totFloors: totalFloors,
		reqsChnl:  make(chan request, totalFloors),
	}
	return eg
}

func (eg *elevatorGroup) addElevator(ele *elevator) {
	eg.elevators = append(eg.elevators, ele)
}

func (eg *elevatorGroup) serve(req request) {
	eg.reqsChnl <- req
}

func (eg *elevatorGroup) findDistance(ele *elevator, req request, disAEleChnl chan<- disAEle) {
	currFloor, from, to := ele.currFloor, req.from, req.to
	stops := make([]int, len(ele.stops))
	copy(stops, ele.stops)

	prev := currFloor
	for i := 0; i < len(stops); i++ {
		if (prev <= from && from <= stops[i]) || (prev >= from && from >= stops[i]) {
			prev = from
			if from != stops[i] {
				stops = slices.Insert(stops, i, from)
			}
			break
		}
	}

	if len(stops) == 0 {
		stops = append(stops, from)
	}
	for i := prev + 1; i < len(stops); i++ {
		if (prev <= to && to <= stops[i]) || (prev >= to && to >= stops[i]) {
			if to != stops[i] {
				stops = slices.Insert(stops, i, to)
			}
			break
		}
	}

	if len(stops) == 1 {
		stops = append(stops, to)
	}

	dis := 0
	prev = currFloor
	for i := 0; i < len(stops)-1; i++ {
		dis += Abs(stops[i] - prev)
		prev = stops[i]
	}

	disAEleChnl <- disAEle{
		dis: dis,
		ele: ele,
	}

}

func (eg *elevatorGroup) SelectBestAndAdd(req request) {

	disAEleChnl := make(chan disAEle, eg.totFloors)
	var wg sync.WaitGroup
	for _, ele := range eg.elevators {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eg.findDistance(ele, req, disAEleChnl)
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
		// If no elevator can take the request, we can either log it or handle it differently.
		// For now, we will just print a message.
		println("No elevator available to serve the request from", req.from, "to", req.to)
		eg.serve(req)
	}

}



func (eg *elevatorGroup) start() {
	for _, ele := range eg.elevators {
		go ele.start()
	}
	for req := range eg.reqsChnl {
		go eg.SelectBestAndAdd(req)
	}
}
