package wann

import (
	"errors"
	"fmt"
	"math/rand"
)

// Neuron is a list of input-neurons, and an activation function.
type Neuron struct {
	Net                    *Network
	InputNodes             []NeuronIndex // pointers to other neurons
	ActivationFunction     ActivationFunctionIndex
	Value                  *float64
	distanceFromOutputNode int // Used when traversing nodes and drawing diagrams
	neuronIndex            NeuronIndex
}

// NewBlankNeuron creates a new Neuron, with the Step activation function as the default
func (net *Network) NewBlankNeuron() (*Neuron, NeuronIndex) {
	// Pre-allocate room for 16 connections and use Linear as the default activation function
	neuron := Neuron{Net: net, InputNodes: make([]NeuronIndex, 0, 16), ActivationFunction: Swish}
	neuron.neuronIndex = NeuronIndex(len(net.AllNodes))
	net.AllNodes = append(net.AllNodes, neuron)
	return &neuron, neuron.neuronIndex
}

// NewNeuron creates a new *Neuron, with a randomly chosen activation function
func (net *Network) NewNeuron() (*Neuron, NeuronIndex) {
	chosenActivationFunctionIndex := ActivationFunctionIndex(rand.Intn(len(ActivationFunctions)))
	inputNodes := make([]NeuronIndex, 0, 16)
	neuron := Neuron{
		Net:                net,
		InputNodes:         inputNodes,
		ActivationFunction: chosenActivationFunctionIndex,
	}
	// The length of net.AllNodes is what will be the last index
	neuronIndex := NeuronIndex(len(net.AllNodes))
	// Assign the neuron index in the net to the neuron
	neuron.neuronIndex = neuronIndex
	// Add this neuron to the net
	net.AllNodes = append(net.AllNodes, neuron)
	return &neuron, neuronIndex
}

// NewUnconnectedNeuron returns a new unconnected neuron with neuronIndex -1 and net pointer set to nil
func NewUnconnectedNeuron() *Neuron {
	// Pre-allocate room for 16 connections and use Linear as the default activation function
	neuron := Neuron{Net: nil, InputNodes: make([]NeuronIndex, 0, 16), ActivationFunction: Linear}
	neuron.neuronIndex = -1
	return &neuron
}

// Connect this neuron to a network, overwriting any existing connections.
// This will also clear any input nodes to this neuron, since the net is different.
// TODO: Find the input nodes from the neuron.Net, save those and re-assign if there are matches?
func (neuron *Neuron) Connect(net *Network) {
	neuron.InputNodes = []NeuronIndex{}
	neuron.Net = net
	for ni := range net.AllNodes {
		// Check if this network already has a pointer to this neuron
		if &net.AllNodes[ni] == neuron {
			// Yes, assign the index
			neuron.neuronIndex = NeuronIndex(ni)
			// All good, bail
			return
		}
	}
	// The neuron was not found in the network
	// Find what will be the last index in net.AllNodes
	neuronIndex := len(net.AllNodes)
	// Add this neuron to the network
	net.AllNodes = append(net.AllNodes, *neuron)
	// Assign the index
	net.AllNodes[neuronIndex].neuronIndex = NeuronIndex(neuronIndex)
}

// RandomizeActivationFunction will choose a random activation function for this neuron
func (neuron *Neuron) RandomizeActivationFunction() {
	chosenActivationFunctionIndex := ActivationFunctionIndex(rand.Intn(len(ActivationFunctions)))
	neuron.ActivationFunction = chosenActivationFunctionIndex
}

// SetValue can be used for setting a value for this neuron instead of using input neutrons.
// This changes how the Evaluation function behaves.
func (neuron *Neuron) SetValue(x float64) {
	neuron.Value = &x
}

// HasInput checks if the given neuron is an input neuron to this one
func (neuron *Neuron) HasInput(e NeuronIndex) bool {
	for _, ni := range neuron.InputNodes {
		if ni == e {
			return true
		}
	}
	return false
}

// FindInput checks if the given neuron is an input neuron to this one,
// and also returns the index to InputNeurons, if found.
func (neuron *Neuron) FindInput(e NeuronIndex) (int, bool) {
	for i, n := range neuron.InputNodes {
		if n == e {
			return i, true
		}
	}
	return -1, false
}

// Is check if the given NeuronIndex points to this neuron
func (neuron *Neuron) Is(e NeuronIndex) bool {
	return neuron.neuronIndex == e
}

// AddInput will add an input neuron
func (neuron *Neuron) AddInput(ni NeuronIndex) error {
	if neuron.Is(ni) {
		return errors.New("adding a neuron as input to itself")
	}
	if neuron.HasInput(ni) {
		return errors.New("neuron already exists")
	}
	neuron.InputNodes = append(neuron.InputNodes, ni)

	return nil
}

// AddInputNeuron both adds a neuron to this network (if needed) and also
// adds its neuron index to the neuron.InputNeurons
func (neuron *Neuron) AddInputNeuron(n *Neuron) error {
	// If n.neuronIndex is known to this network, just add the NeuronIndex to neuron.InputNeurons
	if neuron.Net.Exists(n.neuronIndex) {
		return neuron.AddInput(n.neuronIndex)
	}
	// If not, add this neuron to the network first
	node := *n
	node.neuronIndex = NeuronIndex(len(neuron.Net.AllNodes))
	neuron.Net.AllNodes = append(neuron.Net.AllNodes, node)
	return neuron.AddInput(n.neuronIndex)
}

// RemoveInput will remove an input neuron
func (neuron *Neuron) RemoveInput(e NeuronIndex) error {
	if i, found := neuron.FindInput(e); found {
		// Found it, remove the neuron at index i
		neuron.InputNodes = append(neuron.InputNodes[:i], neuron.InputNodes[i+1:]...)
		return nil
	}
	return errors.New("neuron does not exist")
}

// Exists checks if the given NeuronIndex exists in this Network
func (net *Network) Exists(ni NeuronIndex) bool {
	for i := range net.AllNodes {
		neuronIndex := NeuronIndex(i)
		if neuronIndex == ni {
			return true
		}
	}
	return false
}

// InputNeuronsAreGood checks if all input neurons of this neuron exists in neuron.Net
func (neuron *Neuron) InputNeuronsAreGood() bool {
	for _, inputNeuronIndex := range neuron.InputNodes {
		if !neuron.Net.Exists(inputNeuronIndex) {
			return false
		}
	}
	return true
}

// evaluate will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
// returns true if the maximum number of evaluation loops is reached
func (neuron *Neuron) evaluate(weight float64, maxEvaluationLoops *int) (float64, bool) {
	if *maxEvaluationLoops <= 0 {
		return 0.0, true
	}
	// Assume this is the Output neuron, recursively evaluating the result
	// For each input neuron, evaluate them
	summed := 0.0
	counter := 0

	for _, inputNeuronIndex := range neuron.InputNodes {
		// Let each input neuron do its own evauluation, using the given weight
		(*maxEvaluationLoops)--
		// TODO: Figure out exactly why this one kicks in (and if it matters)
		//       It only seems to kick in during "go test" and not in evolve/main.go
		if int(inputNeuronIndex) >= len(neuron.Net.AllNodes) {
			continue
			//panic("TOO HIGH INPUT NEURON INDEX")
		}
		result, stopNow := neuron.Net.AllNodes[inputNeuronIndex].evaluate(weight, maxEvaluationLoops)
		summed += result * weight
		counter++
		if stopNow || (*maxEvaluationLoops < 0) {
			break
		}
	}
	// No input neurons. Use the .Value field if it's not nil.
	if counter == 0 && neuron.Value != nil {
		return *(neuron.Value), false
	}
	// Return the averaged sum, or 0
	if counter == 0 {
		return 0.0, false
	}
	f := neuron.GetActivationFunction()
	// Run the average input through the activation function
	return f(summed / float64(counter)), false
}

// GetActivationFunction returns the activation function for this neuron
func (neuron *Neuron) GetActivationFunction() func(float64) float64 {
	return ActivationFunctions[neuron.ActivationFunction]
}

// In checks if this neuron is in the given collection
func (neuron *Neuron) In(collection []NeuronIndex) bool {
	for _, existingNodeIndex := range collection {
		if neuron.Is(existingNodeIndex) {
			return true
		}
	}
	return false
}

// IsInput returns true if this is an input node or not
// Returns false if nil
func (neuron *Neuron) IsInput() bool {
	if neuron.Net == nil {
		return false

	}
	return neuron.Net.IsInput(neuron.neuronIndex)
}

// IsOutput returns true if this is an output node or not
// Returns false if nil
func (neuron *Neuron) IsOutput() bool {
	if neuron.Net == nil {
		return false
	}
	return neuron.Net.OutputNode == neuron.neuronIndex
}

// Copy a Neuron to a new Neuron, and assign the pointer to the given network to .Net
func (neuron Neuron) Copy(net *Network) Neuron {
	var newNeuron Neuron
	newNeuron.Net = net
	newNeuron.InputNodes = neuron.InputNodes
	newNeuron.ActivationFunction = neuron.ActivationFunction
	newNeuron.Value = neuron.Value
	newNeuron.distanceFromOutputNode = neuron.distanceFromOutputNode
	newNeuron.neuronIndex = neuron.neuronIndex
	return newNeuron
}

// String will return a string containing both the pointer address and the number of input neurons
func (neuron *Neuron) String() string {
	nodeType := "       Node"
	if neuron.IsInput() {
		nodeType = " Input node"
	} else if neuron.IsOutput() {
		nodeType = "Output node"
	}
	return fmt.Sprintf("%s ID %d has these input connections: %v", nodeType, neuron.neuronIndex, neuron.InputNodes)
}
