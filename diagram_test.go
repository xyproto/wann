package wann

import (
	"fmt"
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
	err := net.SaveDiagram("/tmp/output.svg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("xdg-open /tmp/output.svg")
}
