package internal

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"
)

type elevator struct {
	name         string
	isOpen       bool
	floorsChnl   chan floorAStop
	readyForNext chan bool
	currFloor    int

	mu    sync.Mutex
	stops []int
}

func NewElevator(name string, totalFloors int) *elevator {
	return &elevator{
		name:         name,
		isOpen:       false,
		floorsChnl:   make(chan floorAStop, totalFloors*2),
		readyForNext: make(chan bool, totalFloors*2),
		currFloor:    0,
		stops:        []int{},
	}
}

func (ele *elevator) openClose() {
	ele.isOpen = true
	fmt.Println("Elevator " + ele.name + " is opening doors")
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " is closing doors")
	ele.isOpen = false
}

func (ele *elevator) sendFloors() {
	for range ele.readyForNext {
		fmt.Println("Elevator " + ele.name + " is processing request")

		if ele.currFloor == ele.stops[0] {
			stop, err := ele.popStops()
			if !err {
				ele.floorsChnl <- floorAStop{floor: stop, isAStop: true}
			}
		}

		if len(ele.stops) == 0 {
			fmt.Println("Elevator", ele.name, "has no stops to make so it is stopping")
			break
		}
		if ele.currFloor < ele.stops[0] {
			ele.floorsChnl <- floorAStop{floor: ele.currFloor + 1, isAStop: false}
		} else if ele.currFloor > ele.stops[0] {
			ele.floorsChnl <- floorAStop{floor: ele.currFloor - 1, isAStop: false}
		} else {
			fmt.Println(ele.stops)
			panic("Elevator " + ele.name + " has no stops to make " + strconv.Itoa(ele.currFloor))
		}
	}
}

func (ele *elevator) start() {
	fmt.Println("Elevator " + ele.name + " started")
	ele.readyForNext <- true
	fmt.Println("Elevator " + ele.name + " is ready for next request")
	for floorStop := range ele.floorsChnl {
		fmt.Println("Elevator "+ele.name+" received floor stop request for floor", floorStop.floor)
		to := floorStop.floor
		isAStop := floorStop.isAStop

		if ele.currFloor != to {
			ele.mu.Lock()
			ele.currFloor = -1
			time.Sleep(5 * time.Second) // Simulate time taken to reach the floor
			ele.currFloor = to
			ele.mu.Unlock()
		}

		fmt.Println("Elevator "+ele.name+" reached floor ", ele.currFloor)

		if isAStop {
			ele.openClose()
		}
		ele.readyForNext <- true
	}
}

func (ele *elevator) popStops() (int, bool) {
	ele.mu.Lock()
	defer ele.mu.Unlock()

	if len(ele.stops) == 0 {
		return 0, true
	}

	stop := ele.stops[0]
	ele.stops = ele.stops[1:]
	return stop, false
}

func (ele *elevator) insertFromTo(from int, to int) []int {

	stops := make([]int, len(ele.stops))
	copy(stops, ele.stops)

	nextInd := 0
	prev := ele.currFloor
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

func (ele *elevator) findDistance(stops []int) int {
	prev, dis := ele.currFloor, 0
	for i := 0; i < len(stops); i++ {
		dis += Abs(stops[i] - prev)
		prev = stops[i]
	}
	return dis
}

func (ele *elevator) addStops(req Request, expDis int) bool {
	ele.mu.Lock()
	defer ele.mu.Unlock()

	stops := ele.insertFromTo(req.From, req.To)
	actDis := ele.findDistance(stops)
	if Abs(actDis-expDis) > 4 {
		fmt.Println("Elevator " + ele.name + " cannot take request from " + strconv.Itoa(req.From) + " to " + strconv.Itoa(req.To) + " as it will take " + strconv.Itoa(actDis) + " distance but expected " + strconv.Itoa(expDis))
		return false
	}

	ele.stops = stops

	if len(ele.stops) <= 2 {
		go ele.sendFloors()
	}
	fmt.Println("Elevator " + ele.name + " added stops from " + strconv.Itoa(req.From) + " to " + strconv.Itoa(req.To) + " with expected distance " + strconv.Itoa(expDis))
	return true
}

func (ele *elevator) Stop() {
	close(ele.floorsChnl)
	close(ele.readyForNext)
}
