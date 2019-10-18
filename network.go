package wann

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

// Initialize the random number generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Network is a collection of neurons, an output neuron pointer and a shared weight.
type Network struct {
	InputNodes []*Neuron
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
	// Pre-allocate room for n neurons and set the shared weight to 0.5
	net := &Network{make([]*Neuron, n), NewNeuron(), w}
	// Initialize n input nodes that all are inputs to the one output node.
	for i := 0; i < n; i++ {
		neuron := NewNeuron()
		net.InputNodes[i] = neuron
		// Make connections for all nodes where a random number between 0 and 1 are larger than r
		if rand.Float64() > r {
			net.OutputNode.AddInput(net.InputNodes[i])
		}
		if !net.OutputNode.HasInput(net.InputNodes[i]) {
			panic("EVERYTHING IS BROKEN")
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

// String creates a simple ASCII representation of the network
func (net *Network) String() string {
	var sb strings.Builder
	for i, n := range net.InputNodes {
		if net.OutputNode.HasInput(n) {
			if i == 0 {
				sb.WriteString("o---o\n")
			} else if i == (len(net.InputNodes) - 1) {
				sb.WriteString("o--/\n")
			} else {
				sb.WriteString("o--/|\n")
			}
		} else {
			if i == 0 {
				sb.WriteString("o   o\n")
			} else if i == (len(net.InputNodes) - 1) {
				sb.WriteString("o  /\n")
			} else {
				sb.WriteString("o  /|\n")
			}
		}
	}
	return sb.String()
}
