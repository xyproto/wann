package wann

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Network is a collection of nodes, an output node and a shared weight.
type Network struct {
	InputNodes []*Neuron
	ExtraNodes []*Neuron
	OutputNode *Neuron
	Weight     float64
}

// NewNetwork creates a new minimal network with n input nodes and ratio of r connections
func NewNetwork(c *Config) *Network {
	n := c.Inputs
	r := c.ConnectionRatio
	w := c.SharedWeight
	if n <= 0 {
		return nil
	}
	// Pre-allocate room for n neurons and set the shared weight to the configured value
	net := &Network{make([]*Neuron, n), make([]*Neuron, 0, 64), NewRandomNeuron(), w}

	// Initialize n input nodes that all are inputs to the one output node.
	for i := 0; i < n; i++ {
		net.InputNodes[i] = NewRandomNeuron()
		// Make connections for all nodes where a random number between 0 and 1 are larger than r
		if rand.Float64() > r {
			err := net.OutputNode.AddInput(net.InputNodes[i])
			if err != nil {
				panic(err)
			}
		}
	}
	return net
}

//
// Operators for searching the space of network topologies
//

// InsertNode takes two neurons and inserts a third neuron between them
func (net *Network) InsertNode(a, b, newNode *Neuron) error {
	// This is done by first checking that a is an input node to b,
	// then setting newNode to be an input node to b,
	// then setting a to be an input node to a.
	if !b.HasInput(a) {
		return errors.New("can not insert node: a is not an input neuron to b")
	}
	err := b.RemoveInput(a)
	if err != nil {
		return errors.New("can not insert node: " + err.Error())
	}
	err = b.AddInput(newNode)
	if err != nil {
		return errors.New("can not insert node: " + err.Error())
	}
	err = newNode.AddInput(a)
	if err != nil {
		return errors.New("can not insert node: " + err.Error())
	}
	net.ExtraNodes = append(net.ExtraNodes, newNode)
	return nil
}

// AddConnection adds a connection from a to b
func (net *Network) AddConnection(a, b *Neuron) error {
	return b.AddInput(a)
}

// ChangeActivationFunction changes the activation function for a given node
func (net *Network) ChangeActivationFunction(n *Neuron, f func(float64) float64) {
	n.ActivationFunction = f
}

// String creates a simple and not very useful ASCII representation of the input nodes and the output node.
// Nodes that are not input nodes are skipped.
// Input nodes that are not connected directly to the output node are drawn as non-connected,
// even if they are connected via another node.
func (net *Network) String() string {
	var sb strings.Builder
	//--- Network ---
	//Input nodes: 5
	//Connections to output node: 5
	//---------------
	sb.WriteString("--- Network ---\n")
	sb.WriteString("Input nodes: " + strconv.Itoa(len(net.InputNodes)) + "\n")
	sb.WriteString("Connections to output node: " + strconv.Itoa(len(net.OutputNode.InputNeurons)) + "\n")
	sb.WriteString("---------------")
	return sb.String()
}

// Evaluate will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
// A shared weight can be given.
func (net *Network) Evaluate(inputValues []float64) float64 {
	for i, n := range net.InputNodes {
		if i < len(inputValues) {
			n.SetValue(inputValues[i])
		}
	}
	return net.OutputNode.evaluate(net.Weight)
}

// SetWeight will set a shared weight for the entire network
func (net *Network) SetWeight(weight float64) {
	net.Weight = weight
}

// Complexity measures the network complexity
func (net *Network) Complexity() float64 {
	// Just return the node count, for now
	return float64(len(net.InputNodes))
}

// Take a deep copy of the network struct
func (net *Network) Copy() *Network {
	var newNet Network
	for _, neuron := range net.InputNodes {
		newNet.InputNodes = append(newNet.InputNodes, neuron.Copy())
	}
	for _, neuron := range net.ExtraNodes {
		newNet.ExtraNodes = append(newNet.ExtraNodes, neuron.Copy())
	}
	newOutputNeuron := *(net.OutputNode)
	newNet.OutputNode = &newOutputNeuron
	newNet.Weight = net.Weight
	return &newNet
}

// All returns a slice with pointers to all nodes in this network
func (net *Network) All() []*Neuron {
	allNodes := make([]*Neuron, 0, len(net.InputNodes)+len(net.ExtraNodes)+1)
	allNodes = append(allNodes, net.InputNodes...)
	allNodes = append(allNodes, net.ExtraNodes...)
	allNodes = append(allNodes, net.OutputNode)
	return allNodes
}

// Start with the output node and traverse the net to gather a list of neurons,
// then choose one at random.
func (net *Network) FindRandomNeuron() *Neuron {
	allNeurons := net.All()
	chosenIndex := rand.Intn(len(allNeurons))
	return allNeurons[chosenIndex]
}

// Modify this network a bit
func (net *Network) Modify() {
	// Use method 0, 1 or 2
	method := rand.Intn(3) // up to and not including 3
	// TODO: Perform a modfification, using one of the three methods outlined in the paper
	switch method {
	case 0:
		//fmt.Println("Modifying the network using method 1 - insert node")
		nodeA := net.FindRandomNeuron()
		nodeB := net.FindRandomNeuron()
		if nodeA != nodeB {
			// Insert a new node with a random activation function
			newNode := NewRandomNeuron()
			net.InsertNode(nodeA, nodeB, newNode)
		}
	case 1:
		//fmt.Println("Modifying the network using method 2 - add connection")
		nodeA := net.FindRandomNeuron()
		nodeB := net.FindRandomNeuron()
		if nodeA != nodeB {
			// Create a new connection
			net.AddConnection(nodeA, nodeB)
		}
	case 2:
		//fmt.Println("Modifying the network using method 3 - change activation")
		node := net.FindRandomNeuron()
		// Change the activation function
		node.RandomizeActivationFunction()
	}
}
