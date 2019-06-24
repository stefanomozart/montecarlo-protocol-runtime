// Package montecarlo implements a few montecarlo simullations
package montecarlo

import "testing"

func Test_ExperimentType1(t *testing.T) {
	tests := []struct {
		name string
		m    int
	}{
		{
			"run the experimento once, just to check file the boilerplate code",
			1,
		},
		{
			"run the experimento once more, just to check file the boilerplate code",
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExperimentType1(tt.m)
		})
	}
}
