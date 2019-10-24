package main

import (
	"fmt"
	"math/rand"

	"github.com/xyproto/wann"
)

func main() {

	fmt.Println("### Up detection ###")

	//  o
	// ooo
	up := []float64{0.0, 1.0, 0.0, 1.0, 1.0, 1.0}

	// ooo
	//  o
	down := []float64{1.0, 1.0, 1.0, 0.0, 1.0, 0.0}

	// ooo
	//   o
	left := []float64{1.0, 1.0, 1.0, 0.0, 0.0, 1.0}

	// ooo
	// o
	right := []float64{1.0, 1.0, 1.0, 0.1, 0.0, 0.0}

	_ = up
	_ = down
	_ = left
	_ = right

	// ---

	// Seed based on the current time
	//seed := time.Now().UTC().UnixNano()

	// Seed that makes the program crash
	var seed int64 = 1571917826405889425
	fmt.Printf("seed: %d\n", seed)
	rand.Seed(seed)

	inputNumbers := up

	// 1. Create a neural network that is supposed to be able to detect "up"
	// 2. Use the inputs from up, down, left, right.
	// 3. The goal is that the output neuron should fire "1" for the up patterns and "0" for the rest. (or at least a high/low value)
	// 4. Train, according to the paper.

	config := &wann.Config{
		Inputs:                          len(inputNumbers),
		ConnectionRatio:                 0.5,
		SharedWeight:                    1.0,
		Generations:                     500,
		PopulationSize:                  100,
		MaxIterationsWithoutImprovement: 20,
		MaxModificationIterations:       100,
		Verbose:                         true,
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

	if err := bestNetwork.SaveDiagram("best.svg"); err != nil {
		panic(err)
	}

	// Now test the best network on 4 different inputs

	fmt.Println("Testing the network.")

	upScore := bestNetwork.Evaluate(up)
	downScore := bestNetwork.Evaluate(down)
	leftScore := bestNetwork.Evaluate(left)
	rightScore := bestNetwork.Evaluate(right)

	if upScore > downScore && upScore > leftScore && upScore > rightScore {
		fmt.Println("Network training complete, the results are good.")
	} else {
		fmt.Println("Network training incomplete, the results are not great.")
	}
}
