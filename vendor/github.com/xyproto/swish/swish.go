package swish

import "math"

// Thanks https://codingforspeed.com/using-faster-exponential-approximation/
func exp256(x float64) float64 {
	x = 1.0 + x/256.0
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	return x
}

// Sigmoid is the 1 / (1 + exp(-x)) activation function
func Sigmoid(x float64) float64 {
	// Uses exp256 instead of math.Exp
	return 1.0 / (1.0 + exp256(-x))
}

// Swish is the x / (1 + exp(-x)) activation function, using math.Exp
func Swish(x float64) float64 {
	return x / (1.0 + exp256(-x))
}

// SigmoidPrecise is the 1 / (1 + exp(-x)) activation function, using math.Exp
func SigmoidPrecise(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// SwishPrecise is the x / (1 + exp(-x)) activation function, using math.Exp
func SwishPrecise(x float64) float64 {
	return x / (1.0 + math.Exp(-x))
}
