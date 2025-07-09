package internal

import (
	"time"
)

type elevator struct {
	currentFloor int
	destination  int
	direction    string
	isOpen       bool
	stops        []bool
	
	destinations chan int
	requests     chan int
	topFloor     int
	maxFloor     int
	minFloor     int
}

func (e *elevator) move(dis int) {
	for dest := range e.destinations {
		e.updateDirection(dest)
		for e.currentFloor != dest {
			time.Sleep(5 * time.Second)
			if e.direction == "up" {
				e.currentFloor++
			}
			if e.direction == "down" {
				e.currentFloor--
			}
			if e.stops[e.currentFloor] {
				e.openClose()
				e.stops[e.currentFloor] = false
			}
		}
		e.openClose()
		e.updateDirection(e.currentFloor)
	}

}
func (e *elevator) updateDirection(dest int) {
	if e.currentFloor < dest {
		e.direction = "up"
		e.maxFloor = dest
	} else if e.currentFloor > dest {
		e.direction = "down"
		e.minFloor = dest
	} else {
		e.direction = "rest"
	}
}

func (e *elevator) openClose() {
	e.isOpen = true
	time.Sleep(3 * time.Second)
	e.isOpen = false
}

func (e *elevator) start() {

	for dest := range e.requests {
		if e.currentFloor <= dest && dest <= e.maxFloor && e.direction == "up" {
			e.stops[dest] = true
		} else if e.minFloor <= dest && dest <= e.currentFloor && e.direction == "down" {
			e.stops[dest] = true
		} else {
			e.destinations <- dest
		}

	}

}
