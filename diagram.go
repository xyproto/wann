package wann

import (
	"bytes"
	"github.com/xyproto/onthefly"
	"io"
	"io/ioutil"
)

// WriteSVG will output the current network as an SVG image to the given io.Writer
func (net *Network) WriteSVG(w io.Writer) (int, error) {
	// Set up margins and the canvas size
	var (
		marginLeft     = 10
		marginTop      = 10
		marginBottom   = 10
		marginRight    = 10
		nodeRadius     = 10
		betweenPadding = 4
		l              = len(net.Nodes)
		width          = marginLeft + nodeRadius + 100 + nodeRadius + marginRight
		height         = marginTop + nodeRadius*2*l + betweenPadding*(l-1) + marginBottom
	)

	// Prepare colors that will be used more than once
	lightYellow := onthefly.ColorByName("#ffffcc")

	// Start a new SVG image
	page, svg := onthefly.NewTinySVG(0, 0, width, height)
	desc := svg.AddNewTag("desc")
	desc.AddContent("generated with github.com/xyproto/wann")

	// White background rounded rectangle
	bg := svg.AddRoundedRect(0, 0, 30, 30, width, height)
	bg.Fill2(onthefly.ColorByName("white"))
	bg.Stroke2(onthefly.ColorByName("black"))

	// Position of output node
	outputx := width - (marginRight + nodeRadius*2)
	outputy := (height - (nodeRadius * 2)) / 2

	// Draw the output node
	output := svg.AddCircle(outputx+nodeRadius, outputy+nodeRadius, nodeRadius)
	output.Fill2(lightYellow)
	output.Stroke2(onthefly.ColorByName("black"))

	// Draw the input nodes
	for i, n := range net.Nodes {

		// Find the position of this node circle
		x := marginLeft
		y := (i * (nodeRadius*2 + betweenPadding)) + marginTop

		// Draw this input node
		input := svg.AddCircle(x+nodeRadius, y+nodeRadius, nodeRadius)
		input.Fill2(lightYellow)
		input.Stroke2(onthefly.ColorByName("black"))

		// Draw the connection from this node to the output node, if applicable
		if net.OutputNode.HasInput(n) {
			svg.Line(x+nodeRadius*2, y+nodeRadius, outputx, outputy+nodeRadius, 2, "#0099ff")
		}
	}

	// Write the data to the given io.Writer
	return w.Write([]byte(page.GetXML(false)))
}

// SaveDiagram saves a drawing of the current network as an SVG file
func (net *Network) SaveDiagram(filename string) error {
	var buf bytes.Buffer
	if _, err := net.WriteSVG(&buf); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
