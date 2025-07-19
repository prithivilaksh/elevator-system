# Elevator Control System

A concurrent elevator control system that manages multiple elevators, efficiently handling floor requests while optimizing for performance and reliability.

## Requirements

1. Multiple elevators operating concurrently
2. Each elevator maintains its own state (stops and current floor)
3. Concurrent operation of elevator group and button panels
4. Button panel at each floor for requesting elevators
5. No buttons inside the elevators (only floor selection)
6. Configurable delay for floor-to-floor movement simulation

## System Architecture

The system follows a decoupled architecture with these main components:

1. **Elevator**: Individual elevator implementation
2. **ElevatorGroup**: Manages multiple elevators and routes requests
3. **Types**: Shared types and interfaces
4. **Utils**: Common utility functions

## Core Components

### 1. Elevator
- **Type**: `elevator`
- **Responsibility**: Manages individual elevator operations
- **Key Features**:
  - Tracks current floor and stops
  - Handles door operations
  - Simulates movement between floors
  - Processes new stop requests and concurrently handles them

### 2. ElevatorGroup
- **Type**: `elevatorGroup`
- **Responsibility**: Manages multiple elevators and distributes requests
- **Key Features**:
  - Maintains a pool of elevators
  - Implements request distribution logic
  - Handles concurrent requests
  - Implements elevator selection algorithm

### 3. Types
- **Package**: `types`
- **Purpose**: Defines shared data structures and interfaces
- **Key Types**:
  - `Elevator` interface
  - `FromTo` (request structure)
  - `StopsACurrFloor` (elevator state)

## Key Features

1. **Concurrent Operation**
   - Each elevator operates in its own goroutine
   - Non-blocking request handling
   - Thread-safe state management

2. **Request Handling**
   - Concurrent request processing
   - Smart elevator selection based on distance
   - Request queuing and retry mechanism

3. **Movement Simulation**
   - Realistic timing for floor-to-floor movement
   - Door open/close simulation
   - Stop announcement and status updates

## SOLID Principles

The system adheres to the SOLID principles of object-oriented design:

1. **Single Responsibility Principle (SRP)**
   - `elevator` handles only elevator-specific operations
   - `elevatorGroup` focuses on request distribution and elevator coordination
   - Each type has a single reason to change

2. **Open/Closed Principle (OCP)**
   - New elevator types can be added without modifying existing code
   - Extensible through composition rather than modification

3. **Liskov Substitution Principle (LSP)**
   - All elevators implement the `Elevator` interface completely
   - Derived types are substitutable for their base types
   - Interface contracts are strictly followed

4. **Interface Segregation Principle (ISP)**
   - Small, focused interfaces
   - Clients only depend on methods they use
   - No "fat" interfaces forcing implementation of unused methods

5. **Dependency Inversion Principle (DIP)**
   - High-level modules don't depend on low-level modules; both depend on abstractions
   - Abstractions don't depend on details; details depend on abstractions
   - Dependencies flow toward the center (domain) of the application

## Design Patterns

1. **Worker Pattern**
   - Each elevator is a worker processing its own queue
   - Non-blocking channel communication

2. **Dependency Injection**
   - Elevators injected into ElevatorGroup
   - Interface-based design for testability

3. **Immutable Data**
   - Shared data structures are passed by value
   - Thread-safe state management

## Interfaces

### Elevator Interface
```go
type Elevator interface {
    GetStopsAndCurrFloor() StopsACurrFloor
    AddStopAndGetNextInd(startInd, afterFloor, stop int, stops []int) ([]int, int)
    FindDistance(stops []int, currFloor int, lastInd int) int
    TryAddStops(req FromTo, dis int) bool
}
```

## Usage Example

```go
// Initialize system
totalFloors := 10
totalElevators := 4

// Create elevators
var elevators []types.Elevator
for i := 0; i < totalElevators; i++ {
    elevators = append(elevators, elevator.NewElevator(fmt.Sprintf("Elevator-%d", i+1), totalFloors))
}

// Create elevator group
elevatorGroup := elevatorGroup.NewElevatorGroup(totalFloors, elevators)

// Make requests
elevatorGroup.Serve(0, 9)   // Request from floor 0 to 9
elevatorGroup.Serve(5, 8)   // Request from floor 5 to 8
```

## Error Handling

- Failed requests are automatically retried
- System maintains state consistency
- Graceful shutdown on interrupt signals

## Performance Considerations

- Uses channels for safe concurrent access
- Efficient stop list management
- Minimal locking for better concurrency

## Future Enhancements

1. Add more sophisticated elevator selection algorithms
2. Implement priority for emergency stops
3. Add metrics collection
4. Implement persistent storage for state recovery
5. Add more detailed logging and monitoring
6. Add Interface for elevator group to be used by other modules

## Dependencies

- Go 1.18+ (for generics support)
- Standard library only (no external dependencies)


