package main

import "github.com/stefanomozart/montecarlo"

func main() {
	for _, m := range []int{1000, 5000} {
		montecarlo.ExperimentType1(m)
	}
	for _, m := range []int{1000, 5000} {
		for _, n := range []int{200, 500, 1000} {
			montecarlo.ExperimentType2(m, n)
		}
	}
}
