package elevator

import (
	// "fmt"
	"slices"
	"sync"
	"time"

	"github.com/prithivilaksh/elevator-system/utils"
)

type LeastDisElevator struct {
	id        int
	stops     []int
	currFloor int
	isIdle    bool
	mu        *sync.RWMutex
}

func NewLeastDisElevator(id int) *LeastDisElevator {
	return &LeastDisElevator{
		id:        id,
		stops:     make([]int, 0),
		currFloor: 0,
		isIdle:    true,
		mu:        &sync.RWMutex{},
	}
}

func (e *LeastDisElevator) GetID() int {
	return e.id
}

func addStopAndGetNextInd(startInd, afterFloor, stop int, stops []int) ([]int, int) {

	nextStartInd := 0
	inserted := false
	for i := startInd; i < len(stops); i++ {
		if (afterFloor <= stop && stop <= stops[i]) || (afterFloor >= stop && stop >= stops[i]) {
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

func addStopsAndFindDistance(from, to, currFloor int, stops []int) ([]int, int) {
	nextStartInd := 0
	stops, nextStartInd = addStopAndGetNextInd(nextStartInd, currFloor, from, stops)
	stops, _ = addStopAndGetNextInd(nextStartInd, from, to, stops)
	prev, dis := currFloor, 0
	vis := false
	for _, x := range stops {
		dis += utils.Abs(x - prev)
		prev = x
		if vis && x == to {
			break
		}
		if x == from {
			vis = true
		}
	}
	return stops, dis
}

func (e *LeastDisElevator) TryAddStops(from, to, expMetric, threshold int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	stops := utils.Copy(e.stops)
	currFloor := e.currFloor

	stops, actDis := addStopsAndFindDistance(from, to, currFloor, stops)

	if utils.Abs(actDis-expMetric) <= threshold {
		// fmt.Println("elevator ", e.id, " added stops from ", from, " to ", to, " actual metric ", actDis, " expected metric ", expMetric, " before stops ", e.stops, " after stops ", stops)
		e.stops = stops
		if e.isIdle {
			e.isIdle = false
			go e.Simulate()
		}
		return true
	}
	return false
}

func (e *LeastDisElevator) GetMetric(from, to int) int {
	e.mu.RLock()
	stops := utils.Copy(e.stops)
	currFloor := e.currFloor
	e.mu.RUnlock()

	_, dis := addStopsAndFindDistance(from, to, currFloor, stops)
	return dis
}

func (e *LeastDisElevator) Simulate() {
	for {
		e.mu.Lock()
		if len(e.stops) == 0 {
			e.isIdle = true
			e.mu.Unlock()

			// fmt.Println("elevator ", e.id, " is idle", e.stops)
			break
		}
		if utils.Abs(e.currFloor-e.stops[0]) <= 1 {
			if utils.Abs(e.currFloor-e.stops[0]) == 1 {
				time.Sleep(2 * time.Second)
			}
			e.currFloor = e.stops[0]
			e.stops = e.stops[1:]
		} else if e.currFloor < e.stops[0] {
			e.currFloor++
		} else if e.currFloor > e.stops[0] {
			e.currFloor--
		} else {
			panic("elevator in inconsistent state")
		}
		e.mu.Unlock()
		time.Sleep(2 * time.Second)
		// fmt.Println("elevator ", e.id, "Current Floor", e.currFloor)
	}
}
