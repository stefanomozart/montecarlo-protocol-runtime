// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"math/big"
	"time"

	"github.com/stefanomozart/paillier"
)

// IntMeanWithHE is the computation of a integer mean
// using homomorphic encryption
type IntMeanWithHE struct {
	bitlen     int
	pk         *paillier.PublicKey
	sk         *paillier.PrivateKey
	ch         chan []*big.Int
	cts        []*big.Int    // ciphertexts
	mean       int64         // computed mean
	setup, run time.Duration // runtimes for client and server
}

// NewIntMeanWithHE returns a new instance of the simulation
func NewIntMeanWithHE(bitlen int) *IntMeanWithHE {
	return &IntMeanWithHE{
		bitlen: bitlen,
	}
}

// Setup the data needed for simullation
func (im *IntMeanWithHE) Setup(args []int64) error {
	start := time.Now()

	var err error
	im.pk, im.sk, err = paillier.GenerateKeyPair(im.bitlen)
	if err != nil {
		return err
	}

	im.cts = make([]*big.Int, len(args))
	for i, m := range args {
		im.cts[i], err = im.pk.Encrypt(m)
		if err != nil {
			return err
		}
	}

	im.ch = make(chan []*big.Int)
	go performHEComputation(im.pk, im.ch, im.sk)

	im.setup = time.Since(start)
	return nil
}

func performHEComputation(pk *paillier.PublicKey, ch chan []*big.Int, sk *paillier.PrivateKey) {
	// read channel from client processes
	cts := <-ch

	// perform homomorphic addition
	ctSum := pk.BatchAdd(cts...)

	// perform homomorphic multiplication (sum * (1/n))
	ctMean, err := pk.DivPlaintext(ctSum, int64(len(cts)))
	if err != nil {
		ch <- []*big.Int{new(big.Int), new(big.Int)}
	}

	// send homomorphic result to the client
	ch <- []*big.Int{ctMean, ctSum}
}

// Run the mean calculation simmulation
func (im *IntMeanWithHE) Run() error {
	start := time.Now()

	// send ciphertexts to server, for homomorphic computations
	im.ch <- im.cts

	// read response from server
	ctMean := <-im.ch

	var err error
	im.mean, err = im.sk.Decrypt(ctMean[0])
	if err != nil {
		return err
	}
	if im.mean < 0 || im.mean > 100 {
		sum, _ := im.sk.Decrypt(ctMean[1])
		im.mean = sum / int64(len(im.cts))
	}

	im.run = time.Since(start)
	return nil
}

// Result of the delegated computation
func (im *IntMeanWithHE) Result() int64 {
	return im.mean
}

// Runtimes for client and server side
func (im *IntMeanWithHE) Runtimes() (int64, int64) {
	return im.setup.Nanoseconds(), im.run.Nanoseconds()
}
