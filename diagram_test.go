package wann

import (
	"math/rand"
	"testing"

	"github.com/xyproto/af"
)

func TestDiagram(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})

	// Set a few activation functions
	net.AllNodes[net.InputNodes[0]].ActivationFunction = af.Linear
	net.AllNodes[net.InputNodes[1]].ActivationFunction = af.Swish
	net.AllNodes[net.InputNodes[2]].ActivationFunction = af.Gaussian01
	net.AllNodes[net.InputNodes[3]].ActivationFunction = af.Sigmoid
	net.AllNodes[net.InputNodes[4]].ActivationFunction = af.ReLU

	// Save the diagram as an image
	err := net.SaveDiagram("test.svg")
	if err != nil {
		t.Error(err)
	}
	//os.Remove("test.svg")
}
