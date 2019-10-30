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
	net.AllNodes[net.InputNodes[0]].ActivationFunction = Linear
	net.AllNodes[net.InputNodes[1]].ActivationFunction = Swish
	net.AllNodes[net.InputNodes[2]].ActivationFunction = Gauss
	net.AllNodes[net.InputNodes[3]].ActivationFunction = Sigmoid
	net.AllNodes[net.InputNodes[4]].ActivationFunction = ReLU

	// Save the diagram as an image
	err := net.WriteSVG("test.svg")
	if err != nil {
		t.Error(err)
	}
	os.Remove("test.svg")
}
