package wann

import (
	"github.com/xyproto/af"
)

const (
	LINEAR = iota + 1
	STEP
	SIN
	GAUSS
	TANH
	SIGMOID
	INV
	ABS
	RELU
	COS
	SQUARED
)

// https://github.com/google/brain-tokyo-workshop/blob/master/WANNRelease/WANN/wann_src/ind.py
var ActivationFunctions = map[int](func(float64) float64){
	LINEAR:  af.Linear,     // Linear
	STEP:    af.Step,       // Unsigned Step Function
	SIN:     af.Sin,        // Sin
	GAUSS:   af.Gaussian01, // Gaussian with mean 0 and sigma 1
	TANH:    af.Tanh,       // Hyperbolic Tangent (signed?)
	SIGMOID: af.Sigmoid,    // Sigmoid (unsigned?)
	INV:     af.Inv,        // Inverse
	ABS:     af.Abs,        // Absolute value
	RELU:    af.ReLU,       // Rectified linear unit
	COS:     af.Cos,        // Cosine
	SQUARED: af.Squared,    // Squared
}

// Calc runs an activation function with the given float64 value.
// The activation function is chosen by one of the constants above.
func Calc(functionIndex int, x float64) float64 {
	if f, ok := ActivationFunctions[functionIndex]; ok {
		return f(x)
	}
	// Use the linear function by default
	return af.Linear(x)
}
