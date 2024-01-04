package pkg

import (
	"math"
)


func findAverage(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	average := sum / float64(len(numbers))
	return average
}


func filterInfinity(numbers []float64) []float64 {
	var result []float64

	for _, num := range numbers {
		// Check if the number is not positive infinity
		if !math.IsInf(num, 1) {
			result = append(result, num)
		}
	}
	return result
}


func highestFloat(f []float64) (index int, hf float64) {

	hf = 0.0
	for i, num := range  f{
		if num > hf {
			index = i
			hf = num			
		}
	}

	return
}