package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// Config is a struct that is used when initializing new Network structs.
// The idea is that referring to fields by name is more explicit, and that it can
// be re-used in connection with having a configuration file, in the future.
type Config struct {
	// Number of input neurons (inputs per slice of floats in inputData in the Evolve function)
	Inputs int
	// When initializing a network, this is the propability that the node will be connected to the output node
	ConnectionRatio float64
	// SharedWeight is the weight that is shared by all nodes, since this is a Weight Agnostic Neural Network
	SharedWeight float64
	// How many generations to train for, at a maximum?
	Generations int
	// How large population sizes to use per generation?
	PopulationSize int
	// For how many generations should the training go on, without any improvement in the best score? Disabled if 0.
	MaxIterationsWithoutBestImprovement int
	// For how many generations should the training go on, without any improvement in the average score? Disabled if 0.
	MaxIterationsWithoutAverageImprovement int
	// Verbose?
	Verbose bool
}

// Evolve evolves a neural network, given a slice of training data and a slice of correct output values.
// Will overwrite config.Inputs.
func (config *Config) Evolve(inputData [][]float64, correctOutputMultipliers []float64) (*Network, error) {

	inputLength := len(inputData)

	if inputLength == 0 {
		return nil, errors.New("no input data")
	}

	if len(correctOutputMultipliers) == 1 && inputLength != 1 {
		// Assume the first slice of floats in the input data is the correct and that the rest are examples of being wrong
		for i := 1; i < inputLength; i++ {
			correctOutputMultipliers = append(correctOutputMultipliers, -1.0)
		}
	} else if inputLength != len(correctOutputMultipliers) {
		// Assume that the list of correct output multipliers should match the length of the float64 slices in inputData
		return nil, errors.New("the length of the input data and the slice of output multipliers differs")
	}

	config.Inputs = len(inputData[0])

	population := make([]*Network, config.PopulationSize)

	// Initialize the population
	for i := 0; i < config.PopulationSize; i++ {
		population[i] = NewNetwork(config)
	}

	var bestNetwork *Network

	// For each generation, evaluate and modify the networks

	bestScore := 0.0
	lastBestScore := 0.0
	noImprovementOfBestScoreCounter := 0

	//bestWeight := 0.0

	averageScore := 0.0
	lastAverageScore := 0.0
	noImprovementOfAverageScoreCounter := 0

	for j := 0; j < config.Generations; j++ {

		if config.Verbose {
			fmt.Println("------ generation " + strconv.Itoa(j) + ", population size " + strconv.Itoa(len(population)))
		}

		//bestWeight = 0.0
		bestNetwork = nil

		bestScore = 0.0
		averageScore = 0.0

		// For each weight, evaluate all networks
		//first := true
		w := rand.Float64()
		//for w := 0.0; w <= 1.0; w += 0.1 {

		scoreMap := make(map[int]float64)
		scoreSum := 0.0
		for i := 0; i < config.PopulationSize; i++ {
			net := population[i]

			// if net.AllNodes[net.OutputNode].InputNeuronsAreGood() {
			// 	fmt.Println("input neurons are good")
			// } else {
			// 	panic("input neurons are not good")
			// }

			net.SetWeight(w)
			result := 0.0
			for i := 0; i < len(inputData); i++ {
				result += net.Evaluate(inputData[i]) * correctOutputMultipliers[i]
			}
			score := result / net.Complexity()
			//if score <= 0 {
			//	score = -net.Complexity()
			//}
			scoreSum += score
			scoreMap[i] = score
		}

		// The scores for this weight
		scoreList := SortByValue(scoreMap)

		if bestNetwork == nil || scoreList[0].Value > bestScore {
			bestScore = scoreList[0].Value
			bestNetwork = population[scoreList[0].Key]
			//bestWeight = w
			//first = false
		}

		//}

		if bestNetwork == nil {
			panic("implementation error: no best network")
		}

		if config.Verbose {
			fmt.Printf("Best score: %f, using weight: %f\n", bestScore, w)
		}

		if bestScore == lastBestScore {
			noImprovementOfBestScoreCounter++
		}
		lastBestScore = bestScore

		// No better score for 20 generations? Stop evolving.
		if config.MaxIterationsWithoutBestImprovement > 0 && noImprovementOfBestScoreCounter > config.MaxIterationsWithoutBestImprovement {
			if config.Verbose {
				fmt.Println("No best score improvement for a while, done training.")
			}
			break
		}

		lastAverageScore = averageScore
		averageScore = scoreSum / float64(config.PopulationSize)
		if averageScore == lastAverageScore {
			noImprovementOfAverageScoreCounter++
		}

		// No better score for 20 generations? Stop evolving.
		if config.MaxIterationsWithoutAverageImprovement > 0 && noImprovementOfAverageScoreCounter > config.MaxIterationsWithoutAverageImprovement {
			if config.Verbose {
				fmt.Println("No average score improvement for a while, done training.")
			}
			break
		}

		if config.Verbose {
			fmt.Println("Average score:", averageScore)
		}

		// The scores for this weight
		scoreList = SortByValue(scoreMap)

		// Now take the best networks and make mutated offspring.
		// Delete the worst networks.

		for networkIndex := 0; networkIndex < config.PopulationSize; networkIndex++ {
			networkScore := scoreMap[networkIndex]
			// Is this network in the best half?
			bestHalf := networkScore >= averageScore
			// If the average score is 0, then modify an arbitrary half of the population
			if averageScore == 0 {
				if networkIndex > (config.PopulationSize / 2) {
					bestHalf = false
				}
			}
			// If not in the best half, take a copy of the best network,
			// then modify it a bit (in a random way)
			if !bestHalf {
				// Take a proper copy, not just the the pointers, because the nodes will be changed
				var newNetwork Network = *bestNetwork
				// Update all the Network pointers
				for i := range newNetwork.AllNodes {
					newNetwork.AllNodes[i].Net = &newNetwork
				}
				// Modify it a bit, with the maximum number of iterations being the same as the network size
				newNetwork.Modify(len(newNetwork.AllNodes))
				// Assign it to the population, replacing the low-scoring one
				population[networkIndex] = &newNetwork
			}
			//fmt.Println(networkIndex, "is in the best half?", bestHalf)
		}

	}
	if bestNetwork == nil {
		return nil, errors.New("the best network is nil")
	}
	return bestNetwork, nil
}
