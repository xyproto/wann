package swish

import (
	"math"
)

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// The SWISH activation function
func F(x float64) float64 {
	return x / (1.0 + math.Exp(-x))
}
