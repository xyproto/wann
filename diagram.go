package wann

import (
	"bytes"
	"fmt"
	"github.com/xyproto/onthefly"
	"io"
	"io/ioutil"
)

// WriteSVG will output the current network as an SVG image to the given io.Writer
func (net *Network) WriteSVG(w io.Writer) (int, error) {
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
		x := marginLeft
		// 24 pixels per node, including padding (4 pixels above, 4 pixels below)
		y := (i * (nodeRadius*2 + betweenPadding)) + marginTop

		//rr := svg.AddRoundedRect(x, y, 5, 5, 20, 20)
		rr := svg.AddCircle(x+nodeRadius, y+nodeRadius, nodeRadius)
		rr.Fill2(lightYellow)
		rr.Stroke2(onthefly.ColorByName("black"))

		if net.OutputNode.HasInput(n) {
			svg.Line(x+nodeRadius*2, y+nodeRadius, outputx, outputy+nodeRadius, 2, "#0099ff")
		}

	}
	return w.Write([]byte(page.GetXML(false)))
}

// SaveDiagram saves a drawing of the current network as an SVG file
func (net *Network) SaveDiagram(filename string) error {
	var buf bytes.Buffer
	_, err := net.WriteSVG(&buf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
