package wann

import (
	"fmt"
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

	// Output the diagram
	fmt.Println("--- <DIAGRAM> ---")
	fmt.Println(net)

	err := net.SaveDiagram("/tmp/output.svg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("xdg-open /tmp/output.svg")
	fmt.Println("--- </DIAGRAM> ---")
}
