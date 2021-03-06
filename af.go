package wann

import (
	"fmt"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/xyproto/af"
)

// ActivationFunctionIndex is a number that represents a specific activation function
type ActivationFunctionIndex int

const (
	// Step is a step. First 0 and then abrubtly up to 1.
	Step ActivationFunctionIndex = iota
	// Linear is the linear activation function. Gradually from 0 to 1.
	Linear
	// Sin is the sinoid activation function
	Sin
	// Gauss is the Gaussian function, with a mean of 0 and a sigma of 1
	Gauss
	// Tanh is math.Tanh
	Tanh
	// Sigmoid is the optimized sigmoid function from github.com/xyproto/swish
	Sigmoid
	// Inv is the inverse linear function
	Inv
	// Abs is math.Abs
	Abs
	// ReLU or ReLU is the rectified linear unit, first 0 and then the linear function
	ReLU
	// Cos is the cosoid (?) activation function
	Cos
	// Squared increases rapidly
	Squared
	// Swish is a later invention than ReLU, _|
	Swish
	// SoftPlus is log(1 + exp(x))
	SoftPlus
)

// ActivationFunctions is a collection of activation functions, where the keys are constants that are defined above
// https://github.com/google/brain-tokyo-workshop/blob/master/WANNRelease/WANN/wann_src/ind.py
var ActivationFunctions = map[ActivationFunctionIndex](func(float64) float64){
	Step:     af.Step,       // Unsigned Step Function
	Linear:   af.Linear,     // Linear
	Sin:      af.Sin,        // Sin
	Gauss:    af.Gaussian01, // Gaussian with mean 0 and sigma 1
	Tanh:     af.Tanh,       // Hyperbolic Tangent (signed?)
	Sigmoid:  af.Sigmoid,    // Sigmoid (unsigned?)
	Inv:      af.Inv,        // Inverse
	Abs:      af.Abs,        // Absolute value
	ReLU:     af.ReLU,       // Rectified linear unit
	Cos:      af.Cos,        // Cosine
	Squared:  af.Squared,    // Squared
	Swish:    af.Swish,      // Swish
	SoftPlus: af.SoftPlus,   // SoftPlus
}

// ComplexityEstimate is a map for having an estimate of how complex each function is,
// based on a quick benchmark of each function.
// The complexity estimates will vary, depending on the performance.
var ComplexityEstimate = make(map[ActivationFunctionIndex]float64)

func (config *Config) estimateComplexity() {
	if config.Verbose {
		fmt.Print("Estimating activation function complexity...")
	}
	startEstimate := time.Now()
	resolution := 0.0001
	durationMap := make(map[ActivationFunctionIndex]time.Duration)
	var maxDuration time.Duration
	for i, f := range ActivationFunctions {
		start := time.Now()
		for x := 0.0; x <= 1.0; x += resolution {
			_ = f(x)
		}
		duration := time.Since(start)
		durationMap[ActivationFunctionIndex(i)] = duration
		if duration > maxDuration {
			maxDuration = duration
		}
	}
	for i := range ActivationFunctions {
		// 1.0 means the function took maxDuration
		ComplexityEstimate[ActivationFunctionIndex(i)] = float64(durationMap[ActivationFunctionIndex(i)]) / float64(maxDuration)
	}
	estimateDuration := time.Since(startEstimate)
	if config.Verbose {
		fmt.Printf(" done. (In %v)\n", estimateDuration)
	}
}

// Call runs an activation function with the given float64 value.
// The activation function is chosen by one of the constants above.
func (afi ActivationFunctionIndex) Call(x float64) float64 {
	if f, ok := ActivationFunctions[afi]; ok {
		return f(x)
	}
	// Use the linear function by default
	return af.Linear(x)
}

// Name returns a name for each activation function
func (afi ActivationFunctionIndex) Name() string {
	switch afi {
	case Step:
		return "Step"
	case Linear:
		return "Linear"
	case Sin:
		return "Sinusoid"
	case Gauss:
		return "Gaussian"
	case Tanh:
		return "Tanh"
	case Sigmoid:
		return "Sigmoid"
	case Inv:
		return "Inverted"
	case Abs:
		return "Absolute"
	case ReLU:
		return "ReLU"
	case Cos:
		return "Cosinusoid"
	case Squared:
		return "Squared"
	case Swish:
		return "Swish"
	case SoftPlus:
		return "SoftPlus"
	default:
		return "Untitled"
	}
}

// goExpression returns the Go expression for this activation function, using the given variable name string as the input variable name
func (afi ActivationFunctionIndex) goExpression(varName string) string {
	switch afi {
	case Step:
		// Using s to not confuse it with the varName
		return "func(s float64) float64 { if s >= 0 { return 1 } else { return 0 } }(" + varName + ")"
	case Linear:
		return varName
	case Sin:
		return "math.Sin(math.Pi * " + varName + ")"
	case Gauss:
		return "math.Exp(-(" + varName + " * " + varName + ") / 2.0)"
	case Tanh:
		return "math.Tanh(" + varName + ")"
	case Sigmoid:
		return "(1.0 / (1.0 + math.Exp(-" + varName + ")))"
	case Inv:
		return "-" + varName
	case Abs:
		return "math.Abs(" + varName + ")"
	case ReLU:
		// Using r to not confuse it with the varName
		return "func(r float64) float64 { if r >= 0 { return r } else { return 0 } }(" + varName + ")"
	case Cos:
		return "math.Cos(math.Pi * " + varName + ")"
	case Squared:
		return "(" + varName + " * " + varName + ")"
	case Swish:
		return "(" + varName + "/ (1.0 + math.Exp(-" + varName + ")))"
	case SoftPlus:
		return "math.Log(1.0 + math.Exp(" + varName + "))"
	default:
		return varName
	}
}

// String returns the Go expression for this activation function, using "x" as the input variable name
func (afi ActivationFunctionIndex) String() string {
	return afi.goExpression("x")
}

// Statement returns the Statement statement for this activation function, using the given inner statement
func (afi ActivationFunctionIndex) Statement(inner *jen.Statement) *jen.Statement {
	switch afi {
	case Step:
		// func(s float64) float64 { if s >= 0 { return 1 } else { return 0 } }(inner)
		// Using s to not confuse it with the varName
		return jen.Func().Params(jen.Id("s").Id("float64")).Id("float64").Block(
			jen.If(jen.Id("s").Op(">=").Id("0")).Block(
				jen.Return(jen.Lit(1)),
			).Else().Block(
				jen.Return(jen.Lit(0)),
			),
		).Call(inner)
	case Cos:
		// math.Cos((inner) * math.Pi)
		return jen.Qual("math", "Cos").Call(jen.Parens(inner).Op("*").Id("math").Dot("Pi"))
	case Sin:
		// math.Sin((inner) * math.Pi)
		return jen.Qual("math", "Sin").Call(jen.Parens(inner).Op("*").Id("math").Dot("Pi"))
	case Gauss:
		// return math.Exp(-(math.Pow(inner, 2.0)) / 2.0)
		return jen.Qual("math", "Exp").Call(jen.Op("-").Parens(
			// Using math.Pow ensures the inner expression is only calculated once, if it's a large expression
			//inner.Op("*").Add(inner),
			jen.Qual("math", "Pow").Params(
				inner,
				jen.Lit(2.0),
			),
		).Op("/").Lit(2.0))
	case Tanh:
		// math.Tanh(inner)
		return jen.Qual("math", "Tanh").Call(inner)
	case Sigmoid:
		// (1.0 / (1.0 + math.Exp(-(inner))))
		return jen.Lit(1.0).Op("/").Parens(jen.Lit(1.0).Op("+").Qual("math", "Exp").Call(jen.Op("-").Parens(inner)))
	case Inv:
		// -(inner)
		return jen.Op("-").Parens(inner)
	case Abs:
		// math.Abs(inner)
		return jen.Qual("math", "Abs").Call(inner)
	case ReLU:
		//return "func(r float64) float64 { if r >= 0 { return r } else { return 0 } }(" + varName + ")"
		// Using r to not confuse it with the varName
		return jen.Func().Params(jen.Id("r").Id("float64")).Id("float64").Block(
			jen.If(jen.Id("r").Op(">=").Id("0")).Block(
				jen.Return(jen.Id("r")),
			).Else().Block(
				jen.Return(jen.Lit(0)),
			),
		).Call(inner)
	case Squared:
		// inner^2
		//return inner.Op("*").Add(inner)
		// Using math.Pow ensures the inner expression is only calculated once, if it's a large expression
		return jen.Qual("math", "Pow").Call(inner, jen.Lit(2.0))
	case Swish:
		// (inner / (1.0 + math.Exp(-inner)))
		return jen.Parens(inner.Op("/").Parens(jen.Lit(1.0).Op("+").Qual("math", "Exp").Call(jen.Op("-").Parens(inner))))
	case SoftPlus:
		// math.Log(1.0 + math.Exp(inner))
		return jen.Qual("math", "Log").Call(jen.Lit(1.0).Op("+").Qual("math", "Exp").Call(inner))
	case Linear:
		// This is also the default case: (inner)
		fallthrough
	default:
		// (inner)
		return jen.Parens(inner)
	}
}

// GoRun will first construct the expression using jennifer and then evaluate the result using "go run" and a source file innn /tmp
func (afi ActivationFunctionIndex) GoRun(x float64) (float64, error) {
	return RunStatementX(afi.Statement(jen.Id("x")), x)
}
