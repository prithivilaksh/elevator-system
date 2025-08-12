package elevator

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLeastDisElevator(t *testing.T) {
	elevator := NewLeastDisElevator(1)
	assert.Equal(t, 1, elevator.GetID())
	assert.Equal(t, 0, elevator.currFloor)
	assert.True(t, elevator.isIdle)
	assert.Empty(t, elevator.stops)
}

func TestGetID(t *testing.T) {
	elevator := &LeastDisElevator{id: 42}
	assert.Equal(t, 42, elevator.GetID())
}

func TestAddStopAndGetNextInd(t *testing.T) {
	tests := []struct {
		name         string
		startInd     int
		afterFloor   int
		stop         int
		stops        []int
		expected     []int
		expectedNext int
	}{
		{
			"Insert in middle - ascending",
			0, 2, 4, []int{2, 6, 8},
			[]int{2, 4, 6, 8}, 2,
		},
		{
			"Insert in middle - descending",
			0, 8, 4, []int{8, 6, 2},
			[]int{8, 6, 4, 2}, 3,
		},
		{
			"Insert at end",
			0, 1, 10, []int{2, 4, 6},
			[]int{2, 4, 6, 10}, 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, nextInd := addStopAndGetNextInd(tt.startInd, tt.afterFloor, tt.stop, tt.stops)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedNext, nextInd)
		})
	}
}

func TestAddStopsAndFindDistance(t *testing.T) {
	tests := []struct {
		name     string
		from     int
		to       int
		currFloor int
		stops    []int
		expectedStops []int
		expectedDist  int
	}{
		{
			"Simple up",
			3, 5, 1,
			[]int{2, 4, 6},
			[]int{2, 3, 4, 5, 6},
			4, // 1->2->3->4->5 (stops at 5)
		},
		{
			"Simple down",
			5, 3, 7,
			[]int{6, 4, 2},
			[]int{6, 5, 4, 3, 2},
			4, // 7->6->5->4->3 (stops at 3)
		},
		{
			"Already existing from",
			3, 5, 1,
			[]int{2,3,7},
			[]int{2, 3, 5, 7},
			4, // 7->6->5->4->3 (stops at 3)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stops, dist := addStopsAndFindDistance(tt.from, tt.to, tt.currFloor, tt.stops)
			assert.Equal(t, tt.expectedStops, stops)
			assert.Equal(t, tt.expectedDist, dist)
		})
	}
}

func TestGetMetric(t *testing.T) {
	elevator := NewLeastDisElevator(1)
	elevator.currFloor = 5
	elevator.stops = []int{3, 7}

	// Test getting metric for a new request
	distance := elevator.GetMetric(2, 4)
	// Expected: 5->3->7->2->4 = 13
	assert.Equal(t, 13, distance)
}

func TestTryAddStops(t *testing.T) {
	elevator := NewLeastDisElevator(1)
	elevator.currFloor = 1

	// Test adding a valid stop
	added := elevator.TryAddStops(3, 5, 6, 2)
	assert.True(t, added)
	assert.Equal(t, []int{3, 5}, elevator.stops)
	assert.False(t, elevator.isIdle)

	// Wait for simulation to complete
	time.Sleep(100 * time.Millisecond)

	// Test adding an invalid stop (metric too different)
	added = elevator.TryAddStops(10, 12, 5, 2)
	assert.False(t, added)
}

func TestConcurrentAccess(t *testing.T) {
	elevator := NewLeastDisElevator(1)
	var wg sync.WaitGroup

	// Start multiple goroutines trying to add stops
	for i := range 5 {
		wg.Add(1)
		go func(floor int) {
			defer wg.Done()
			elevator.TryAddStops(floor, floor+2, 2, 1)
		}(i * 2)
	}

	wg.Wait()
	// Just verify no race conditions occurred
	assert.True(t, len(elevator.stops) > 0)
}