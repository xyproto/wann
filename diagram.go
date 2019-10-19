package wann

import (
	"bytes"
	"github.com/xyproto/onthefly"
	"io"
	"io/ioutil"
)

type Pos struct {
	x float64
	y float64
}

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
		//p := onthefly.NewPosf(0.0, 0.0)
		//prevp := onthefly.NewPosf(0.0, 0.0)
		//first := true
		//fmt.Println("---")
		points := make([]*onthefly.Pos, 0)
		for xpos := startx; xpos < stopx; xpos += 0.2 {
			xr := float64(xpos-startx) / float64(stopx-startx)
			xv := (xr * 8.0) - 5.0
			yv := n.ActivationFunction(xv)
			// plot
			yp := float64(ypos) + float64(nodeRadius) - (yv * 2.5) + float64(nodeRadius)*0.1
			//yz := float64(ypos)+float64(nodeRadius)-(0.0*2.5)+float64(nodeRadius)*0.1
			xp := xpos
			//dot := svg.AddCirclef(xp, yp, 0.5)
			//prevp = p
			p := onthefly.NewPosf(xp, yp)
			//if first {
			//	prevp = p
			//	first = false
			//}
			//svg.Line2(prevp, p, 1, onthefly.ColorByName("red"))
			points = append(points, p)
			//svg.Line2(onthefly.NewPosf(xp, yp), onthefly.NewPosf(xp+0.2, yp-0.2), 2, onthefly.ColorByName("black"))
			//svg.Circle2(onthefly.NewPosf(xp, yp), 1, onthefly.ColorByName("black"))
			//dot.Fill2(onthefly.ColorByName("black"))
			//dot := svg.AddRectf(float64(xpos), float64(ypos)+float64(nodeRadius)-(yv*2.5)+float64(nodeRadius)*0.1, 0.5)
			//dot.Fill2(onthefly.ColorByName("black"))
			//dot.Stroke2(onthefly.ColorByName("black"))
		}

		// Now append all the points in reverse, to make it a line
		//pc2 := points[:]
		//for i := len(pc2)-1; i > 0; i-- {
		//	points = append(points, pc2[i])
		//}

		//points = append(points, onthefly.NewPosf(100, 100))

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
