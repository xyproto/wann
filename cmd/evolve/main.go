package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xyproto/wann"
)

var currentTime = time.Now().UTC().UnixNano()

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

	rand.Seed(currentTime)

	inputNumbers := up

	// 1. Create a neural network that is supposed to be able to detect "up"
	// 2. Use the inputs from up, down, left, right.
	// 3. The goal is that the output neuron should fire "1" for the up patterns and "0" for the rest. (or at least a high/low value)
	// 4. Train, according to the paper.

	config := &wann.Config{
		Inputs:          len(inputNumbers),
		ConnectionRatio: 0.5,
		SharedWeight:    1.0,
	}

	// 100 generations
	G := 100

	// population of 100 networks
	N := 100
	population := make([]*wann.Network, N)

	// Initialize the population
	for i := 0; i < N; i++ {
		population[i] = wann.NewNetwork(config)
	}

	var bestNetwork *wann.Network

	// For each generation, evaluate and modify the networks
	for j := 0; j < G; j++ {

		bestWeight := 0.0
		bestScore := 0.0
		bestNetwork = nil

		// For each weight, evaluate all networks
		for w := 0.0; w <= 1.0; w += 0.1 {

			scoreMap := make(map[int]float64)
			scoreSum := 0.0
			for i := 0; i < N; i++ {
				net := population[i]
				net.SetWeight(w)
				result := net.Evaluate(up) - (net.Evaluate(down) + net.Evaluate(right) + net.Evaluate(left))
				score := result / (net.Complexity() * 0.1)
				scoreSum += score
				scoreMap[i] = score
			}

			// The scores for this weight
			scoreList := wann.SortByValue(scoreMap)

			if scoreList[0].Value > bestScore {
				bestScore = scoreList[0].Value
				bestNetwork = population[scoreList[0].Key]
				bestWeight = w
			}

		}

		if bestNetwork == nil {
			panic("no best network?")
		}

		fmt.Println("Best score:", bestScore)
		fmt.Println("Best weight:", bestWeight)

		// Now evaluate the network, but only for the best weight

		w := bestWeight

		scoreMap := make(map[int]float64)
		scoreSum := 0.0
		for i := 0; i < N; i++ {
			net := population[i]
			net.SetWeight(w)
			result := net.Evaluate(up) - (net.Evaluate(down) + net.Evaluate(right) + net.Evaluate(left))
			score := result / (net.Complexity() * 0.1)
			scoreSum += score
			scoreMap[i] = score
		}
		averageScore := scoreSum / float64(N)

		// The scores for this weight
		scoreList := wann.SortByValue(scoreMap)

		//		// For each network, for each weight, evaluate the result
		//		scoreMap := make(map[int]float64)
		//		scoreSum := 0.0
		//		for i := 0; i < N; i++ {
		//			net := population[i]
		//			bestResult := 0.0
		//			complexity := net.Complexity()
		//			for w := 0.0; w <= 1.0; w += 0.1 {
		//				net.SetWeight(w)
		//				result := net.Evaluate(inputNumbers)
		//				if result > bestResult {
		//					bestResult = result
		//					bestWeight = w
		//				}
		//			}
		//			score := bestResult / (complexity * 0.1)
		//			scoreMap[i] = score
		//			scoreSum += score
		//			//fmt.Println("Best weight for network", i, "is", bestWeight, "with score", score, "(best result", bestResult, ", complexity", complexity, ")")
		//		}
		//

		// sort the population index map -> score, by value (the scorE)
		//scoreList := sortByValue(scoreMap)

		//for _, pair := range scoreList {
		//	populationIndex := pair.Key
		//	score := pair.Value
		//	//fmt.Println("score", score, "population index", populationIndex)
		//}
		//if len(scoreList) == 0 {
		//	panic("NO SCORES!")
		//}

		//fmt.Println("Best score A: ", scoreList[0].Value)
		//fmt.Println("Best score B: ", scoreList[len(scoreList)-1].Value)
		//bestIndex := scoreList[0].Key
		//fmt.Println("Best population index: ", bestIndex)
		//bestNetwork = *(population[bestIndex])
		//fmt.Println(bestNetwork)

		// Now take the best networks and make mutated offspring.
		// Delete the worst networks.

		// For now, don't weight anything, just delete the bad half,
		// then add modified versions of the best 3 until the population is full.
		//
		// This method is probably buggy, since the score is the key for the indexes.
		//
		for networkIndex := 0; networkIndex < N; networkIndex++ {
			// Is this network in the best half?
			bestHalf := false
			for _, pair := range scoreList {
				score := pair.Value
				scoreIndex := pair.Key
				if scoreIndex == networkIndex {
					if score >= averageScore {
						bestHalf = true
						break
					}
				}
			}
			// If not in the best half, take a copy of the best network,
			// then modify it a bit (in a random way)
			if !bestHalf {
				// Take a deep copy, not just the the pointers
				newNetwork := bestNetwork.Copy()
				// Modify it a bit
				newNetwork.Modify()
				// Assign it to the population, replacing the low-scoring one
				population[networkIndex] = newNetwork
			}
			//fmt.Println(networkIndex, "is in the best half?", bestHalf)
		}

		//fmt.Println("Weight of best network:", bestNetwork.Weight)

		////		// Output a diagram of the best network
		////		err := bestNetwork.SaveDiagram("best" + strconv.Itoa(j) + ".svg")
		////		if err != nil {
		////			panic(err)
		////		}
	}

	err := bestNetwork.SaveDiagram("best.svg")
	if err != nil {
		panic(err)
	}

	// Now test the best network on 4 different inputs

	fmt.Println("Testing the up-detector:")
	upScore := bestNetwork.Evaluate(up)
	downScore := bestNetwork.Evaluate(down)
	leftScore := bestNetwork.Evaluate(left)
	rightScore := bestNetwork.Evaluate(right)

	if upScore > downScore && upScore > leftScore && upScore > rightScore {
		fmt.Println("Network success")
	} else {
		fmt.Println("Network failure")
	}
}
