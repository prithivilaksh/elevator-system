package internal

import (
	"fmt"
	"slices"
	"strconv"
	"time"
)

type elevator struct {
	name         string
	currFloor    int
	stops []int
	readyForNextChnl chan bool
	getStopsAndCurrFloorChnl chan chan stopsAndCurrFloor
	tryAddStopsChnl chan fromToIsAdded
}

func NewElevator(name string, totalFloors int) *elevator {
	return &elevator{
		name:         name,
		currFloor:    0,
		stops:        []int{},
		getStopsAndCurrFloorChnl: make(chan chan stopsAndCurrFloor, totalFloors*2),
		tryAddStopsChnl: make(chan fromToIsAdded, totalFloors*2),
		readyForNextChnl: make(chan bool, totalFloors*2),
	}
}

func (ele *elevator) openClose() {
	fmt.Println("Elevator " + ele.name + " is opening doors")
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " is closing doors")
}


func (ele *elevator) AddStopsAndGet(from, to int, currFloor int, stops []int ) []int {

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
	for i := 0; i < len(stops); i++ {
		dis += Abs(stops[i] - prev)
		prev = stops[i]
	}
	return dis
}


func (ele *elevator) Move(to int, isAStop bool) {
	time.Sleep(5 * time.Second)
	fmt.Println("Elevator " + ele.name + " moved to floor " + strconv.Itoa(to))
	if isAStop{
		ele.openClose()
	}
	ele.readyForNextChnl <- true
}

func (ele *elevator) Start() {
	for {
		select {
		case chnl:=<-ele.getStopsAndCurrFloorChnl:
			chnl <- stopsAndCurrFloor{stops: ele.stops, currFloor: ele.currFloor}
		
		case <-ele.readyForNextChnl:
			isAStop := false
			if len(ele.stops)==0{
				continue
			}
			if Abs(ele.currFloor-ele.stops[0])<=1{
				ele.currFloor=ele.stops[0]
				ele.stops=ele.stops[1:]
				isAStop=true
			}else if ele.currFloor>ele.stops[0]{
				ele.currFloor--
			}else{
				ele.currFloor++
			}
			go ele.Move(ele.currFloor, isAStop)
		
		case ins := <-ele.tryAddStopsChnl:
			stops:=ele.AddStopsAndGet(ins.from, ins.to, ele.currFloor, ele.stops)
			dis:=ele.findDistance(stops, ele.currFloor)
			if Abs(ins.dis-dis)<=2{
				ele.stops=stops
				ins.isAdded <- true	
				if len(ele.readyForNextChnl)==0{
					ele.readyForNextChnl <- true
				}
			}else{
				ins.isAdded <- false
			}
		}
	}
}
