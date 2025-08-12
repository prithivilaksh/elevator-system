# Elevator Control System

A concurrent elevator control system that manages multiple elevators, efficiently handling floor requests while optimizing for performance and reliability. The system is implemented in Go and supports different elevator scheduling algorithms.

## Features

- **Concurrent Operation**: Multiple elevators operate simultaneously in their own goroutines
- **Thread-Safe**: Uses mutexes for safe concurrent access to shared resources
- **Extensible Design**: Supports different elevator scheduling algorithms
- **Efficient Request Handling**: Smart elevator selection based on configurable metrics
- **Simulation Ready**: Includes simulation capabilities for testing and demonstration

## System Architecture

The system follows a clean architecture with these main components:

1. **Elevator**: Interface and implementations for individual elevators
2. **ElevatorGroup**: Manages multiple elevators and routes requests
3. **Utils**: Common utility functions

## Core Components

### 1. Elevator Interface
```go
type Elevator interface {
    GetID() int
    TryAddStops(from, to, expMetric, threshold int) bool
    GetMetric(from, to int) int
}
```

### 2. Implementations

#### Least Distance Elevator (`LeastDisElevator`)
- **Package**: `elevator`
- **Algorithm**: Selects elevators based on the shortest distance to the requested floor
- **Features**:
  - Tracks current floor and stops
  - Implements efficient stop management
  - Thread-safe operations with mutexes

### 3. Elevator Group

#### MetricElevatorGroup
- **Package**: `elevatorgroup`
- **Responsibility**: Manages multiple elevators and distributes requests
- **Features**:
  - Maintains a pool of elevators
  - Implements concurrent request processing
  - Uses a metric-based approach to select the best elevator

## Usage Example

```go
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
	// Create a new metric-based elevator group
	metricElevatorGroup := elevatorgroup.NewMetricElevatorGroup()

	// Add elevators to the group
	for i := 1; i <= 4; i++ {
		metricElevatorGroup.AddElevator(elevator.NewLeastDisElevator(i))
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Elevator system running... Press Ctrl+C to stop")

	// Example requests (in a real application, these would come from user input or API)
	go metricElevatorGroup.GetElevatorID(1, 10)  // Request elevator from floor 1 to 10
	go metricElevatorGroup.GetElevatorID(4, 9)   // Request elevator from floor 4 to 9

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("Shutting down elevator system...")
}
```

## How It Works

1. **Request Handling**:
   - When a floor request is made via `GetElevatorID(from, to)`, the system:
     1. Calculates a metric for each elevator using `GetMetric(from, to)`
     2. Sorts elevators based on their metrics
     3. Attempts to assign the request to the best elevator using `TryAddStops()`

2. **Elevator Selection**:
   - The system uses a metric-based approach to select the most appropriate elevator
   - The `LeastDisElevator` implementation selects the elevator that will have the shortest travel distance

3. **Concurrency**:
   - Each elevator operates in its own goroutine
   - Mutexes are used to ensure thread safety when accessing shared resources

## Error Handling

- Failed requests are automatically retried with the next best elevator
- The system maintains consistency even with concurrent requests
- Graceful shutdown on interrupt signals

## Performance Considerations

- Uses mutexes for thread-safe operations
- Efficient stop list management using slices
- Minimal locking for better concurrency

## Future Enhancements

1. Implement additional elevator algorithms (e.g., least time, round-robin)
2. Add more sophisticated metrics for elevator selection
3. Implement priority for emergency stops
4. Add metrics collection and monitoring
5. Add persistent storage for state recovery
6. Implement a more sophisticated simulation mode

## Dependencies

- Go 1.18+ (for generics support)
- Standard library only (no external dependencies)

## Running the Example

1. Make sure you have Go 1.18 or later installed
2. Clone the repository
3. Run the main program:
   ```bash
   go run main.go
   ```
4. The system will process the example requests and wait for an interrupt signal

## License

This project is open source and available under the [MIT License](LICENSE).

