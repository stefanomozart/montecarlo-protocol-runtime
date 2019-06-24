// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"fmt"
	"math/rand"
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
	/*fHE, _ := os.Create(fmt.Sprintf("./datasets/exp1he%v.bank.csv", m))
	defer fHE.Close()
	for i := 0; i < m; i++ {
		imHE := NewIntMeanWithHE(securityParameter)
		imHE.Setup(dataset)
		imHE.Run()
		imHEres := imHE.Result()
		imHEclient, imHEserver := imHE.Runtimes()
		fHE.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", imHEres, imHEclient, imHEserver, imHEclient+imHEserver))
	}*/

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

// ExperimentType2 runs the two protocols with a dataset sampled at each
// of the `m` iteractions
func ExperimentType2(m, n int) {
	// determines the bitsize of the public key used in homomorphic
	// encryption of in secret particioning in the MPC protocol
	securityParameter := 1024

	// File to record HE protocol runtimes
	fHE, _ := os.Create(fmt.Sprintf("./datasets/exp2he_m%v.csv", m))
	defer fHE.Close()

	// File to record MPC protocol runtimes
	fMPC, _ := os.Create(fmt.Sprintf("./datasets/exp2mpc_m%v.csv", m))
	defer fMPC.Close()

	for i := 0; i < m; i++ {
		// sample a dataset size
		s := rand.Intn(n-100) + 100

		// sample the data set of the size defined in the previous step
		dataset := getIntNormalSample(s, 100, 20)

		// run and record the HE protocol
		imHE := NewIntMeanWithHE(securityParameter)
		imHE.Setup(dataset)
		imHE.Run()
		imHEres := imHE.Result()
		imHEcli, imHEsrv := imHE.Runtimes()
		fHE.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%d\n", n, s, imHEres, imHEcli, imHEsrv, imHEcli+imHEsrv))

		// run and record the MPC protocol
		imMPC := NewIntMeanWithMPC(securityParameter)
		imMPC.Setup(dataset)
		imMPC.Run()
		imMPCres := imMPC.Result()
		imMPCcli, imMPCsrv := imMPC.Runtimes()
		fMPC.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%d\n", n, s, imMPCres, imMPCcli, imMPCsrv, imMPCcli+imMPCsrv))
	}
}

// ExperimentType3 runs the two protocols with a dataset sampled at each
// of the `m` iteractions
func ExperimentType3(m, n int) {
	// determines the bitsize of the public key used in homomorphic
	// encryption of in secret particioning in the MPC protocol
	securityParameter := 1024

	// File to record HE protocol runtimes
	fHE, _ := os.Create(fmt.Sprintf("./datasets/exp3he_m%v.csv", m))
	defer fHE.Close()

	// File to record MPC protocol runtimes
	fMPC, _ := os.Create(fmt.Sprintf("./datasets/exp3mpc_m%v.csv", m))
	defer fMPC.Close()

	for i := 0; i < m; i++ {
		// sample the data set of the size defined in the previous step
		dataset := getIntNormalSample(n, 100, 20)

		// run and record the HE protocol
		imHE := NewIntMeanWithHE(securityParameter)
		imHE.Setup(dataset)
		imHE.Run()
		imHEres := imHE.Result()
		imHEcli, imHEsrv := imHE.Runtimes()
		fHE.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d\n", n, imHEres, imHEcli, imHEsrv, imHEcli+imHEsrv))

		// run and record the MPC protocol
		imMPC := NewIntMeanWithMPC(securityParameter)
		imMPC.Setup(dataset)
		imMPC.Run()
		imMPCres := imMPC.Result()
		imMPCcli, imMPCsrv := imMPC.Runtimes()
		fMPC.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d\n", n, imMPCres, imMPCcli, imMPCsrv, imMPCcli+imMPCsrv))
	}
}
