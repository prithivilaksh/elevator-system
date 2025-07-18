# Elevator Group

## Requirements

1. Can have multiple elevators.
2. Each one has its own direction, currentFloor, destinations slice, door status etc.
3. The elevators should operate concurrently.
4. Button group should be present once at each floor which can operate concurrently.
5. Button group consist of X Selector at each floor.
6. Elevator group should have number of Floors x Button group buttons totally
7. No buttons inside the lift.
9. To simulate elevator movement add a delay of Y seconds to move from floor a to a+1 or vice versa




# Elevator Group

### Creates a new elevator group
NewElevatorGroup(totalFloors int, maxReqs int, totalElevators int) *elevatorGroup

### Adds an elevator to the group
AddElevator(ele *elevator)

### Accepts requests and sends to reqs channel
Request(from,to int)

### Reads requests from reqs channel and processes them in a separate go routine
ProcessRequests()

### Selects the best elevator and assigns the request to it
SelectBestAndAssign()

### Start - Entry point
Start()




# Elevator

### Creates a new elevator
NewElevator(name string, totalFloors int) *elevator

### Returns a deep copy of stops and current floor
GetStopsAndCurrFloor() []int, int

### Inserts a request(from,to) to the copy of stops
AddStopsAndGet(from, to , currFloor int, stops []int ) []int

### Finds the distance between stops and also the current floor
FindDistance(stops []int, currFloor int) int


<!-- ### Gets the current floor
GetCurrFloor() int -->

### Move - Moves the elevator to the destination
Move(to int)

### Manages the state of the elevator
ManageState()


### Start - Entry point
Start()


