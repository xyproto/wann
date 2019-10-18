package wann

import (
	"math/rand"
	"testing"
)

func TestDiagram(t *testing.T) {
	rand.Seed(currentTime)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	net.OutputDiagram("/tmp/output.svg")
}
