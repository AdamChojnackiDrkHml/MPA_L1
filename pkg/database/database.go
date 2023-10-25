package database

import (
	"bufio"
	"os"
	"strings"
)

type Database struct {
	Products map[string]*Product
}

func Create_Database(
	sampleSizes int,
	seed uint64,
	setOfProductNames []string,
	pricesMatrix [][]float64,
	amountMatrix [][]int,
	chanceOfAnomaly float64,
) *Database {

	db := &Database{Products: make(map[string]*Product)}

	for i, productName := range setOfProductNames {
		db.Products[productName] = Create_Product(sampleSizes, seed, pricesMatrix[i][0], pricesMatrix[i][1], amountMatrix[i][0], amountMatrix[i][1], chanceOfAnomaly, setOfProductNames[i])
	}

	return db
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

// func Create_FromFile(fileName string) *Database {
// 	dbSpec, err := os.ReadFile(fileName)

// 	if err != nil {
// 		return nil
// 	}

// 	return databaseFromString(string(dbSpec))
// }

func (db *Database) String() string {

	productsString := make([]string, 0)

	for _, product := range db.Products {
		productsString = append(productsString, product.String())
	}

	return strings.Join(productsString, "\n")
}

// func databaseFromString(dbSpec string) *Database {
// 	dbSpecLines := strings.Split(dbSpec, "\n")

// 	db := &Database{
// 		Products: make(map[string]*Product),
// 	}

// 	for _, productSpec := range dbSpecLines {
// 		productEntry, productName := productEntryFromString(productSpec)

// 		if db.Products[productName] != nil {
// 			db.Products[productName] = append(db.Products[productName], productEntry)
// 		} else {
// 			db.Products[productName] = []*ProductEntry{productEntry}
// 		}
// 	}

// 	return db
// }
