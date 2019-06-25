package montecarlo

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"

	"github.com/atgjack/prob"
)

func loadDataset(filename string, column int) []int64 {

	return []int64{}
}

func loadCSVDataset(filename string, column int) []int64 {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	c := csv.NewReader(f)

	records, err := c.ReadAll()
	if err != nil {
		panic(err)
	}
	dataset := make([]int64, len(records))
	for i, r := range records {
		dataset[i], err = strconv.ParseInt(r[column], 10, 64)
		if err != nil {
			dataset[i] = 0
		}
	}

	return dataset
}

// return a normally distributed dataset with size u ~ U(600, 16000)
func randomSizeDataset() []int64 {
	return getIntNormalSample(rand.Intn(1000)+600, 10, 3)
}

func getIntBetaSample(n, fact, alpha, beta int) []int64 {
	b, err := prob.NewBeta(float64(alpha), float64(beta))
	if err != nil {
		panic(err)
	}

	dataset := make([]int64, n)
	for i := 0; i < n; i++ {
		val := int64(float64(fact) * b.Random())
		if val < 1 {
			val = 1
		}
		dataset[i] = val
	}

	return dataset
}

func getIntGammaSample(n, fact, alpha, beta int) []int64 {
	g, err := prob.NewGamma(float64(alpha), float64(beta))
	if err != nil {
		panic(err)
	}

	dataset := make([]int64, n)
	for i := 0; i < n; i++ {
		val := int64(float64(fact) * g.Random())
		if val < 1 {
			val = 1
		}
		dataset[i] = val
	}

	return dataset
}

func getIntNormalSample(n int, mu, sigma int64) []int64 {
	dataset := make([]int64, n)
	for i := 0; i < n; i++ {
		val := (int64(rand.NormFloat64()) * sigma) + mu
		if val < 1 {
			val = 1
		}
		dataset[i] = val
	}

	return dataset
}

func getIntUniformSample(n int, min, max int64) []int64 {
	dataset := make([]int64, n)
	for i := 0; i < n; i++ {
		dataset[i] = rand.Int63n(max-min) + min
	}

	return dataset
}
