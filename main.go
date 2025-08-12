package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/prithivilaksh/elevator-system/elevator"
	"github.com/prithivilaksh/elevator-system/elevatorgroup"
)

func main() {
	metricElevatorGroup := elevatorgroup.NewMetricElevatorGroup()

	for i := 1; i <= 4; i++ {
		metricElevatorGroup.AddElevator(elevator.NewLeastDisElevator(i))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Running... Press Ctrl+C to stop")

	go metricElevatorGroup.GetElevatorID(1, 10)
	go metricElevatorGroup.GetElevatorID(4, 9)
	go metricElevatorGroup.GetElevatorID(7, 3)
	go metricElevatorGroup.GetElevatorID(2, 6)
	go metricElevatorGroup.GetElevatorID(5, 8)
	go metricElevatorGroup.GetElevatorID(8, 2)
	go metricElevatorGroup.GetElevatorID(1, 3)
	go metricElevatorGroup.GetElevatorID(5, 10)
	go metricElevatorGroup.GetElevatorID(2, 4)
	go metricElevatorGroup.GetElevatorID(4, 10)
	go metricElevatorGroup.GetElevatorID(10, 2)
	go metricElevatorGroup.GetElevatorID(2, 1)
	go metricElevatorGroup.GetElevatorID(1, 4)
	go metricElevatorGroup.GetElevatorID(9, 10)
	go metricElevatorGroup.GetElevatorID(3, 9)
	go metricElevatorGroup.GetElevatorID(5, 3)
	go metricElevatorGroup.GetElevatorID(7, 8)
	go metricElevatorGroup.GetElevatorID(3, 10)
	go metricElevatorGroup.GetElevatorID(4, 6)
	go metricElevatorGroup.GetElevatorID(5, 4)
	go metricElevatorGroup.GetElevatorID(2, 9)
	go metricElevatorGroup.GetElevatorID(4, 9)

	<-sigChan
	fmt.Println("Shutting down...")
}
