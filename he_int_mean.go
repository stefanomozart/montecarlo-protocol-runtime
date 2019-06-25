// Package montecarlo implements a few montecarlo simullations
package montecarlo

import (
	"math/big"
	"sync"
	"time"

	"github.com/stefanomozart/paillier"
)

// IntMeanWithHE is the computation of a integer mean
// using homomorphic encryption
type IntMeanWithHE struct {
	bitlen     int
	pk         *paillier.PublicKey
	sk         *paillier.PrivateKey
	srv        chan []*big.Int
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

	// generate key pair
	var err error
	im.pk, im.sk, err = paillier.GenerateKeyPair(im.bitlen)
	if err != nil {
		return err
	}

	// cipher the dataset
	im.cts = make([]*big.Int, len(args))
	for i, m := range args {
		im.cts[i], err = im.pk.Encrypt(m)
		if err != nil {
			return err
		}
	}

	// Instantiate server
	im.srv = make(chan []*big.Int)
	go performHEParalellComputation(im.pk, im.srv)

	im.setup = time.Since(start)
	return nil
}

func performHEComputation(pk *paillier.PublicKey, ch chan []*big.Int) {
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

func performHEParalellComputation(pk *paillier.PublicKey, cli chan []*big.Int) {
	// read channel from client
	cts := <-cli

	// internallly, the server may implement several optimizations,
	// including parallel execution. We are going to split computation
	// onto three threads
	k := 3

	// instantiate k subprocesses
	var parts []chan []*big.Int
	wg := new(sync.WaitGroup)
	for i := 0; i < k; i++ {
		parts = append(parts, make(chan []*big.Int))
		wg.Add(1)
		go bachAdd(pk, parts[i], wg, i)
	}

	// split the dataset and send to the k parts
	batchSize := (len(cts) + k - 1) / k

	for i, j := 0, 0; i < len(cts); i, j = i+batchSize, j+1 {
		end := i + batchSize

		if end > len(cts) {
			end = len(cts)
		}

		parts[j] <- cts[i:end]
	}

	// receive k partials, and sum
	btSum := make([]*big.Int, k)
	for i, ch := range parts {
		res := <-ch
		btSum[i] = res[0]
	}
	wg.Wait()

	// Add partials
	ctSum := pk.BatchAdd(btSum...)

	// perform homomorphic multiplication (sum * (1/n))
	ctMean, err := pk.DivPlaintext(ctSum, int64(len(cts)))
	if err != nil {
		cli <- []*big.Int{new(big.Int), new(big.Int)}
	}

	// send homomorphic result to the client
	cli <- []*big.Int{ctMean, ctSum}
}

func bachAdd(pk *paillier.PublicKey, part chan []*big.Int, wg *sync.WaitGroup, w int) {
	// read from channel
	cts := <-part

	// perform homomorphic addition
	ctSum := pk.BatchAdd(cts...)

	// write result to channel
	part <- []*big.Int{ctSum}
	wg.Done()
}

// Run the mean calculation simmulation
func (im *IntMeanWithHE) Run() error {
	start := time.Now()

	// send ciphertexts to server, for homomorphic computations
	im.srv <- im.cts

	// read response from server (this is a blocking operation)
	ctMean := <-im.srv

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
