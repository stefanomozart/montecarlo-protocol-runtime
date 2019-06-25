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
	datasets := map[string]int{
		//"./datasets/wine.csv": 12, // Proline
		"./datasets/bank.csv":            11, // balance: average yearly balance, in euros (numeric)
		"./datasets/dow_jones_index.csv": 7,  // volume: the number of shares of stock that traded hands in the week
	}
	for filename, column := range datasets {
		dataset := loadCSVDataset(filename, column)

		// determines the bitsize of the public key used in homomorphic
		// encryption of in secret particioning in the MPC protocol
		securityParameter := 1024

		// Repeat experiment HE integer mean `m` times
		fHE, _ := os.Create(fmt.Sprintf("./datasets/exp1_he_m%v.%v.csv", m, column))
		defer fHE.Close()
		for i := 0; i < m; i++ {
			he := NewIntMeanWithHE(securityParameter)
			he.Setup(dataset)
			he.Run()
			heMu := he.Result()
			heTcli, heTsrv := he.Runtimes()
			fHE.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", heMu, heTcli, heTsrv, heTcli+heTsrv))
		}

		// Repeat experiment Secret-sharing integer mean for `m` times
		fMPC, _ := os.Create(fmt.Sprintf("./datasets/exp1_mpc_m%v.%v.csv", m, column))
		defer fMPC.Close()
		for i := 0; i < m; i++ {
			mpc := NewIntMeanWithMPC(securityParameter)
			mpc.Setup(dataset)
			mpc.Run()
			mpcMu := mpc.Result()
			mpcTcli, mpcTsrv := mpc.Runtimes()
			fMPC.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", mpcMu, mpcTcli, mpcTsrv, mpcTcli+mpcTsrv))
		}
	}
}

// ExperimentType2 runs the two protocols with a dataset sampled at each
// of the `m` iteractions
func ExperimentType2(m, n int) {
	// File to record HE protocol runtimes
	fHE, _ := os.Create(fmt.Sprintf("./datasets/exp2he_m%v_n%v.csv", m, n))
	defer fHE.Close()

	// File to record MPC protocol runtimes
	fMPC, _ := os.Create(fmt.Sprintf("./datasets/exp2mpc_m%v_n%v.csv", m, n))
	defer fMPC.Close()

	for i := 0; i < m; i++ {
		// Unif
		experimentType2ProtocolStep(fHE, fMPC, getIntUniformSample(n, 80, 160), "unif")

		// Normal
		experimentType2ProtocolStep(fHE, fMPC, getIntNormalSample(n, 120, 30), "norm")

		// Gamma
		experimentType2ProtocolStep(fHE, fMPC, getIntGammaSample(n, 120, 2, 2), "gamma")

		// Betta
		experimentType2ProtocolStep(fHE, fMPC, getIntBetaSample(n, 130, 30, 2), "beta")
	}
}

// experimentType2_ProtocolStep runs the two protocols with the given dataset
// and record the times
func experimentType2ProtocolStep(fHE, fMPC *os.File, dataset []int64, dist string) {
	// determines the bitsize of the public key used in homomorphic
	// encryption of in secret particioning in the MPC protocol
	securityParameter := 1024

	// run and record the HE protocol
	he := NewIntMeanWithHE(securityParameter)
	he.Setup(dataset)
	he.Run()
	heMu := he.Result()
	heTcli, heTsrv := he.Runtimes()
	fHE.WriteString(fmt.Sprintf("%s,%d,%d,%d,%d\n", dist, heMu, heTcli, heTsrv, heTcli+heTsrv))

	// run and record the MPC protocol
	mpc := NewIntMeanWithMPC(securityParameter)
	mpc.Setup(dataset)
	mpc.Run()
	mpcMu := mpc.Result()
	mpcTcli, mpcTsrv := mpc.Runtimes()
	fMPC.WriteString(fmt.Sprintf("%s,%d,%d,%d,%d\n", dist, mpcMu, mpcTcli, mpcTsrv, mpcTcli+mpcTsrv))
}
