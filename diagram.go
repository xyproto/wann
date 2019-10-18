package wann

import (
	"fmt"
	"github.com/xyproto/onthefly"
)

// OutputDiagram will output a diagram as an SVG image
func (net *Network) OutputDiagram(filename string) error {
	l := len(net.Nodes)
	fmt.Printf("%d input nodes\n", l)

	marginLeft := 10
	marginTop := 10
	marginBottom := 10
	marginRight := 10
	nodeRadius := 10
	betweenPadding := 4

	lightYellow := onthefly.ColorByName("#ffffcc")

	// Start out by creating a document, 4px padding

	width := marginLeft + nodeRadius + 100 + nodeRadius + marginRight
	height := marginTop + nodeRadius*2*l + betweenPadding*(l-1) + marginBottom

	page, svg := onthefly.NewTinySVG(0, 0, width, height)
	desc := svg.AddNewTag("desc")
	desc.AddContent("generated with github.com/xyproto/wann")

	// White background rounded rectangle
	bg := svg.AddRoundedRect(0, 0, 30, 30, width, height)
	bg.Fill2(onthefly.ColorByName("white"))
	bg.Stroke2(onthefly.ColorByName("black"))

	outputx := width - (marginRight + nodeRadius*2)
	outputy := (height - (nodeRadius * 2)) / 2

	// The output node
	//output := svg.AddRoundedRect(outputx, outputy, 5, 5, 20, 20)
	output := svg.AddCircle(outputx+nodeRadius, outputy+nodeRadius, nodeRadius)
	output.Fill2(lightYellow)
	output.Stroke2(onthefly.ColorByName("black"))

	for i, n := range net.Nodes {
		x := 10
		// 24 pixels per node, including padding (4 pixels above, 4 pixels below)
		y := (i * (20 + 5)) + 5

		//rr := svg.AddRoundedRect(x, y, 5, 5, 20, 20)
		rr := svg.AddCircle(x+10, y+10, 10)
		rr.Fill2(lightYellow)
		rr.Stroke2(onthefly.ColorByName("black"))

		if net.OutputNode.HasInput(n) {
			svg.Line(x+20, y+10, outputx, outputy+10, 1, "#0099ff")
		}

	}

	return page.SaveSVG(filename)
}
