package database

import (
	"fmt"
	"strings"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type Product struct {
	Name    string
	Amounts []int
	Prices  []float64
}

func Create_Product(
	sampleSize int,
	seed uint64,
	minPrice float64,
	maxPrice float64,
	minAmount int,
	maxAmount int,
	chanceOfAnomaly float64,
	name string,
) *Product {
	return &Product{
		Amounts: createAmounts(sampleSize, seed, minAmount, maxAmount, chanceOfAnomaly),
		Prices:  createPrices(sampleSize, seed, minPrice, maxPrice, chanceOfAnomaly),
		Name:    name,
	}
}

func createPrices(sampleSize int, seed uint64, minPrice, maxPrice float64, chanceOfAnomaly float64) []float64 {

	dist := distuv.Uniform{
		Min: -0.2,
		Max: 0.2,
		Src: rand.NewSource(seed),
	}

	randomGenerator := rand.New(rand.NewSource(seed))

	currPrice := minPrice
	priceDiff := maxPrice - minPrice
	priceDelta := priceDiff / float64(sampleSize)

	// Generating random prices
	prices := make([]float64, sampleSize)

	for i := 0; i < sampleSize; i++ {
		prices[i] = currPrice + (dist.Rand() * priceDiff)

		if randomGenerator.Float64() < chanceOfAnomaly {
			prices[i] = createAnomaly(prices[i], randomGenerator)
		}

		currPrice += priceDelta
	}

	// Append start and end price as specified
	return prices
}

func createAnomaly(price float64, randomGenerator *rand.Rand) float64 {
	possibleAnomalies := []float64{0.01, 0.1, 10, 100}
	randomOracle := randomGenerator.Intn(len(possibleAnomalies))

	return price * possibleAnomalies[randomOracle]
}

func createAmounts(sampleSize int, seed uint64, minAmount, maxAmount int, chanceOfAnomaly float64) []int {

	randomGenerator := rand.New(rand.NewSource(seed))

	// Generating random amounts
	amounts := make([]int, sampleSize)

	for i := range amounts {
		amounts[i] = rand.Intn(maxAmount-minAmount) + minAmount

		if randomGenerator.Float64() < chanceOfAnomaly {
			amounts[i] = int(createAnomaly(float64(amounts[i]), randomGenerator))
		}
	}

	return amounts
}

func (p *Product) CreateDataPoints() ([][]float64, [][]float64) {
	dataPointsPrice := make([][]float64, len(p.Amounts))
	dataPointsAmounts := make([][]float64, len(p.Amounts))

	for i := range p.Amounts {
		dataPointsPrice[i] = []float64{float64(i), p.Prices[i]}
		dataPointsAmounts[i] = []float64{float64(i), float64(p.Amounts[i])}
	}

	return dataPointsPrice, dataPointsAmounts
}

func (p *Product) String() string {
	returnString := make([]string, len(p.Amounts))

	for i := range p.Amounts {
		returnString[i] = fmt.Sprintf("%v %v %v", p.Name, p.Amounts[i], p.Prices[i])
	}

	return strings.Join(returnString, "\n")
}

// func ProductEntryFromString(productEntrySpecString string) (*ProductEntry, string) {
// 	productSpec := strings.Split(productEntrySpecString, " ")

// 	productName := productSpec[0]
// 	productAmount, _ := strconv.Atoi(productSpec[1])
// 	productPrice, _ := strconv.ParseFloat(productSpec[2], 64)

// 	return &ProductEntry{
// 			Amount: productAmount,
// 			Price:  productPrice,
// 		},
// 		productName
// }

func (p *Product) ReplaceValues(clearedPrices [][]float64, clearedAmounts [][]float64) {
	p.Prices = make([]float64, len(clearedPrices))
	for i, elem := range clearedPrices {
		p.Prices[i] = elem[1]
	}

	p.Amounts = make([]int, len(clearedAmounts))
	for i, elem := range clearedAmounts {
		p.Amounts[i] = int(elem[1])
	}

}
