package wann

import (
	"fmt"
	"github.com/xyproto/onthefly"
)

func (net *Network) OutputDiagram(filename string) {
	l := len(net.InputNodes)
	fmt.Printf("%d input nodes\n", l)

	// Start out by creating a document, 4px padding

	width := 4 + 24 + 64 + 24 + 4
	height := 4 + l*24 + 4

	page, svg := onthefly.NewTinySVG(0, 0, width, height)
	desc := svg.AddNewTag("desc")
	desc.AddContent("generated with github.com/xyproto/wann")

	bg := svg.AddRoundedRect(0, 0, 10, 10, width, height)
	bg.Fill2(onthefly.ColorByName("white"))

	for i, n := range net.InputNodes {
		x := 6
		// 24 pixels per node, including padding (4 pixels above, 4 pixels below)
		y := (i * (20 + 4)) + 4

		rr := svg.AddRoundedRect(x, y, 5, 5, 20, 20)
		rr.Fill2(onthefly.ColorByName("yellow"))

		if net.OutputNode.HasInput(n) {
			svg.Line(x, y, 100, 4, 8, "red")
		}

	}

	page.SaveSVG(filename)
}
