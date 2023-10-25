package main

import (
	"L1/pkg/database"
	"time"
)

func main() {
	filepath := "./data/exampleGeneratedDatabase"

	sampleSize := 60
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

	db := database.Create_Database(sampleSize, seed, productNames, prices, amounts, 0.1)

	// db := database.Create_FromFile(filepath)
	db.SaveToFile(filepath)

}
