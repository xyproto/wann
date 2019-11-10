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
