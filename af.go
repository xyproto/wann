package wann

import (
	"github.com/xyproto/af"
)

const (
	Linear = iota + 1
	Step
	Sin
	Gauss
	Tanh
	Sigmoid
	Inv
	Abs
	Relu
	Cos
	Squared
)

// ActivationFunctions is a collection of activation functions, where the keys are constants that are defined above
// https://github.com/google/brain-tokyo-workshop/blob/master/WANNRelease/WANN/wann_src/ind.py
var ActivationFunctions = map[int](func(float64) float64){
	Linear:  af.Linear,     // Linear
	Step:    af.Step,       // Unsigned Step Function
	Sin:     af.Sin,        // Sin
	Gauss:   af.Gaussian01, // Gaussian with mean 0 and sigma 1
	Tanh:    af.Tanh,       // Hyperbolic Tangent (signed?)
	Sigmoid: af.Sigmoid,    // Sigmoid (unsigned?)
	Inv:     af.Inv,        // Inverse
	Abs:     af.Abs,        // Absolute value
	Relu:    af.ReLU,       // Rectified linear unit
	Cos:     af.Cos,        // Cosine
	Squared: af.Squared,    // Squared
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
