package elevator

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	. "github.com/prithivilaksh/elevator-system/utils"
	. "github.com/prithivilaksh/elevator-system/types"
)

type FromToDisIsAdded struct {
	FromTo
	dis         int
	isAddedChnl chan bool
}

type elevator struct {
	name                     string
	currFloor                int
	isIdle                   bool
	stops                    []int
	readyForNextChnl         chan bool
	getStopsAndCurrFloorChnl chan (chan StopsACurrFloor)
	tryAddStopsChnl          chan FromToDisIsAdded
}

func NewElevator(name string, totalFloors int) *elevator {
	ele := &elevator{
		name:                     name,
		currFloor:                0,
		isIdle:                   true,
		stops:                    make([]int, 0, totalFloors*2),
		getStopsAndCurrFloorChnl: make(chan (chan StopsACurrFloor), totalFloors*2),
		tryAddStopsChnl:          make(chan FromToDisIsAdded, totalFloors*2),
		readyForNextChnl:         make(chan bool, totalFloors*2),
	}
	go ele.Start()
	return ele
}

func (ele *elevator) AddStopAndGetNextInd(startInd, prev, stop int, stops []int) ([]int, int) {

	nextStartInd := 0
	inserted := false
	for i := startInd; i < len(stops); i++ {
		if (prev <= stop && stop <= stops[i]) || (prev >= stop && stop >= stops[i]) {
			if stop != stops[i] {
				stops = slices.Insert(stops, i, stop)
			}
			nextStartInd = i + 1
			inserted = true
			break
		}
	}

	if !inserted {
		stops = append(stops, stop)
		nextStartInd = len(stops)
	}

	return stops, nextStartInd
}

func (ele *elevator) FindDistance(stops []int, prev int, lastInd int) int {
	// tmp:=prev
	dis := 0
	for i := range lastInd {
		dis += Abs(stops[i] - prev)
		prev = stops[i]
	}
	// fmt.Println("prev", tmp, "stops", stops, "dis", dis, "lastInd", lastInd)
	return dis
}

func (ele *elevator) GetStopsAndCurrFloor() StopsACurrFloor {
	stopsAndCurrFloorChnl := make(chan StopsACurrFloor)
	defer close(stopsAndCurrFloorChnl)
	ele.getStopsAndCurrFloorChnl <- stopsAndCurrFloorChnl
	return <-stopsAndCurrFloorChnl
}

func (ele *elevator) getStopsAndCurrFloor(stopsAndCurrFloorChnl chan StopsACurrFloor) {
	stops := DeepCopy(ele.stops)
	stopsAndCurrFloorChnl <- StopsACurrFloor{Stops: stops, CurrFloor: ele.currFloor}
}

func (ele *elevator) TryAddStops(req FromTo, dis int) bool {
	isAddedChnl := make(chan bool)
	defer close(isAddedChnl)
	ele.tryAddStopsChnl <- FromToDisIsAdded{FromTo: req, dis: dis, isAddedChnl: isAddedChnl}
	return <-isAddedChnl
}

func (ele *elevator) tryAddStops(req FromToDisIsAdded) {
	stops, isAdded, nextStartInd := DeepCopy(ele.stops), false, 0
	stops, nextStartInd = ele.AddStopAndGetNextInd(nextStartInd, ele.currFloor, req.From, stops)
	stops, nextStartInd = ele.AddStopAndGetNextInd(nextStartInd, req.From, req.To, stops)
	dis := ele.FindDistance(stops, ele.currFloor, nextStartInd)
	if Abs(req.dis-dis) <= 2 {
		fmt.Println("Request from", req.From, "to", req.To, "assigned to elevator", ele.name, "with distance", dis, "and current distance", req.dis, "stops before", ele.stops, " current floor", ele.currFloor, " stops after", stops)
		ele.stops = stops
		isAdded = true
		if ele.isIdle {
			ele.readyForNextChnl <- true
			ele.isIdle = false
			fmt.Println("Elevator " + ele.name + " started")
		}
	}
	req.isAddedChnl <- isAdded
}

func (ele *elevator) simulate(to int, isAStop bool) {
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " moved to floor " + strconv.Itoa(to))
	if isAStop {
		fmt.Println("Elevator " + ele.name + " is opening doors")
		time.Sleep(5 * time.Second)
		fmt.Println("Elevator " + ele.name + " is closing doors")
	}
	ele.readyForNextChnl <- true
}

func (ele *elevator) moveToNextFloor() {
	if len(ele.stops) == 0 {
		ele.isIdle = true
		fmt.Println("Elevator " + ele.name + " stopped at floor " + strconv.Itoa(ele.currFloor))
		return
	}
	isAStop := false
	if Abs(ele.currFloor-ele.stops[0]) <= 1 {
		ele.currFloor = ele.stops[0]
		ele.stops = ele.stops[1:]
		isAStop = true
	} else if ele.currFloor > ele.stops[0] {
		ele.currFloor--
	} else {
		ele.currFloor++
	}
	go ele.simulate(ele.currFloor, isAStop)

}

func (ele *elevator) Start() {
	defer close(ele.readyForNextChnl)
	defer close(ele.getStopsAndCurrFloorChnl)
	defer close(ele.tryAddStopsChnl)
	for {
		select {
		case stopsAndCurrFloorChnl := <-ele.getStopsAndCurrFloorChnl:
			ele.getStopsAndCurrFloor(stopsAndCurrFloorChnl)

		case <-ele.readyForNextChnl:
			ele.moveToNextFloor()

		case req := <-ele.tryAddStopsChnl:
			ele.tryAddStops(req)
		}
	}
}
