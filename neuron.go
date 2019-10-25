package wann

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/xyproto/af"
)

// Neuron is a list of input-neurons, and an activation function.
type Neuron struct {
	Net                    *Network
	InputNeurons           []NeuronIndex // pointers to other neurons
	ActivationFunction     func(float64) float64
	Value                  *float64
	distanceFromOutputNode int // Used when traversing nodes and drawing diagrams
	neuronIndex            NeuronIndex
}

// NewNeuron creates a new Neuron
func (net *Network) NewNeuron() *Neuron {
	// Pre-allocate room for 64 connections and use Linear as the default activation function
	neuron := Neuron{Net: net, InputNeurons: make([]NeuronIndex, 0, 4), ActivationFunction: af.Linear}
	neuron.neuronIndex = NeuronIndex(len(net.AllNodes))
	net.AllNodes = append(net.AllNodes, neuron)
	return &neuron
}

// NewRandomNeuron creates a new *Neuron, with a randomly chosen activation function
func (net *Network) NewRandomNeuron() *Neuron {
	n := net.NewNeuron()
	n.RandomizeActivationFunction()
	return n
}

// RandomizeActivationFunction will choose a random activation function for this neuron
func (neuron *Neuron) RandomizeActivationFunction() {
	chosenIndex := rand.Intn(len(ActivationFunctions))
	neuron.ActivationFunction = ActivationFunctions[chosenIndex]
}

// SetValue can be used for setting a value for this neuron instead of using input neutrons.
// This changes how the Evaluation function behaves.
func (neuron *Neuron) SetValue(x float64) {
	neuron.Value = &x
}

// HasInput checks if the given neuron is an input neuron to this one
func (neuron *Neuron) HasInput(e NeuronIndex) bool {
	for _, ni := range neuron.InputNeurons {
		if ni == e {
			return true
		}
	}
	return false
}

// FindInput checks if the given neuron is an input neuron to this one,
// and also returns the index to InputNeurons, if found.
func (neuron *Neuron) FindInput(e NeuronIndex) (int, bool) {
	for i, n := range neuron.InputNeurons {
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
func (neuron *Neuron) AddInput(e NeuronIndex) error {
	if neuron.HasInput(e) {
		return errors.New("neuron already exists")
	}
	if neuron.Is(e) {
		return errors.New("adding a neuron as input to itself")
	}
	neuron.InputNeurons = append(neuron.InputNeurons, e)
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
		neuron.InputNeurons = append(neuron.InputNeurons[:i], neuron.InputNeurons[i+1:]...)
		return nil
	}
	return errors.New("neuron does not exist")
}

// Index finds the NeuronIndex for this node, if available
func (net *Network) Index(neuron *Neuron) NeuronIndex {
	return neuron.neuronIndex
	//for i := range net.AllNodes {
	//	if neuron.Is(NeuronIndex(i)) {
	//		return NeuronIndex(i), nil
	//	}
	//}
	//return NeuronIndex(-1), errors.New("neuron not found")
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

// String will return a string containing both the pointer address and the number of input neurons
func (neuron *Neuron) String() string {
	return fmt.Sprintf("node [%d] with %d inputs", neuron.neuronIndex, len(neuron.InputNeurons))
	// inputCount := len(neuron.InputNeurons)
	// switch inputCount {
	// case 0:
	// 	return fmt.Sprintf("Neuron [%d].", neuron.neuronIndex)
	// case 1:
	// 	return fmt.Sprintf("Neuron [%d] has 1 input: %d", neuron.neuronIndex, neuron.InputNeurons[0])
	// default:
	// 	var sb strings.Builder
	// 	sb.WriteString(fmt.Sprintf("Neuron [%d] has %d inputs:", neuron.neuronIndex, len(neuron.InputNeurons)))
	// 	for _, inputNeuronIndex := range neuron.InputNeurons {
	// 		inputNeuron := neuron.Net.AllNodes[inputNeuronIndex]
	// 		sb.WriteString("\n\t" + inputNeuron.String())
	// 	}
	// 	return sb.String()
	// }
}

// AnyInputConnectionsOrPanic will check if there are any input nodes
// in the network, or else just panic.
func (net *Network) AnyInputConnectionsOrPanic() {
	// DEBUG
	maxLengthInputNode := 0
	for _, node := range net.AllNodes {
		if len(node.InputNeurons) > maxLengthInputNode {
			maxLengthInputNode = len(node.InputNeurons)
		}
	}
	if maxLengthInputNode == 0 {
		panic("no input nodes in the entire network")
	}
}

// evaluate will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
func (neuron *Neuron) evaluate(weight float64, maxEvaluationLoops *int) (float64, bool) {
	//fmt.Println("Evaluate. Countdown: ", *maxEvaluationLoops)
	if *maxEvaluationLoops <= 0 {
		return 0.0, true
	}
	// Assume this is the Output neuron, recursively evaluating the result
	// For each input neuron, evaluate them
	summed := 0.0
	counter := 0

	for _, inputNeuronIndex := range neuron.InputNeurons {
		// Let each input neuron do its own evauluation, using the given weight
		(*maxEvaluationLoops)--
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
	return neuron.ActivationFunction(summed / float64(counter)), false
}

// Copy takes a deep copy of this neuron
func (neuron *Neuron) Copy() *Neuron {
	var newNeuron Neuron
	// for _, inputNeuron := range neuron.InputNeurons {
	// 	if inputNeuron == neuron {
	// 		newNeuron.InputNeurons = append(newNeuron.InputNeurons, inputNeuron.Copy())
	// 	} else {
	// 		// This neuron is an input to itself!? okay. Don't make a copy.
	// 		newNeuron.InputNeurons = append(newNeuron.InputNeurons, inputNeuron)
	// 	}
	// }
	newNeuron.InputNeurons = neuron.InputNeurons
	newNeuron.ActivationFunction = neuron.ActivationFunction
	if neuron.Value != nil {
		v := *neuron.Value
		newNeuron.Value = &v
	}
	return &newNeuron
}
