package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/prithivilaksh/elevator-group/internal"
)

func main() {
	fmt.Println("Elevator Group Simulation Started")

	totalFloors := 20
	maxReqs := 100
	totalElevators := 5

	elevatorGroup := internal.NewElevatorGroup(totalFloors, maxReqs, totalElevators)

	for i := range totalElevators {
		elevator := internal.NewElevator(fmt.Sprintf("Elevator-%d", i+1), maxReqs)
		elevatorGroup.AddElevator(elevator)
	}

	elevatorGroup.Start()
	elevatorGroup.Serve(internal.Request{From: 0, To: 9})
	elevatorGroup.Serve(internal.Request{From: 5, To: 8})
	elevatorGroup.Serve(internal.Request{From: 7, To: 13})
	// elevatorGroup.Serve(internal.Request{From: 3, To: 5})
	// elevatorGroup.Serve(internal.Request{From: 4, To: 6})
	// elevatorGroup.Serve(internal.Request{From: 5, To: 1})
	// elevatorGroup.Serve(internal.Request{From: 6, To: 3})
	// elevatorGroup.Serve(internal.Request{From: 7, To: 8})
	// elevatorGroup.Serve(internal.Request{From: 8, To: 0})
	// elevatorGroup.Serve(internal.Request{From: 9, To: 7})

	// elevatorGroup.Serve(internal.Request{From: 0, To: 2})
	// elevatorGroup.Serve(internal.Request{From: 1, To: 9})
	// elevatorGroup.Serve(internal.Request{From: 2, To: 5})
	// elevatorGroup.Serve(internal.Request{From: 3, To: 6})
	// elevatorGroup.Serve(internal.Request{From: 4, To: 1})
	// elevatorGroup.Serve(internal.Request{From: 5, To: 3})
	// elevatorGroup.Serve(internal.Request{From: 6, To: 8})
	// elevatorGroup.Serve(internal.Request{From: 7, To: 0})
	// elevatorGroup.Serve(internal.Request{From: 8, To: 7})
	// elevatorGroup.Serve(internal.Request{From: 9, To: 4})

	
	// time.Sleep(20 * time.Second) // Allow some time for the elevators to process requests
	sigChan := make(chan os.Signal, 1)

	// Notify the channel on Interrupt (Ctrl+C) or SIGTERM
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Running... Press Ctrl+C to stop")

	// Block until a signal is received
	<-sigChan
	elevatorGroup.Stop()
	// time.Sleep(20 * time.Second) // Allow some time for the elevators to stop gracefully

	fmt.Println("\nTermination signal received. Exiting gracefully.")

}

