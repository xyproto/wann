package wann

import (
	"math/rand"
	"os"
	"testing"
)

func TestDiagram(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.5,
		sharedWeight:           0.5,
	})

	// Set a few activation functions
	net.AllNodes[net.InputNodes[0]].ActivationFunctionIndex = Linear
	net.AllNodes[net.InputNodes[1]].ActivationFunctionIndex = Swish
	net.AllNodes[net.InputNodes[2]].ActivationFunctionIndex = Gauss
	net.AllNodes[net.InputNodes[3]].ActivationFunctionIndex = Sigmoid
	net.AllNodes[net.InputNodes[4]].ActivationFunctionIndex = ReLU

	// Save the diagram as an image
	err := net.WriteSVG("test.svg")
	if err != nil {
		t.Error(err)
	}
	os.Remove("test.svg")
}
