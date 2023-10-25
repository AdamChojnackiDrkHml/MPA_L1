package dataFunctions

import (
	"math"

	"github.com/sajari/regression"
)

func FindLinearRegressionForDataPoints(dataPoints [][]float64) []float64 {
	r := new(regression.Regression)
	r.SetVar(0, "MeasureIndex")
	dataPointsP := regression.MakeDataPoints(dataPoints, 0)
	r.Train(dataPointsP...)
	r.Run()

	return r.GetCoeffs()
}

func RemoveOutliers(dataPoints [][]float64, coeffs []float64) (clearedDataPoints [][]float64) {
	expectedValue := make([]float64, len(dataPoints))

	for i := range expectedValue {
		expectedValue[i] = dataPoints[i][0]*coeffs[0] + coeffs[1]
	}

	sumOfSquaresError := 0.0

	for i := range expectedValue {
		sumOfSquaresError += math.Pow((dataPoints[i][1])-expectedValue[i], 2.0)
	}

	s := math.Sqrt(sumOfSquaresError / float64(len(dataPoints)))

	for i := range dataPoints {
		if math.Abs(dataPoints[i][1]-expectedValue[i]) < math.Abs(2.0*s) {
			clearedDataPoints = append(clearedDataPoints, dataPoints[i])
		}
	}

	return clearedDataPoints
}

func ClearData(dataPoints [][]float64) [][]float64 {

	coeffPrice := FindLinearRegressionForDataPoints(dataPoints)

	clearedPrices := RemoveOutliers(dataPoints, coeffPrice)

	return clearedPrices
}
