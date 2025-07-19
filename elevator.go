package main

import (
	"fmt"
	"slices"
	"strconv"
	"time"
)

type elevator struct {
	name                     string
	currFloor                int
	isIdle                   bool
	stops                    []int
	readyForNextChnl         chan bool
	getStopsAndCurrFloorChnl chan chan stopsAndCurrFloor
	tryAddStopsChnl          chan fromToIsAdded
}

func NewElevator(name string, totalFloors int) *elevator {
	ele := &elevator{
		name:                     name,
		currFloor:                0,
		isIdle:                   true,
		stops:                    make([]int, 0, totalFloors*2),
		getStopsAndCurrFloorChnl: make(chan (chan stopsAndCurrFloor), totalFloors*2),
		tryAddStopsChnl:          make(chan fromToIsAdded, totalFloors*2),
		readyForNextChnl:         make(chan bool, totalFloors*2),
	}
	go ele.Start()
	return ele
}

func (ele *elevator) AddStopsAndGet(from, to int, currFloor int, stops []int) []int {

	nextInd := 0
	prev := currFloor
	inserted := false
	for i := 0; i < len(stops); i++ {
		if (prev <= from && from <= stops[i]) || (prev >= from && from >= stops[i]) {
			nextInd = i + 1
			if from != stops[i] {
				stops = slices.Insert(stops, i, from)
			}
			inserted = true
			break
		}
	}

	if !inserted {
		stops = append(stops, from)
		nextInd = len(stops)
	}

	prev = from
	inserted = false
	for i := nextInd; i < len(stops); i++ {
		if (prev <= to && to <= stops[i]) || (prev >= to && to >= stops[i]) {
			if to != stops[i] {
				stops = slices.Insert(stops, i, to)
			}
			inserted = true
			break
		}
	}

	if !inserted {
		stops = append(stops, to)
	}

	return stops
}

func (ele *elevator) findDistance(stops []int, currFloor int) int {
	prev, dis := currFloor, 0
	for _, stop := range stops {
		dis += Abs(stop - prev)
		prev = stop
	}
	return dis
}

func (ele *elevator) Move(to int, isAStop bool) {
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " moved to floor " + strconv.Itoa(to))
	if isAStop {
		fmt.Println("Elevator " + ele.name + " is opening doors")
		time.Sleep(5 * time.Second)
		fmt.Println("Elevator " + ele.name + " is closing doors")
	}
	ele.readyForNextChnl <- true
}

func (ele *elevator) Start() {
	defer close(ele.readyForNextChnl)
	defer close(ele.getStopsAndCurrFloorChnl)
	defer close(ele.tryAddStopsChnl)
	for {
		select {
		case chnl := <-ele.getStopsAndCurrFloorChnl:
			stops := make([]int, len(ele.stops))
			copy(stops, ele.stops)
			chnl <- stopsAndCurrFloor{stops: stops, currFloor: ele.currFloor}

		case <-ele.readyForNextChnl:
			if len(ele.stops) == 0 {
				ele.isIdle = true
				fmt.Println("Elevator " + ele.name + " stopped at floor " + strconv.Itoa(ele.currFloor))
				continue
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
			go ele.Move(ele.currFloor, isAStop)

		case ins := <-ele.tryAddStopsChnl:
			x := make([]int, len(ele.stops))
			copy(x, ele.stops)
			stops := make([]int, len(ele.stops))
			copy(stops, ele.stops)
			stops = ele.AddStopsAndGet(ins.from, ins.to, ele.currFloor, stops)
			dis := ele.findDistance(stops, ele.currFloor)
			if Abs(ins.dis-dis) <= 2 {
				ele.stops = stops
				fmt.Println("Request from", ins.from, "to", ins.to, "assigned to elevator", ele.name, "with distance", dis, "and current distance", ins.dis, "stops before", x, " current floor", ele.currFloor, " stops after", stops)
				ins.isAdded <- true
				if ele.isIdle {
					ele.isIdle = false
					fmt.Println("Elevator " + ele.name + " started")
					ele.readyForNextChnl <- true
				}
			} else {
				ins.isAdded <- false
			}
		}
	}
}
