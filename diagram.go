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
		imgPadding     = 5
		lineWidth      = 2
	)

	// Prepare colors that will be used more than once
	lightYellow := onthefly.ColorByName("#ffffcc")

	// Start a new SVG image
	page, svg := onthefly.NewTinySVG(0, 0, width+imgPadding*2, height+imgPadding*2)

	desc := svg.AddNewTag("desc")
	desc.AddContent("generated with github.com/xyproto/wann")

	// White background rounded rectangle
	bg := svg.AddRoundedRect(imgPadding, imgPadding, 30, 30, width, height)
	bg.Fill2(onthefly.ColorByName("white"))
	bg.Stroke2(onthefly.ColorByName("black"))

	// Position of output node
	outputx := width - (marginRight + nodeRadius*2) + imgPadding
	outputy := (height-(nodeRadius*2))/2 + imgPadding

	// Draw the input nodes as circles, and connections to the output node as lines
	for i, n := range net.Nodes {

		// Find the position of this node circle
		x := marginLeft + imgPadding
		y := (i * (nodeRadius*2 + betweenPadding)) + marginTop + imgPadding

		// Draw the connection from the center of this node to the center of the output node, if applicable
		if net.OutputNode.HasInput(n) {
			svg.Line(x+nodeRadius, y+nodeRadius, outputx+nodeRadius, outputy+nodeRadius, lineWidth, "#0099ff")
		}

		// Draw this input node
		input := svg.AddCircle(x+nodeRadius, y+nodeRadius, nodeRadius)
		input.Fill2(lightYellow)
		input.Stroke2(onthefly.ColorByName("black"))

		// Plot the activation function inside this node
		startx := float64(x) + float64(nodeRadius)*0.5
		stopx := float64(x+nodeRadius*2) - float64(nodeRadius)*0.5
		ypos := float64(y)
		var points []*onthefly.Pos
		for xpos := startx; xpos < stopx; xpos += 0.2 {
			// xr is from 0 to 1
			xr := float64(xpos-startx) / float64(stopx-startx)
			// xv is from -5 to 3
			//xv := (xr * 8.0) - 5.0
			// xv is from -2 to 2
			//xv := (xr * 4.0) - 2.0
			// xv is from -5 to 5
			xv := (xr - 0.5) * float64(nodeRadius)
			yv := n.ActivationFunction(xv)
			// plot, 3.0 is the amplitude along y
			yp := float64(ypos) + float64(nodeRadius)*1.35 - (yv * 0.6 * float64(nodeRadius))

			if yp < (ypos + float64(nodeRadius)*0.1) {
				continue
			} else if yp > (ypos + float64(nodeRadius)*1.9) {
				continue
			}
			p := onthefly.NewPosf(xpos, yp)
			points = append(points, p)
		}
		// Draw the polyline (graph)
		pl := svg.Polyline(points, onthefly.ColorByName("black"))
		pl.Stroke2(onthefly.ColorByName("black"))
		pl.Fill2(onthefly.ColorByName("none"))
	}

	// Draw the output node
	output := svg.AddCircle(outputx+nodeRadius, outputy+nodeRadius, nodeRadius)
	output.Fill2(lightYellow)
	output.Stroke2(onthefly.ColorByName("black"))

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
