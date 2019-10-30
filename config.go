package wann

// Config is a struct that is used when initializing new Network structs.
// The idea is that referring to fields by name is more explicit, and that it can
// be re-used in connection with having a configuration file, in the future.
type Config struct {
	// Number of input neurons (inputs per slice of floats in inputData in the Evolve function)
	inputs int
	// When initializing a network, this is the propability that the node will be connected to the output node
	InitialConnectionRatio float64
	// sharedWeight is the weight that is shared by all nodes, since this is a Weight Agnostic Neural Network
	sharedWeight float64
	// How many generations to train for, at a maximum?
	Generations int
	// How large population sizes to use per generation?
	PopulationSize int
	// For how many generations should the training go on, without any improvement in the best score? Disabled if 0.
	MaxIterationsWithoutBestImprovement int
	// RandomSeed, for initializing the random number generator. The current time is used for the seed if this is set to 0.
	RandomSeed int64
	// Verbose output
	Verbose bool
	// Has the pseudo-random number generator been seeded and the activation function complexity been estimated yet?
	initialized bool
}
