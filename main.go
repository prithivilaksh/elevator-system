package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Elevator Group Simulation Started")

	totalFloors := 10
	totalElevators := 4

	var elevators []*elevator
	for i := range totalElevators {
		elevator := NewElevator(fmt.Sprintf("Elevator-%d", i+1), totalFloors)
		elevators = append(elevators, elevator)
	}

	elevatorGroup := NewElevatorGroup(totalFloors, elevators)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Running... Press Ctrl+C to stop")

	go elevatorGroup.Serve(0, 9)
	go elevatorGroup.Serve(5, 8)
	go elevatorGroup.Serve(3, 5)
	go elevatorGroup.Serve(4, 6)
	go elevatorGroup.Serve(5, 1)
	go elevatorGroup.Serve(6, 3)
	go elevatorGroup.Serve(7, 8)
	go elevatorGroup.Serve(8, 0)
	go elevatorGroup.Serve(9, 7)
	go elevatorGroup.Serve(0, 2)
	go elevatorGroup.Serve(1, 9)
	go elevatorGroup.Serve(2, 5)
	go elevatorGroup.Serve(3, 6)
	go elevatorGroup.Serve(4, 1)
	go elevatorGroup.Serve(5, 3)
	go elevatorGroup.Serve(6, 8)
	go elevatorGroup.Serve(7, 0)
	go elevatorGroup.Serve(8, 7)
	go elevatorGroup.Serve(9, 4)

	<-sigChan

	go elevatorGroup.Serve(7, 8)
	go elevatorGroup.Serve(0, 9)
	go elevatorGroup.Serve(5, 8)
	go elevatorGroup.Serve(2, 5)
	go elevatorGroup.Serve(4, 1)
	go elevatorGroup.Serve(3, 5)
	go elevatorGroup.Serve(6, 3)
	go elevatorGroup.Serve(8, 0)
	go elevatorGroup.Serve(0, 2)
	go elevatorGroup.Serve(1, 9)
	go elevatorGroup.Serve(3, 6)
	go elevatorGroup.Serve(5, 3)
	go elevatorGroup.Serve(4, 6)
	go elevatorGroup.Serve(6, 8)
	go elevatorGroup.Serve(7, 0)
	go elevatorGroup.Serve(5, 1)
	go elevatorGroup.Serve(8, 7)
	go elevatorGroup.Serve(9, 7)
	go elevatorGroup.Serve(9, 4)

	<-sigChan

	fmt.Println("\nTermination signal received. Exiting gracefully.")
}
