package main

import (
	"L1/pkg/database"
	"L1/pkg/predictor"
	"fmt"
	"time"
)

func main() {
	filepath := "./data/exampleDatabase"

	sampleSize := []int{1000, 1000}
	productNames := []string{"Parasol", "Lornetka"}

	prices := [][]float64{
		{100.0, 124.0},
		{50.0, 80.0},
	}

	amounts := [][]int{
		{30, 50},
		{100, 150},
	}

	seed := uint64(time.Now().UnixNano())

	db := database.Create_Database(prices, amounts, 0.003)
	db.Populate(sampleSize, seed, productNames)
	db.SaveToFile(filepath)

	p := predictor.Create_Predictor(db, 0.99)

	newValidProduct := &database.ProductEntry{
		ProductName: "Parasol",
		Price:       100,
		Amount:      40,
	}

	newInvalidPriceProduct := &database.ProductEntry{
		ProductName: "Parasol",
		Price:       1000,
		Amount:      40,
	}

	newInvalidAmountProduct := &database.ProductEntry{
		ProductName: "Parasol",
		Price:       100,
		Amount:      4,
	}

	products := []*database.ProductEntry{
		newValidProduct,
		newInvalidPriceProduct,
		newInvalidAmountProduct,
	}

	for _, product := range products {
		if res, err := p.InsertProductToDatabase(product); !res {
			fmt.Println(err)
		}
	}

	// db := database.Create_FromFile(filepath)
	// db.SaveToFile(filepath)
}
