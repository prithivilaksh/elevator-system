package internal

import (
	"fmt"
	"slices"
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

func (ele *elevator) openClose() {
	ele.isOpen = true
	fmt.Println("Elevator " + ele.name + " is opening doors")
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " is closing doors")
	ele.isOpen = false
}

func (ele *elevator) sendFloors() {
	for range ele.readyForNext {
		if len(ele.stops) == 0 {
			fmt.Printf("Elevator %s has no stops to make so it is stopping\n", ele.name)
			break
		}
		if Abs(ele.currFloor-ele.stops[0]) == 1 {
			stop, err := ele.popStops()
			if !err {
				ele.floorsChnl <- floorAStop{floor: stop, isAStop: true}
			}

		}
		if ele.currFloor < ele.stops[0] {
			ele.floorsChnl <- floorAStop{floor: ele.currFloor + 1, isAStop: false}
		} else if ele.currFloor > ele.stops[0] {
			ele.floorsChnl <- floorAStop{floor: ele.currFloor - 1, isAStop: false}
		} else {
			panic("Elevator " + ele.name + " has no stops to make")
		}
	}
}

func (ele *elevator) start() {
	ele.readyForNext <- true
	for floorStop := range ele.floorsChnl {
		to := floorStop.floor
		isAStop := floorStop.isAStop
		time.Sleep(5 * time.Second) // Simulate time taken to reach the floor

		func() {
			ele.mu.Lock()
			ele.currFloor = to
			ele.mu.Unlock()
		}()

		fmt.Printf("Elevator "+ele.name+" reached floor %d\n", ele.currFloor)
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

func (ele *elevator) addStops(req request, expDis int) bool {
	ele.mu.Lock()
	defer ele.mu.Unlock()

	currFloor, from, to := ele.currFloor, req.from, req.to
	stops := ele.stops

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

	actDis := 0
	prev = currFloor
	for i := 0; i < len(stops)-1; i++ {
		actDis += Abs(stops[i] - prev)
		prev = stops[i]
	}

	if Abs(actDis-expDis) > 4 {
		fmt.Printf("Elevator %s cannot take request from %d to %d as it will take %d distance but expected %d\n", ele.name, from, to, actDis, expDis)
		return false
	}

	if len(ele.stops) == 2 {
		go ele.sendFloors()
	}
	fmt.Printf("Elevator %s added stops from %d to %d with expected distance %d\n", ele.name, from, to, expDis)
	return true

}
