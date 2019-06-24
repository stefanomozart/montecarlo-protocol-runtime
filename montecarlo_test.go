// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"testing"
)

func Test_ExperimentType1(t *testing.T) {
	tests := []struct {
		name string
		m    int
	}{
		{
			"run the experiment once, just to check file the boilerplate code",
			1,
		},
		{
			"run the experimento 5 more times, just to check file the boilerplate code",
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExperimentType1(tt.m)
		})
	}
}

func TestExperimentType2(t *testing.T) {
	tests := []struct {
		name string
		m    int //number of repetitions
		n    int // average size of the dataset
	}{
		{
			"run the experiment once, just to check file the boilerplate code",
			1,
			150,
		},
		{
			"run the experiment 5 more times, just to check file the boilerplate code",
			5,
			150,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExperimentType2(tt.m, tt.n)
		})
	}
}

func TestExperimentType3(t *testing.T) {
	tests := []struct {
		name string
		m    int //number of repetitions
		n    int // fixed size of the dataset
	}{
		{
			"run the experiment once, just to check file the boilerplate code",
			1,
			150,
		},
		{
			"run the experiment 5 more times, just to check file the boilerplate code",
			5,
			150,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExperimentType3(tt.m, tt.n)
		})
	}
}
