package wann

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNetwork_GoExpression_first(t *testing.T) {
	// First create a network with only one output node, that has a step function
	net := NewNetwork()
	net.AllNodes[net.OutputNode].ActivationFunction = Step

	// Output a Go expression for this network, using the given input variable names
	expression := net.GoExpression("x")
	if expression != "func(s float64) float64 { if s >= 0 { return 1 } else { return 0 } }(x)" {
		t.Fail()
	}
	//fmt.Println(expression)
}

func TestNetwork_GoExpression_second(t *testing.T) {
	// Then create a network with an input node that has a sigmoid function and an output node that has an invert function
	net := NewNetwork()
	net.NewInputNode(Sigmoid, true)
	net.AllNodes[net.OutputNode].ActivationFunction = Inv

	// Output a Go expression for this network, using the given input variable names
	expression := net.GoExpression("x")
	if expression != "-(1.0 / (1.0 + math.Exp(-x)))" {
		t.Fail()
	}
	//fmt.Println(expression)
}

func ExampleNetwork_GoFunction_first() {
	rand.Seed(999)
	net := NewNetwork(&Config{
		inputs:                 1,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.GoFunction())

	// Output:
	// func f(x float64) float64 { return math.Exp(-(x * x) / 2.0) }
}
func ExampleNetwork_GoFunction_second() {
	rand.Seed(1111113)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.GoFunction("x1", "x2", "x3", "x4", "x5"))

	// Output:
	// func f(x1, x2, x3, x4, x5 float64) float64 { return func(r float64) float64 { if r >= 0 { return r } else { return 0 } }((((x1/ (1.0 + math.Exp(-x1))) * 0.5 + (x2/ (1.0 + math.Exp(-x2))) * 0.5 + math.Log(1.0 + math.Exp(x3)) * 0.5 + (1.0 / (1.0 + math.Exp(-x4))) * 0.5 + math.Sin(math.Pi * x5) * 0.5) / 5.0)) }
}
