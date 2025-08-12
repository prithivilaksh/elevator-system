package elevatorgroup

import "github.com/prithivilaksh/elevator-system/elevator"

type ElevatorGroup interface {
	AddElevator(elevator *elevator.Elevator)
	GetElevatorID(from int, to int) (int, error)
}
