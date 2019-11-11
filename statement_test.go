package wann

import (
	"fmt"
	"math/rand"
	"testing"
)

// ExampleNetwork_Trace
func TestNetwork_Trace(t *testing.T) {
	rand.Seed(1)
	net := NewNetwork(&Config{
		inputs:                 6,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})

	net.SetInputValues([]float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5})
	statement, err := net.StatementWithInputValues()
	if err != nil {
		panic(err)
	}
	fmt.Println(Render(statement))
}

// ExampleNetwork_Trace2
func TestNetwork_Trace2(t *testing.T) {
	rand.Seed(1)
	net := NewNetwork(&Config{
		inputs:                 6,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})

	// 1.234 should not appear in the output statement
	net.SetInputValues([]float64{1.234, 1.234, 1.234, 1.234, 1.234, 1.234})

	statement, err := net.StatementWithInputDataVariables()
	if err != nil {
		panic(err)
	}
	fmt.Println(Render(statement))
}

// ExampleNeuron_InputStatement
func ExampleNeuron_InputStatement() {
	rand.Seed(1)
	net := NewNetwork(&Config{
		inputs:                 6,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})

	// 1.234 should not appear in the output statement
	net.SetInputValues([]float64{1.234, 1.234, 1.234, 1.234, 1.234, 1.234})

	inputStatement2, err := net.AllNodes[net.InputNodes[2]].InputStatement()
	if err != nil {
		panic(err)
	}
	fmt.Println(Render(inputStatement2))
	// Output:
	// inputData[2]
}

func ExampleNetwork_OutputNodeStatementX_first() {
	// First create a network with only one output node, that has a step function
	net := NewNetwork()
	net.AllNodes[net.OutputNode].ActivationFunction = Step

	fmt.Println(net.OutputNodeStatementX("score"))

	// Output:
	// score := func(s float64) float64 {
	// 	if s >= 0 {
	// 		return 1
	// 	} else {
	// 		return 0
	// 	}
	// }(x)
}

func ExampleNetwork_OutputNodeStatementX_second() {
	// Then create a network with an input node that has a sigmoid function and an output node that has an invert function
	net := NewNetwork()
	net.NewInputNode(Sigmoid, true)
	net.AllNodes[net.OutputNode].ActivationFunction = Inv

	// Output a Go expression for this network, using the given input variable names
	fmt.Println(net.OutputNodeStatementX("score"))

	// Output:
	// score := -(x)

}

func ExampleNetwork_OutputNodeStatementX_third() {
	rand.Seed(999)
	net := NewNetwork(&Config{
		inputs:                 1,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.OutputNodeStatementX("score"))

	// Output:
	// score := math.Exp(-(math.Pow(x, 2.0)) / 2.0)
}
func ExampleNetwork_OutputNodeStatementX_fourth() {

	rand.Seed(1111113)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.OutputNodeStatementX("score"))

	// Output:
	// score := func(r float64) float64 {
	//	if r >= 0 {
	//		return r
	//	} else {
	//		return 0
	//	}
	// }(x)
}

func ExampleNetwork_OutputNodeStatementX() {
	rand.Seed(1)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.OutputNodeStatementX("f"))
	//fmt.Println(net.Score())

	// Output:
	// f := math.Pow(x, 2.0)
}
