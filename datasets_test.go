// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"reflect"
	"testing"
)

func Test_loadDataset(t *testing.T) {
	tests := []struct {
		name, filename string
		colunm         int
		want           []int64
	}{
		{
			"teste tem que funcionar, senão, não tem graça",
			"teste.data",
			2,
			[]int64{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadDataset(tt.filename, tt.colunm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadDataset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadCSVDataset(t *testing.T) {
	type args struct {
		filename string
		column   int
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			"um teste qualquer",
			args{
				"teste.csv",
				0,
			},
			[]int64{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadCSVDataset(tt.args.filename, tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadCSVDataset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntNormalSample(t *testing.T) {
	type args struct {
		n      int
		mu     int64
		sigma2 int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"generate 1",
			args{100, 120, 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIntNormalSample(tt.args.n, tt.args.mu, tt.args.sigma2)
			if len(got) != tt.args.n {
				t.Errorf("getIntNormalSample() = %v, want %v", len(got), tt.args.n)
			}
			for _, val := range got {
				if val < 1 {
					t.Errorf("getIntNormalSample() off limits = %v", val)
				}
			}
		})
	}
}

func Test_getIntGammaSample(t *testing.T) {
	type args struct {
		n, fact, alpha, beta int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Testing just the limits",
			args{100, 120, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIntGammaSample(tt.args.n, tt.args.fact, tt.args.alpha, tt.args.beta)

			if len(got) != tt.args.n {
				t.Errorf("getIntGammaSample() = %v, want %v", len(got), tt.args.n)
			}
			for _, val := range got {
				if val < 1 {
					t.Errorf("getIntGammaSample() off limits = %v", val)
				}
			}
		})
	}
}

func Test_getIntBetaSample(t *testing.T) {
	type args struct {
		n, fact, alpha, beta int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Testing just the limits",
			args{
				100,
				130,
				30,
				2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIntBetaSample(tt.args.n, tt.args.fact, tt.args.alpha, tt.args.beta)
			if len(got) != tt.args.n {
				t.Errorf("getIntBetaSample() = %v, want %v", len(got), tt.args.n)
			}

			for _, val := range got {
				if val < 1 {
					t.Errorf("getIntBetaSample() off limits = %v", val)
				}
			}
		})
	}
}
