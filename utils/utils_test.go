package utils

import (
	"reflect"
	"testing"
)

type AbsTestCase[T number] struct {
	name  string
	input T
	want  T
}

func TestAbs(t *testing.T) {
	t.Run("with int", func(t *testing.T) {
		tests := []AbsTestCase[int]{
			{"positive number", 1, 1},
			{"negative number", -1, 1},
			{"zero", 0, 0},
		}
		runTests(t, tests)
	})

	t.Run("with float64", func(t *testing.T) {
		tests := []AbsTestCase[float64]{
			{"positive float", 3.5, 3.5},
			{"negative float", -2.75, 2.75},
			{"float zero", 0.0, 0.0},
		}
		runTests(t, tests)
	})

	t.Run("with int64", func(t *testing.T) {
		tests := []AbsTestCase[int64]{
			{"positive int64", 100, 100},
			{"negative int64", -200, 200},
		}
		runTests(t, tests)
	})
}

func runTests[T number](t *testing.T, tests []AbsTestCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.input); got != tt.want {
				t.Errorf("Abs(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

type DeepCopyTestCase struct {
	name  string
	input []int
	want  []int
}

func TestDeepCopy(t *testing.T) {

	same := []int{1, 2}
	tests := []DeepCopyTestCase{
		{"empty slice", []int{}, []int{}},
		{"slice with one element", []int{1}, []int{1}},
		{"slice with multiple elements", []int{1, 2, 3}, []int{1, 2, 3}},
		{"slice with same elements", same, same},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeepCopy(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepCopy(%v) = %v, want %v", tt.input, got, tt.want)
			}
			if &tt.input == &got {
				t.Errorf("got and want are same(Address)")
			}
		})
	}
}
