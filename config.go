package wann

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
