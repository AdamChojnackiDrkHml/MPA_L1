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
	minSamplePrice  float64
	maxSamplePrice  float64
	minSampleAmount int
	maxSampleAmount int
	chanceOfAnomaly float64
	isPopulated     bool
	Products        []*ProductEntry
}

type ProductEntry struct {
	ProductName string
	Amount      int
	Price       float64
}

func Create_Database(
	minSamplePrice float64,
	maxSamplePrice float64,
	minSampleAmount int,
	maxSampleAmount int,
	chanceOfAnomaly float64,
) *Database {

	return &Database{
		minSamplePrice:  minSamplePrice,
		maxSamplePrice:  maxSamplePrice,
		minSampleAmount: minSampleAmount,
		maxSampleAmount: maxSampleAmount,
		chanceOfAnomaly: chanceOfAnomaly,
		isPopulated:     false,
		Products:        make([]*ProductEntry, 0),
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
		prices := db.createPrices(sampleSizes[i], seed)
		amounts := db.createAmounts(sampleSizes[i], seed)

		for j := 0; j < sampleSizes[i]; j++ {
			db.Products = append(
				db.Products,
				&ProductEntry{
					ProductName: productName,
					Amount:      amounts[j],
					Price:       prices[j],
				})
		}
	}

	db.isPopulated = true
}

func (db *Database) createPrices(sampleSize int, seed uint64) []float64 {
	dist := distuv.Uniform{
		Min: db.minSamplePrice,
		Max: db.maxSamplePrice,

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
	prices = prependFloat(prices, db.minSamplePrice)
	prices = append(prices, db.maxSamplePrice)

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

func (db *Database) createAmounts(sampleSize int, seed uint64) []int {
	dist := distuv.Uniform{
		Min: float64(db.minSampleAmount),
		Max: float64(db.maxSampleAmount),
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

	retString := fmt.Sprintf("%v %v %v %v %v", db.minSamplePrice, db.maxSamplePrice, db.minSampleAmount, db.maxSampleAmount, db.chanceOfAnomaly)

	productsString := make([]string, 0)

	for _, pe := range db.Products {
		productsString = append(productsString, pe.String())
	}

	return retString + "\n" + strings.Join(productsString, "\n")
}

func databaseFromString(dbSpec string) *Database {
	dbSpecLines := strings.Split(dbSpec, "\n")
	header := strings.Split(dbSpecLines[0], " ")
	dbSpecLines = dbSpecLines[1:]

	minSamplePrice, _ := strconv.ParseFloat(header[0], 64)
	maxSamplePrice, _ := strconv.ParseFloat(header[1], 64)
	minSampleAmount, _ := strconv.Atoi(header[2])
	maxSampleAmount, _ := strconv.Atoi(header[3])
	chanceOfAnomaly, _ := strconv.ParseFloat(header[4], 64)

	db := Create_Database(
		minSamplePrice,
		maxSamplePrice,
		minSampleAmount,
		maxSampleAmount,
		chanceOfAnomaly,
	)

	for _, productSpec := range dbSpecLines {
		db.Products = append(db.Products, productEntryFromString(productSpec))
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
