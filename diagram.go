package wann

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/xyproto/tinysvg"
)

// OutputSVG will output the current network as an SVG image to the given io.Writer
// TODO: Clean up and refactor
func (net *Network) OutputSVG(w io.Writer) (int, error) {
	// Set up margins and the canvas size
	var (
		marginLeft     = 10
		marginTop      = 10
		marginBottom   = 10
		marginRight    = 10
		nodeRadius     = 10
		betweenPadding = 4
		d              = float64(net.Depth()) * 2.5
		width          = marginLeft + int(float64(nodeRadius)*2.0*d) + betweenPadding*(int(d)-1) + nodeRadius + marginRight
		l              = float64(len(net.InputNodes))
		height         = marginTop + int(float64(nodeRadius)*1.5*l) + betweenPadding*(int(l)-1) + marginBottom
		imgPadding     = 5
		lineWidth      = 2
	)

	if width < 128 {
		width = 128
	}
	if height < 128 {
		height = 128
	}

	// Start a new SVG image
	document, svg := tinysvg.NewTinySVG(width+imgPadding*2, height+imgPadding*2)
	svg.Describe("generated with github.com/xyproto/wann")

	// White background rounded rectangle
	bg := svg.AddRoundedRect(imgPadding, imgPadding, 30, 30, width, height)
	bg.Fill2(tinysvg.ColorByName("white"))
	bg.Stroke2(tinysvg.ColorByName("black"))

	// Position of output node
	outputx := width - (marginRight + nodeRadius*2) + imgPadding
	outputy := (height-(nodeRadius*2))/2 + imgPadding

	// For each connected neuron, store it with the distance from the output neuron as the key in a map
	layerNeurons := make(map[int][]NeuronIndex)
	maxDistance := 0
	net.ForEachConnectedNodeIndex(func(ni NeuronIndex) {
		distanceFromOutput := net.AllNodes[ni].distanceFromOutputNode
		layerNeurons[distanceFromOutput] = append(layerNeurons[distanceFromOutput], ni)
		if distanceFromOutput > maxDistance {
			maxDistance = distanceFromOutput
		}
	})

	// Draw the input nodes as circles, and connections to the output node as lines
	//for i, n := range net.InputNodes {
	columnOffset := 50

	getPosition := func(givenNeuron NeuronIndex) (int, int) {
		for outputDistance, neurons := range layerNeurons {
			for neuronLayerIndex, otherNeuron := range neurons {
				if otherNeuron == givenNeuron {
					x := marginLeft + imgPadding + columnOffset*(maxDistance-outputDistance)
					y := (neuronLayerIndex * (nodeRadius*2 + betweenPadding)) + marginTop + imgPadding
					return x, y
				}
			}
		}
		panic("implementation error: neuron index not found")
	}

	// Draw node lines first
	for _, neurons := range layerNeurons {
		for _, neuronIndex := range neurons {
			if neuronIndex == net.OutputNode {
				continue
			}
			// Find the position of this node circle
			x, y := getPosition(neuronIndex)
			// Draw the connection from the center of this node to the center of all input nodes, if applicable
			for _, inputNeuron := range (net.AllNodes[neuronIndex]).InputNodes {
				ix, iy := getPosition(inputNeuron)
				svg.Line(ix+nodeRadius, iy+nodeRadius, x+nodeRadius, y+nodeRadius, lineWidth, "orange")
			}
			// Draw the connection to the output node, if it has this node as input
			if net.AllNodes[net.OutputNode].HasInput(neuronIndex) {
				svg.Line(x+nodeRadius, y+nodeRadius, outputx+nodeRadius, outputy+nodeRadius, lineWidth, "#0099ff")
			}
		}
	}

	// Then draw the nodes on top, including graph plots
	for _, neurons := range layerNeurons {
		for _, neuronIndex := range neurons {
			if neuronIndex == net.OutputNode {
				continue
			}

			// Find the position of this node circle
			x, y := getPosition(neuronIndex)

			// Draw this node
			input := svg.AddCircle(x+nodeRadius, y+nodeRadius, nodeRadius)
			switch net.AllNodes[neuronIndex].distanceFromOutputNode {
			case 1, 6:
				input.Fill("lightblue")
			case 2, 7:
				input.Fill("lightgreen")
			case 3, 8:
				input.Fill("lightyellow")
			case 4, 9:
				input.Fill("orange")
			case 5, 10:
				input.Fill("red")
			default:
				input.Fill("gray")
			}
			input.Stroke2(tinysvg.ColorByName("black"))

			// Plot the activation function inside this node
			var points []*tinysvg.Pos
			startx := float64(x) + float64(nodeRadius)*0.5
			stopx := float64(x+nodeRadius*2) - float64(nodeRadius)*0.5
			ypos := float64(y)
			for xpos := startx; xpos < stopx; xpos += 0.2 {
				// xr is from 0 to 1
				xr := float64(xpos-startx) / float64(stopx-startx)
				// xv is from -5 to 5
				xv := (xr - 0.5) * float64(nodeRadius)
				node := net.AllNodes[neuronIndex]
				f := ActivationFunctions[node.ActivationFunction]
				yv := f(xv)
				// plot, 3.0 is the amplitude along y
				yp := float64(ypos) + float64(nodeRadius)*1.35 - (yv * 0.6 * float64(nodeRadius))

				if yp < (ypos + float64(nodeRadius)*0.1) {
					continue
				} else if yp > (ypos + float64(nodeRadius)*1.9) {
					continue
				}

				// Label
				name := node.ActivationFunction.Name()
				if net.IsInput(neuronIndex) {
					// Add a the input number to the name
					for i, ni := range net.InputNodes {
						if neuronIndex == ni {
							name += " [" + strconv.Itoa(i) + "]"
						}
					}
				} else if neuronIndex == net.OutputNode {
					name += " !"
				}
				box := svg.AddRect(int(startx-float64(nodeRadius)*0.4), int(ypos+float64(nodeRadius)*2.5)-5, len(name)*5, 6)
				box.Fill("black")
				svg.Text(int(startx-float64(nodeRadius)*0.4), int(ypos+float64(nodeRadius)*2.5), 8, "Courier", name, "white")

				p := tinysvg.NewPosf(xpos, yp)
				points = append(points, p)
			}
			// Draw the polyline (graph)
			pl := svg.Polyline(points, tinysvg.ColorByName("black"))
			pl.Stroke2(tinysvg.ColorByName("black"))
			pl.Fill2(tinysvg.ColorByName("none"))

		}
	}

	// Draw the output node
	output := svg.AddCircle(outputx+nodeRadius+1, outputy+nodeRadius+1, nodeRadius)
	output.Fill("magenta")
	output.Stroke2(tinysvg.ColorByName("black"))

	// Label
	name := net.AllNodes[net.OutputNode].ActivationFunction.Name() + " [o]"
	box := svg.AddRect(outputx-nodeRadius/2, (nodeRadius*2)+outputy+1, len(name)*5, 6)
	box.Fill("black")
	svg.Text(outputx-nodeRadius/2, (nodeRadius*2)+outputy+6, 8, "Courier", name, "white")

	// Write the data to the given io.Writer
	return w.Write(document.Bytes())
}

// WriteSVG saves a drawing of the current network as an SVG file
func (net *Network) WriteSVG(filename string) error {
	var buf bytes.Buffer
	if _, err := net.OutputSVG(&buf); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
