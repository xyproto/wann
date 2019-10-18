package wann

// Config is a struct that is used when initializing new Network structs.
// The idea is that referring to fields by name is more explicit, and that it can
// be re-used in connection with having a configuration file, in the future.
type Config struct {
	// Number of input neurons
	Inputs int
	// When initializing a network, this is the propability that the node will be connected to the output node
	ConnectionRatio float64
	// SharedWeight is the weight that is shared by all nodes, since this is a Weight Agnostic Neural Network
	SharedWeight float64
}
