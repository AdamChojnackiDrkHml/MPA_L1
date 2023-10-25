package predictor

import (
	"L1/pkg/database"
)

type Predictor struct {
	database                *database.Database
	sensitivity             float64
	expectedValuesOfPrices  map[string]float64
	expectedValuesOfAmounts map[string]float64
	varianceOfPrices        map[string]float64
	varianceOfAmounts       map[string]float64
}

// func Create_Predictor(database *database.Database, sensitivity float64) *Predictor {
// 	predictor := &Predictor{
// 		database:                database,
// 		sensitivity:             sensitivity,
// 		expectedValuesOfPrices:  make(map[string]float64),
// 		expectedValuesOfAmounts: make(map[string]float64),
// 		varianceOfPrices:        make(map[string]float64),
// 		varianceOfAmounts:       make(map[string]float64),
// 	}

// 	for productName, productTypeList := range database.Products {

// 		sumOfPrices := 0.0
// 		sumOfAmounts := 0
// 		numberOfProductsF := float64(len(productTypeList))

// 		for _, productEntry := range productTypeList {
// 			sumOfPrices += productEntry.Price
// 			sumOfAmounts += productEntry.Amount
// 		}

// 		predictor.expectedValuesOfPrices[productName] = sumOfPrices / numberOfProductsF
// 		predictor.expectedValuesOfAmounts[productName] = (float64(sumOfAmounts) / numberOfProductsF)

// 		sumOfSquaredDiffFromMeanPrices := 0.0
// 		sumOfSquaredDiffFromMeanAmounts := 0.0

// 		for _, productEntry := range productTypeList {
// 			sumOfSquaredDiffFromMeanPrices += math.Pow(productEntry.Price-predictor.expectedValuesOfPrices[productName], 2.0)
// 			sumOfSquaredDiffFromMeanAmounts += math.Pow(float64(productEntry.Amount)-predictor.expectedValuesOfAmounts[productName], 2.0)
// 		}

// 		predictor.varianceOfPrices[productName] = sumOfSquaredDiffFromMeanPrices / numberOfProductsF
// 		predictor.varianceOfAmounts[productName] = float64(sumOfSquaredDiffFromMeanAmounts) / numberOfProductsF
// 	}

// 	return predictor
// }

// func (p *Predictor) InsertProductToDatabase(pe *database.ProductEntry, prodcutName string) (bool, string) {
// 	if p.chebyshevInequalityCheck(p.expectedValuesOfPrices[prodcutName], p.varianceOfPrices[prodcutName], pe.Price) {
// 		return false, "Value of price is unlikely, are you sure?"
// 	}

// 	if p.chebyshevInequalityCheck(p.expectedValuesOfAmounts[prodcutName], p.varianceOfAmounts[prodcutName], float64(pe.Amount)) {
// 		return false, "Value of amount is unlikely, are you sure?"
// 	}

// 	p.database.Products[prodcutName] = append(p.database.Products[prodcutName], pe)
// 	p.recalculate(pe, prodcutName)

// 	return true, ""
// }

// func (p *Predictor) chebyshevInequalityCheck(exVal, variance float64, element float64) bool {
// 	fmt.Println(exVal, variance, p.sensitivity)
// 	a := math.Sqrt(variance / p.sensitivity)
// 	fmt.Println(a, element)
// 	return math.Abs(element-exVal) > a
// }

// func (p *Predictor) recalculate(pe *database.ProductEntry, prodcutName string) {
// 	currNoOfItems := float64(len(p.database.Products[prodcutName]))
// 	prevNoOfItems := currNoOfItems - 1.0

// 	previousExValOfPrices := p.expectedValuesOfPrices[prodcutName]
// 	previousExValOfAmounts := p.expectedValuesOfAmounts[prodcutName]

// 	p.expectedValuesOfPrices[prodcutName] = (pe.Price / currNoOfItems) + ((prevNoOfItems / currNoOfItems) * previousExValOfPrices)
// 	p.expectedValuesOfAmounts[prodcutName] = (float64(pe.Amount) / currNoOfItems) + ((prevNoOfItems / currNoOfItems) * previousExValOfAmounts)

// 	p.varianceOfPrices[prodcutName] = (1.0 / (currNoOfItems)) * (prevNoOfItems*p.varianceOfPrices[prodcutName] +
// 		prevNoOfItems*(prevNoOfItems-1)*math.Pow(previousExValOfPrices-p.expectedValuesOfPrices[prodcutName], 2.0))

// 	p.varianceOfAmounts[prodcutName] = (1.0 / (currNoOfItems)) * (prevNoOfItems*p.varianceOfAmounts[prodcutName] +
// 		prevNoOfItems*(prevNoOfItems-1)*math.Pow(previousExValOfAmounts-p.expectedValuesOfAmounts[prodcutName], 2.0))
// }
