// Package af provides several activation functions that can be used in neural networks
package af

import (
	"github.com/xyproto/swish"
	"math"
)

// The swish package offers optimized Swish, Sigmoid
// SoftPlus and Gaussian01 activation functions
var (
	Sigmoid    = swish.Sigmoid
	Swish      = swish.Swish
	SoftPlus   = swish.SoftPlus
	Gaussian01 = swish.Gaussian01
	Linear     = func(x float64) float64 { return x }
	Inv        = func(x float64) float64 { return -x }
	Sin        = func(x float64) float64 { return math.Sin(math.Pi * x) }
	Cos        = func(x float64) float64 { return math.Cos(math.Pi * x) }
	Squared    = func(x float64) float64 { return x * x }
	Tanh       = math.Tanh
	Abs        = math.Abs
)

// Step function
func Step(x float64) float64 {
	if x >= 0 {
		return 1
	}
	return 0
}

// ReLU is the "rectified linear unit"
// `x >= 0 ? x : 0`
func ReLU(x float64) float64 {
	if x >= 0 {
		return x
	}
	return 0
}

// PReLU is the parametric rectified linear unit.
// `x >= 0 ? x : a * x`
func PReLU(x, a float64) float64 {
	if x >= 0 {
		return x
	}
	return a * x
}
