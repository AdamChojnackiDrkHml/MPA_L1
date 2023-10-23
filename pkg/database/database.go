package database

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Database struct {
	pricesMatrix    [][]float64
	amountMatrix    [][]int
	chanceOfAnomaly float64
	isPopulated     bool
	Products        map[string][]*ProductEntry
}

type ProductEntry struct {
	ProductName string
	Amount      int
	Price       float64
}

func Create_Database(
	pricesMatrix [][]float64,
	amountMatrix [][]int,
	chanceOfAnomaly float64,
) *Database {

	return &Database{
		pricesMatrix:    pricesMatrix,
		amountMatrix:    amountMatrix,
		chanceOfAnomaly: chanceOfAnomaly,
		isPopulated:     false,
		Products:        make(map[string][]*ProductEntry),
	}
}

func (db *Database) SaveToFile(fileName string) bool {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return false
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(db.String())
	writer.Flush()

	return err == nil
}

func Create_FromFile(fileName string) *Database {
	dbSpec, err := os.ReadFile(fileName)

	if err != nil {
		return nil
	}

	return databaseFromString(string(dbSpec))
}

func (db *Database) Populate(sampleSizes []int, seed uint64, setOfProductNames []string) {
	if db.isPopulated {
		fmt.Print("Database already populated")
		return
	}

	for i, productName := range setOfProductNames {
		prices := db.createPrices(sampleSizes[i], seed, i)
		amounts := db.createAmounts(sampleSizes[i], seed, i)

		db.Products[productName] = make([]*ProductEntry, 0)

		for j := 0; j < sampleSizes[i]; j++ {
			db.Products[productName] = append(
				db.Products[productName],
				&ProductEntry{
					ProductName: productName,
					Amount:      amounts[j],
					Price:       prices[j],
				})
		}
	}

	db.isPopulated = true
}

func (db *Database) createPrices(sampleSize int, seed uint64, productIndex int) []float64 {
	dist := distuv.Uniform{
		Min: db.pricesMatrix[productIndex][0],
		Max: db.pricesMatrix[productIndex][1],

		Src: rand.NewSource(seed),
	}

	// Generating random prices
	prices := make([]float64, sampleSize-2)

	for i := range prices {
		prices[i] = dist.Rand()
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i] < prices[j]
	})

	// Creating anomalies
	randomGenerator := rand.New(rand.NewSource(seed))

	for i := 0; i < len(prices); i++ {
		if randomGenerator.Float64() < db.chanceOfAnomaly {
			prices[i] = priceAnomaly(prices[i], randomGenerator)
		}
	}

	// Append start and end price as specified
	prices = prependFloat(prices, db.pricesMatrix[productIndex][0])
	prices = append(prices, db.pricesMatrix[productIndex][1])

	return prices
}

func priceAnomaly(price float64, randomGenerator *rand.Rand) float64 {
	randomOracle := randomGenerator.Float64()

	if randomOracle < 0.3 {
		return price * 100.0 // missing comma
	} else if randomOracle < 0.9 {
		return price * math.Pow10((randomGenerator.Intn(4) - 3)) // comma in wrong place
	} else {
		return 0 // didn't enter price
	}
}

func (db *Database) createAmounts(sampleSize int, seed uint64, productIndex int) []int {
	dist := distuv.Uniform{
		Min: float64(db.amountMatrix[productIndex][0]),
		Max: float64(db.amountMatrix[productIndex][1]),
		Src: rand.NewSource(seed),
	}

	// Generating random amounts
	amounts := make([]int, sampleSize)

	for i := range amounts {
		amounts[i] = int(dist.Rand())
	}

	// Creating anomalies
	randomGenerator := rand.New(rand.NewSource(seed))

	for i := 0; i < len(amounts); i++ {
		if randomGenerator.Float64() < db.chanceOfAnomaly {
			amounts[i] = amountAnomaly(amounts[i], randomGenerator)
		}
	}

	return amounts
}

func amountAnomaly(amount int, randomGenerator *rand.Rand) int {
	randomOracle := randomGenerator.Float64()

	if randomOracle < 0.9 {
		return amount * int(math.Pow10((randomGenerator.Intn(4) - 2))) // comma in wrong place
	} else {
		return 0 // didn't enter amount
	}
}

func prependFloat(x []float64, y float64) []float64 {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}

func (db *Database) String() string {

	productsString := make([]string, 0)

	for _, productTypeList := range db.Products {
		for _, pe := range productTypeList {
			productsString = append(productsString, pe.String())
		}
	}

	return strings.Join(productsString, "\n")
}

func databaseFromString(dbSpec string) *Database {
	dbSpecLines := strings.Split(dbSpec, "\n")

	db := &Database{
		isPopulated: true,
		Products:    make(map[string][]*ProductEntry),
	}

	for _, productSpec := range dbSpecLines {
		productEntry := productEntryFromString(productSpec)

		if db.Products[productEntry.ProductName] != nil {
			db.Products[productEntry.ProductName] = append(db.Products[productEntry.ProductName], productEntry)
		} else {
			db.Products[productEntry.ProductName] = []*ProductEntry{productEntry}
		}
	}

	return db
}

func (pe *ProductEntry) String() string {
	return fmt.Sprintf("%v %v %v", pe.ProductName, pe.Amount, pe.Price)
}

func productEntryFromString(productEntrySpecString string) *ProductEntry {
	productSpec := strings.Split(productEntrySpecString, " ")

	productName := productSpec[0]
	productAmount, _ := strconv.Atoi(productSpec[1])
	productPrice, _ := strconv.ParseFloat(productSpec[2], 64)

	return &ProductEntry{
		ProductName: productName,
		Amount:      productAmount,
		Price:       productPrice,
	}
}
