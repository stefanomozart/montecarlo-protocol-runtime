// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"fmt"
	"os"
)

// Simullation is a runnable
type Simullation interface {
	Setup(args ...interface{})
	Run() error
	TimeElapsed() int64
}

// Run `sim` (the given Simullation function) `m` times, with the given `args`
func Run(sim Simullation, m int, args ...interface{}) {
	sim.Setup(args)
	times := make([]int64, m)
	for i := 0; i < m; i++ {
		sim.Run()
		times[i] = sim.TimeElapsed()
	}
}

// ExperimentType1 simply runs the two protocols with a fixed dataset,
// for `m` iterations
func ExperimentType1(m int) {
	//dataset := loadDataset("./datasets/wine.csv", 12) // Proline
	dataset := loadCSVDataset("./datasets/bank.csv", 11) // balance: average yearly balance, in euros (numeric)
	//dataset := loadCSVDataset("./datasets/dow_jones_index.csv", 7) // volume: the number of shares of stock that traded hands in the week

	// determines the bitsize of the public key used in homomorphic
	// encryption of in secret particioning in the MPC protocol
	securityParameter := 1024

	// Repeat experiment HE integer mean `m` times
	fHE, _ := os.Create(fmt.Sprintf("./datasets/exp1he%v.bank.csv", m))
	defer fHE.Close()
	for i := 0; i < m; i++ {
		imHE := NewIntMeanWithHE(securityParameter)
		imHE.Setup(dataset)
		imHE.Run()
		imHEres := imHE.Result()
		imHEclient, imHEserver := imHE.Runtimes()
		fHE.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", imHEres, imHEclient, imHEserver, imHEclient+imHEserver))
	}

	// Repeat experiment Secret-sharing integer mean for `m` times
	fMPC, _ := os.Create(fmt.Sprintf("./datasets/exp1mpc%v.bank.csv", m))
	defer fMPC.Close()
	for i := 0; i < m; i++ {
		imMPC := NewIntMeanWithMPC(securityParameter)
		imMPC.Setup(dataset)
		imMPC.Run()
		imMPCres := imMPC.Result()
		imMPCclient, imMPCserver := imMPC.Runtimes()
		fMPC.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", imMPCres, imMPCclient, imMPCserver, imMPCclient+imMPCserver))
	}
}

// ExperimentType1: run the two protocols with a dataset sampled at each
// of the `m` iteraations
func ExperimentType2(m int) {
	for i := 0; i < m; i++ {
		// sample a dataset size

		// sample the data set of the size defined in the previous step

		// run both algorithms on the sampled dataset

		// record runtimes
	}
}
