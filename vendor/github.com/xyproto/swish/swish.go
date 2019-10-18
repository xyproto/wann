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

// Swish is the x / (1 + exp(-x)) activation function, using exp256
func Swish(x float64) float64 {
	return x / (1.0 + exp256(-x))
}

// Sigmoid is the 1 / (1 + exp(-x)) activation function, using exp256
func Sigmoid(x float64) float64 {
	// Uses exp256 instead of math.Exp
	return 1.0 / (1.0 + exp256(-x))
}

// SoftPlus is the log(1 + exp(x)) function, using exp256
func SoftPlus(x float64) float64 {
	return math.Log(1.0 + exp256(x))
}

// Gaussian01 is the Gaussian function with mean 0 and sigma 1, using exp256
func Gaussian01(x float64) float64 {
	return exp256(-(x * x) / 2.0)
}

// SwishPrecise is the x / (1 + exp(-x)) activation function, using math.Exp
func SwishPrecise(x float64) float64 {
	return x / (1.0 + math.Exp(-x))
}

// SigmoidPrecise is the 1 / (1 + exp(-x)) activation function, using math.Exp
func SigmoidPrecise(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// SoftPlusPrecise is the log(1 + exp(x)) function, using math.Exp
func SoftPlusPrecise(x float64) float64 {
	return math.Log(1.0 + math.Exp(x))
}

// Gaussian01 is the Gaussian function with mean 0 and sigma 1, using math.Exp
func Gaussian01Precise(x float64) float64 {
	return math.Exp(-(x * x) / 2.0)
}
