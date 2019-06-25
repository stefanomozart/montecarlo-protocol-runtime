// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"testing"
)

func TestIntMeanWithHE_Setup(t *testing.T) {
	sim := NewIntMeanWithHE(1024)

	tests := []struct {
		name       string
		simullator *IntMeanWithHE
		args       []int64
		wantErr    bool
	}{
		{
			"a test always helps",
			sim,
			[]int64{1, 2, 3},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.simullator.Setup(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntMeanWithHE.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntMeanWithHE_Run(t *testing.T) {
	cap := 200
	small := make([]int64, cap)
	for i := 0; i < cap; i++ {
		small[i] = 20
	}
	tests := []struct {
		name    string
		bitlen  int
		dataset []int64
		wantErr bool
	}{
		{
			"test with extremely small dataset",
			1024,
			[]int64{2, 2, 3, 3, 4, 4},
			false,
		},
		{
			"test with extremely small dataset",
			1024,
			[]int64{2, 2, 2, 2, 4, 4, 4, 4},
			false,
		},
		{
			"test with very small dataset",
			1024,
			small,
			false,
		},
		{
			"test with random datase",
			1024,
			randomSizeDataset(),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simullation := NewIntMeanWithHE(tt.bitlen)
			simullation.Setup(tt.dataset)
			// run
			if err := simullation.Run(); (err != nil) != tt.wantErr {
				t.Errorf("IntMeanWithHE.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			// test result
			var m int64
			for _, p := range tt.dataset {
				m += p
			}
			m = m / int64(len(tt.dataset))
			r := simullation.Result()
			if r <= 0 || r > (m+3) {
				t.Errorf("IntMeanWithHE.Run() got %v, want %v", r, m)
			}
		})
	}
}
