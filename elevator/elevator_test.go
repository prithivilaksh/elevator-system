package elevator

import "testing"

type FindDisTestCase struct {
	name      string
	elevator  *elevator
	stops     []int
	currFloor int
	lastInd   int
	want      int
}

func TestFindDistance(t *testing.T) {

	ele1 := NewElevator("Elevator-1", 10)
	currFloor1 := 3
	stops1 := []int{2, 6, 4, 9, 10}
	want1 := 13

	ele2 := NewElevator("Elevator-2", 30)
	currFloor2 := 5
	stops2 := []int{2, 22, 9, 10, 25}
	want2 := 52

	cases := []FindDisTestCase{
		{"Find Distance case 1", ele1, stops1, currFloor1, len(stops1), want1},
		{"Find Distance case 2", ele2, stops2, currFloor2, len(stops2), want2},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			dis := tt.elevator.FindDistance(tt.stops, tt.currFloor, tt.lastInd)
			if dis != tt.want {
				t.Errorf("FindDistance() = %v, want %v", dis, tt.want)
			}
		})
	}
}
