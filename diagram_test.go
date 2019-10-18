package wann

import (
	"fmt"
	"testing"
)

func TestDiagram(t *testing.T) {
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	net.OutputDiagram("/tmp/output.svg")
	fmt.Println("xdg-open /tmp/output.svg")
}
