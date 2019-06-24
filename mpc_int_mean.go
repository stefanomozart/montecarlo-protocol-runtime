// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"time"

	"github.com/stefanomozart/mpc"
)

// IntMeanWithMPC is the computation of a integer mean
// using homomorphic encryption
type IntMeanWithMPC struct {
	bitlen     int
	mean       int64
	dim        mpc.IntProtocol
	setup, run time.Duration
}

// NewIntMeanWithMPC returns a new instance of the simulation
func NewIntMeanWithMPC(bitlen int) *IntMeanWithMPC {
	return &IntMeanWithMPC{
		bitlen: bitlen,
	}
}

// Setup the data needed for simullation
func (im *IntMeanWithMPC) Setup(dataset []int64) error {
	start := time.Now()

	params := mpc.NewParameters(im.bitlen)
	im.dim = mpc.NewDistributedIntMean()
	if err := im.dim.Setup(params, dataset); err != nil {
		return err
	}

	im.setup = time.Since(start)
	return nil
}

// Run the mean calculation simmulation
func (im *IntMeanWithMPC) Run() error {
	start := time.Now()

	err := im.dim.Run()
	if err != nil {
		return err
	}

	im.mean = im.dim.Output()
	im.run = time.Since(start)
	return nil
}

// Result returns the computed mean
func (im *IntMeanWithMPC) Result() int64 {
	return im.mean
}

// Runtimes for the setup and the online phase of the simullation
func (im *IntMeanWithMPC) Runtimes() (int64, int64) {
	return im.setup.Nanoseconds(), im.run.Nanoseconds()
}
