package main

import (
	"fmt"
	"os"

	"github.com/xyproto/wann"
)

func main() {
	// Here are four shapes, representing: up, down, left and right:

	up := []float64{
		0.0, 1.0, 0.0, //  o
		1.0, 1.0, 1.0} // ooo

	down := []float64{
		1.0, 1.0, 1.0, // ooo
		0.0, 1.0, 0.0} //  o

	left := []float64{
		1.0, 1.0, 1.0, // ooo
		0.0, 0.0, 1.0} //   o

	right := []float64{
		1.0, 1.0, 1.0, // ooo
		1.0, 0.0, 0.0} // o

	// Prepare the input data as a 2D slice
	inputData := [][]float64{
		up,
		down,
		left,
		right,
	}

	// Which of the elements in the input data are we trying to identify?
	correctResultsForUp := []float64{0.0, -1.0, -1.0, -1.0}

	// Prepare a neural network configuration struct
	config := &wann.Config{
		InitialConnectionRatio: 0.2,
		Generations:            2000,
		PopulationSize:         500,
		Verbose:                true,
	}

	// Evolve a network, using the input data and the sought after results
	trainedNetwork, err := config.Evolve(inputData, correctResultsForUp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	// Now to test the trained network on 4 different inputs and see if it passes the test
	upScore := trainedNetwork.Evaluate(up)
	downScore := trainedNetwork.Evaluate(down)
	leftScore := trainedNetwork.Evaluate(left)
	rightScore := trainedNetwork.Evaluate(right)

	if config.Verbose {
		if upScore > downScore && upScore > leftScore && upScore > rightScore {
			fmt.Println("Network training complete, the results are good.")
		} else {
			fmt.Println("Network training complete, but the results did not pass the test.")
		}
	}

	// Save the trained network as an SVG image
	if config.Verbose {
		fmt.Print("Writing network.svg...")
	}
	if err := trainedNetwork.WriteSVG("network.svg"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	if config.Verbose {
		fmt.Println("ok")
	}
}
