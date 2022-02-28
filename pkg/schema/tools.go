package schema

import (
	"math/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

func getRandomLength() int {
	max := 10.0
	min := 0.0
	return int(distuv.Uniform{Min: min, Max: max}.Rand())
}

func getRandomProbability() float64 {
	return rand.Float64()
}
