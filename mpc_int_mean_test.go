// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"testing"
)

func TestIntMeanWithMPC_Setup(t *testing.T) {
	tests := []struct {
		name    string
		bitlen  int
		dataset []int64
		wantErr bool
	}{
		{
			"test, just to be shure setup happens",
			1024,
			[]int64{1, 2, 3},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntMeanWithMPC(tt.bitlen)
			err := im.Setup(tt.dataset)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntMeanWithMPC.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntMeanWithMPC_Run(t *testing.T) {
	cap := 100
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
			[]int64{2, 2, 4, 4},
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
			simulation := NewIntMeanWithMPC(tt.bitlen)
			if err := simulation.Setup(tt.dataset); (err != nil) != tt.wantErr {
				t.Errorf("IntMeanWithHE.Run()::Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			// run
			if err := simulation.Run(); (err != nil) != tt.wantErr {
				t.Errorf("IntMeanWithHE.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			// test result
			var m int64
			for _, p := range tt.dataset {
				m += p
			}
			m = m / int64(len(tt.dataset))
			r := simulation.Result()
			if r <= 0 || r > (m+3) {
				t.Errorf("IntMeanWithHE.Run() got %v, want %v", r, m)
			}
		})
	}
}
