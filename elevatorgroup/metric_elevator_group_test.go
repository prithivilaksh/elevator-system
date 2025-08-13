package elevatorgroup_test

import (
	"math/rand"
	"testing"

	"github.com/prithivilaksh/elevator-system/elevator"
	"github.com/prithivilaksh/elevator-system/elevatorgroup"
)

func setupBenchmarkElevators(elevatorCount int) *elevatorgroup.MetricElevatorGroup {
	eg := elevatorgroup.NewMetricElevatorGroup()
	for i := range elevatorCount {
		elevator := elevator.NewLeastDisElevator(i+1)
		eg.AddElevator(elevator)
	}
	return eg
}

func benchmarkGetElevatorID(elevatorCount int, b *testing.B) {
	eg := setupBenchmarkElevators(elevatorCount)

	for b.Loop() {
		// Generate random floor numbers between 1 and 20
		from := rand.Intn(20) + 1
		to := rand.Intn(20) + 1
		// Ensure from and to are different
		for from == to {
			to = rand.Intn(20) + 1
		}
		eg.GetElevatorID(from, to)
	}
}

func BenchmarkGetElevatorID_1Elevator(b *testing.B)   { benchmarkGetElevatorID(1, b) }
func BenchmarkGetElevatorID_5Elevators(b *testing.B)  { benchmarkGetElevatorID(5, b) }
func BenchmarkGetElevatorID_10Elevators(b *testing.B) { benchmarkGetElevatorID(10, b) }
func BenchmarkGetElevatorID_50Elevators(b *testing.B) { benchmarkGetElevatorID(50, b) }
