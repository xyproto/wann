package wann

import (
	"math/rand"
	"testing"
)

func TestDiagram(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})

	// net.AllNodes[net.InputNodes[0]].ActivationFunction = af.Linear
	// net.AllNodes[net.InputNodes[1]].ActivationFunction = af.Swish
	// net.AllNodes[net.InputNodes[2]].ActivationFunction = af.Gaussian01
	// net.AllNodes[net.InputNodes[3]].ActivationFunction = af.Sigmoid
	// net.AllNodes[net.InputNodes[4]].ActivationFunction = af.ReLU

	// Set a few activation functions
	net.AllNodes[net.InputNodes[0]].ActivationFunctionIndex = Linear
	net.AllNodes[net.InputNodes[1]].ActivationFunctionIndex = Swish
	net.AllNodes[net.InputNodes[2]].ActivationFunctionIndex = Gauss
	net.AllNodes[net.InputNodes[3]].ActivationFunctionIndex = Sigmoid
	net.AllNodes[net.InputNodes[4]].ActivationFunctionIndex = ReLU

	// Save the diagram as an image
	err := net.SaveDiagram("test.svg")
	if err != nil {
		t.Error(err)
	}
	//os.Remove("test.svg")
}
