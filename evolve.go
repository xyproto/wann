package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// ScorePopulation evaluates a population, given a slice of input numbers.
// It returns a map with scores, together with the sum of scores.
func ScorePopulation(population []*Network, weight float64, inputData [][]float64, incorrectOutputMultipliers []float64) (map[int]float64, float64) {

	scoreMap := make(map[int]float64)
	scoreSum := 0.0

	for i := 0; i < len(population); i++ {
		net := population[i]

		if len(net.AllNodes[net.OutputNode].InputNodes) == 0 {
			// The output node has no input nodes, not great
			scoreMap[i] = 0.0
			continue
		}

		net.SetWeight(weight)

		// Evaluate all the input data examples for this network
		result := 0.0
		for i := 0; i < len(inputData); i++ {
			result += net.Evaluate(inputData[i]) * incorrectOutputMultipliers[i]
		}

		// The score is how well the network is doing, divided by the network complexity rating
		score := result / net.Complexity()

		scoreSum += score
		scoreMap[i] = score
	}
	return scoreMap, scoreSum
}

// Modify the network using one of the three methods outlined in the paper:
// * Insert node
// * Add connection
// * Change activation function
func (net *Network) Modify(maxIterations int) {

	// Use method 0, 1 or 2
	method := rand.Intn(3) // up to and not including 3

	// Perform a modfification, using one of the three methods outlined in the paper
	switch method {
	case 0:
		// Insert a node, replacing a randomly chosen existing connection
		counter := 0
		for net.InsertRandomNode() == false {
			counter++
			if maxIterations > 0 && counter > maxIterations {
				break
			}
		}
	case 1:
		nodeA, nodeB := net.GetRandomNode(), net.GetRandomNode()
		// Continue finding random neurons until they work out or until maxIterations is reached
		// Create a new connection
		counter := 0
		for net.AddConnection(nodeA, nodeB) != nil {
			nodeA, nodeB = net.GetRandomNode(), net.GetRandomNode()
			counter++
			if maxIterations > 0 && counter > maxIterations {
				// Could not add a connection. The possibilities for connections might be saturated.
				return
			}
		}
	case 2:
		// Change the activation function to a randomly selected one
		net.RandomizeActivationFunctionForRandomNeuron()
	default:
		panic("implementation error: invalid method number: " + strconv.Itoa(method))
	}
}

// Complexity measures the network complexity
// Will return 1.0 at a minimum
func (net *Network) Complexity() float64 {

	// TODO: These two constants really affect the results. Place them in the Config struct instead.

	// How much should the function complexity matter in relation to the number of connected nodes?
	const functionComplexityMultiplier = 7.0

	// How much should the complexity score matter in relation to the network results, when scoring the network?
	const complexityMultiplier = 5.0

	sum := 0.0
	// Sum the complexity of all activation functions.
	// This penalizes both slow activation functions and
	// unconnected nodes.
	for _, n := range net.AllNodes {
		if n.Value == nil {
			sum += ComplexityEstimate[n.ActivationFunction] * functionComplexityMultiplier
		}
	}
	// The number of connected nodes should also carry some weight
	connectedNodes := float64(len(net.Connected()))
	// This must always be larger than 0, to avoid divide by zero later
	return (connectedNodes+sum)*complexityMultiplier + 1.0
}

// Evolve evolves a neural network, given a slice of training data and a slice of correct output values.
// Will overwrite config.Inputs.
func (config *Config) Evolve(inputData [][]float64, correctOutputMultipliers []float64) (*Network, error) {

	// TODO: If the config.initialConnectionRatio field is too low (0.0, for instance), then this function will fail.
	//       Return with an error if none of the networks in a population has any connections left, then get rid of the "no improvement counter".

	// Initialize, if needed
	if !config.initialized {
		config.Init()
	}

	const maxModificationInterationsWhenMutating = 10

	incorrectOutputMultipliers := make([]float64, len(correctOutputMultipliers))
	for i := range correctOutputMultipliers {
		// Convert from having 0..1 for meaning from incorrect to correct, to -1..1 to mean the same
		incorrectOutputMultipliers[i] = correctOutputMultipliers[i]*2.0 - 1.0
		// Convert from having 0..1 for meaning from incorrect to correct, to 1..0 to mean the same
		//incorrectOutputMultipliers[i] = -correctOutputMultipliers[i] + 1.0
	}

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

	config.inputs = len(inputData[0])

	population := make([]*Network, config.PopulationSize)

	// Initialize the population
	for i := 0; i < config.PopulationSize; i++ {
		n := NewNetwork(config)
		population[i] = &n
		population[i].UpdateNetworkPointers()
	}

	var (
		bestNetwork *Network

		// Keep track of the best scores
		bestScore     float64
		lastBestScore float64

		noImprovementCounter int // Counts how many times the best score has been stagnant

		// Keep track of the average scores
		averageScore float64

		// Keep track of the worst scores
		worstScore float64
	)

	if config.Verbose {
		fmt.Printf("Starting evolution with population size %d, for %d generations.\n", config.PopulationSize, config.Generations)
	}

	// For each generation, evaluate and modify the networks
	for j := 0; j < config.Generations; j++ {

		bestNetwork = nil

		// Initialize the scores with unlikely values
		// TODO: Use the first network in the population for initializing these instead
		first := true

		// Random weight from -2.0 to 2.0
		w := rand.Float64()

		// The scores for this generation (using a random shared weight within ScorePopulation).
		// CorrectOutputMultipliers gives weight to the "correct" or "wrong" results, with the same index as the inputData
		// Score each network in the population.
		scoreMap, scoreSum := ScorePopulation(population, w, inputData, incorrectOutputMultipliers)

		// Sort by score
		scoreList := SortByValue(scoreMap)

		// Handle the best score stats
		if first {
			lastBestScore = 0.0
			bestScore = scoreList[0].Value
			worstScore = scoreList[len(scoreList)-1].Value
			bestNetwork = population[scoreList[0].Key]
			bestNetwork.SetWeight(w)
			first = false
		} else {
			lastBestScore = bestScore
			if scoreList[0].Value > bestScore {
				bestScore = scoreList[0].Value
			}
		}
		if bestScore > lastBestScore {
			bestNetwork = population[scoreList[0].Key]
			bestNetwork.SetWeight(w)
			noImprovementCounter = 0
		} else {
			noImprovementCounter++
		}

		// Handle the average score stats
		averageScore = scoreSum / float64(config.PopulationSize)

		// Handle the worst score stats
		if scoreList[len(scoreList)-1].Value < worstScore {
			worstScore = scoreList[len(scoreList)-1].Value
		}

		if bestNetwork == nil {
			panic("implementation error: no best network")
		}

		if config.Verbose {
			fmt.Printf("[generation %d] worst score = %f, average score = %f, best score = %f\n", j, worstScore, averageScore, bestScore)
			//fmt.Printf("[generation %d] worst score = %f, average score = %f, best score = %f, no improvement counter for this generation = %d\n", j, worstScore, averageScore, bestScore, noImprovementCounter)
			if noImprovementCounter > 0 {
				fmt.Printf("No improvement in the best score for the last %d generations\n", noImprovementCounter)
			}
		}

		// Only keep the best 7%
		bestFractionCountdown := int(float64(len(population)) * 0.07)

		goodNetworks := make([]*Network, 0, bestFractionCountdown)

		// Now loop over all networks, sorted by score (descending order)
		// p.Key is the network index
		// p.Value is the network score
		for _, p := range scoreList {
			networkIndex := p.Key
			if bestFractionCountdown > 0 {
				bestFractionCountdown--
				// In the best third of the networks
				goodNetworks = append(goodNetworks, population[networkIndex])
				continue
			}
			// // If there has not been any improvement to the best score lately, randomize the bad half
			// if noImprovementCounter > 100 {
			// 	n := NewNetwork(config)
			// 	population[networkIndex] = &n
			// 	continue
			// }
			randomGoodNetwork := goodNetworks[rand.Intn(len(goodNetworks))]
			randomGoodNetworkCopy := randomGoodNetwork.Copy()
			randomGoodNetworkCopy.Modify(maxModificationInterationsWhenMutating)
			// Replace the "bad" network with the modified copy of a "good" one
			// It's important that this is a pointer to a Network and not
			// a bare Network, so that the node .Net pointers are correct.
			population[networkIndex] = randomGoodNetworkCopy
		}
		// if noImprovementCounter > 100 {
		// 	noImprovementCounter = 0
		// }
	}
	if config.Verbose {
		fmt.Printf("[all time best network, random weight  ] weight=%f score=%f\n", bestNetwork.Weight, bestScore)
	}

	// Now find the best weight for the best network, using a population of 1
	// and a step size of 0.0001 for the weight
	population = []*Network{bestNetwork}
	bestWeight := -2.0
	for w := -2.0; w <= 2.0; w += 0.0001 {
		scoreMap, _ := ScorePopulation(population, w, inputData, incorrectOutputMultipliers)
		// Handle the best score stats
		if scoreMap[0] > bestScore {
			bestScore = scoreMap[0]
			population[0].SetWeight(w)
			bestWeight = w
		}
	}

	// Check if the best network is nil, just in case
	if bestNetwork == nil {
		return nil, errors.New("the total best network is nil")
	}

	// Save the best weight for the network
	bestNetwork.SetWeight(bestWeight)

	if config.Verbose {
		fmt.Printf("[all time best network, optimal weight ] weight=%f best score=%f\n", bestNetwork.Weight, bestScore)
	}

	return bestNetwork, nil
}
