package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xyproto/wann"
)

func init() {
	// Seed based on the current time
	seed := time.Now().UTC().UnixNano()

	// Use a specific seed
	//var seed int64 = 1571917826405889425

	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

func main() {
	// Here are four shapes, up, down, left and right:

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
		0.1, 0.0, 0.0} // o

	// 1. Create a neural network that is supposed to be able to detect "up"
	// 2. Use the inputs from up, down, left, right.
	// 3. The goal is that the output neuron should fire "1" for the up patterns and "0" for the rest. (or at least a high/low value)
	// 4. Train, according to the paper.

	config := &wann.Config{
		Inputs:                                 0,
		ConnectionRatio:                        0.1,
		SharedWeight:                           0.0,
		Generations:                            2000,
		PopulationSize:                         100,
		MaxIterationsWithoutBestImprovement:    1000,
		MaxIterationsWithoutAverageImprovement: 1000,
		Verbose:                                true,
	}

	inputData := make([][]float64, 4)
	inputData[0] = up
	inputData[1] = down
	inputData[2] = left
	inputData[3] = right

	bestNetwork, err := config.Evolve(inputData, []float64{1.0, -1.0, -1.0, -1.0})
	if err != nil {
		panic(err)
	}

	// Now test the best network on 4 different inputs and see if it passes the test

	fmt.Println("Testing the network.")

	upScore := bestNetwork.Evaluate(up)
	downScore := bestNetwork.Evaluate(down)
	leftScore := bestNetwork.Evaluate(left)
	rightScore := bestNetwork.Evaluate(right)

	if upScore > downScore && upScore > leftScore && upScore > rightScore {
		fmt.Println("Network training complete, the results are good.")
	} else {
		fmt.Println("Network training complete, but the results did not pass the test.")
	}

	// Save the image as an SVG image
	if err := bestNetwork.WriteSVG("best.svg"); err != nil {
		panic(err)
	}
}
