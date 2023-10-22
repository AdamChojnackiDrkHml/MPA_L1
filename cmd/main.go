package main

import (
	"L1/pkg/database"
)

func main() {
	filepath := "./data/exampleDatabase"

	// sampleSize := []int{1000}
	// productNames := []string{"Parasol"}
	// seed := uint64(time.Now().UnixNano())
	// db := database.Create_Database(100.0, 124.0, 30, 50, 0.03)
	// db.Populate(sampleSize, seed, productNames)

	db := database.Create_FromFile(filepath)
	db.SaveToFile(filepath)
}
