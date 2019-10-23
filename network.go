package wann

import (
	"errors"
	"fmt"
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
	OutputNode *Neuron
	Weight     float64
	nodeCount  uint
}

// NewNetwork creates a new minimal network with n input nodes and ratio of r connections
func NewNetwork(c *Config) *Network {
	n := c.Inputs
	r := c.ConnectionRatio
	w := c.SharedWeight
	if n <= 0 {
		return nil
	}
	// Pre-allocate room for n neurons and set the shared weight to the configured value.
	net := &Network{make([]*Neuron, n), NewRandomNeuron(), w, 0}

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

	// Store the total number of nodes so far
	net.nodeCount = uint(n) + 1

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
	net.nodeCount++
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
	sb.WriteString("Network\n")
	sb.WriteString("\tInput nodes: " + strconv.Itoa(len(net.InputNodes)) + "\n")
	sb.WriteString("\tConnections to output node: " + strconv.Itoa(len(net.OutputNode.InputNeurons)) + "\n")
	sb.WriteString("\tOutput neuron: " + fmt.Sprintf("%p", net.OutputNode) + "\n")
	for _, node := range net.All() {
		sb.WriteString("\t" + node.String() + "\n")
	}
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

// Copy will take a deep copy of the network struct
func (net *Network) Copy() *Network {
	var newNet Network
	for _, neuron := range net.InputNodes {
		newNet.InputNodes = append(newNet.InputNodes, neuron.Copy())
	}
	newOutputNeuron := *(net.OutputNode)
	newNet.OutputNode = &newOutputNeuron
	newNet.Weight = net.Weight
	newNet.nodeCount = net.nodeCount
	return &newNet
}

// All returns a slice with pointers to all nodes in this network
func (net *Network) All() []*Neuron {
	allNodes := net.InputNodes[:]
	// For each node that is connected to the output node
	net.ForEachConnected(func(neuron *Neuron, _ int) {
		for _, existingNode := range allNodes {
			if neuron == existingNode {
				// Skip this one by returning from the anonymouse function
				return
			}
		}
		// Add this node to the collection
		allNodes = append(allNodes, neuron)
	})
	// Return all nodes in this network
	return allNodes
}

// GetRandomNeuron will select a random neuron.
// This can be any node, including the output node.
func (net *Network) GetRandomNeuron() *Neuron {
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
		nodeA := net.GetRandomNeuron()
		nodeB := net.GetRandomNeuron()
		if nodeA != nodeB {
			// Insert a new node with a random activation function
			newNode := NewRandomNeuron()
			net.InsertNode(nodeA, nodeB, newNode)
		}
	case 1:
		//fmt.Println("Modifying the network using method 2 - add connection")
		nodeA := net.GetRandomNeuron()
		nodeB := net.GetRandomNeuron()
		if nodeA != nodeB {
			// Create a new connection
			net.AddConnection(nodeA, nodeB)
		}
	case 2:
		//fmt.Println("Modifying the network using method 3 - change activation")
		node := net.GetRandomNeuron()
		// Change the activation function
		node.RandomizeActivationFunction()
	}
}

// In checks if this neuron is in the given collection
func (node *Neuron) In(collection []*Neuron) bool {
	for _, existingNode := range collection {
		if node == existingNode {
			return true
		}
	}
	return false
}

func combine(a, b []*Neuron) []*Neuron {
	lena := len(a)
	lenb := len(b)
	// Allocate the exact size needed
	res := make([]*Neuron, lena+lenb)
	// Add the elements from a
	for i := 0; i < lena; i++ {
		res[i] = a[i]
	}
	// Add the elements from b
	for i := 0; i < lenb; i++ {
		res[i+lena] = b[i]
	}
	return res
}

// TODO: Never add an input to the input nodes!

// getAllNodes is a helper function for the recursive network traversal.
// Given the output node and the number 0, it will return a slice of all
// connected nodes, where the distance from the output node has been stored in
// node.distanceFromOutputNode.
func getAllNodes(node *Neuron, distanceFromFirstNode int) []*Neuron {
	allNodes := make([]*Neuron, 1)
	node.distanceFromOutputNode = distanceFromFirstNode
	allNodes[0] = node
	for _, inputNode := range node.InputNeurons {
		if !inputNode.In(allNodes) {
			allNodes = combine(allNodes, getAllNodes(inputNode, distanceFromFirstNode+1))
		}
	}
	return allNodes
}

// ForEachConnected will only go through nodes that are connected to the output node (directly or indirectly)
// Unconnected input nodes are not covered.
func (net *Network) ForEachConnected(f func(n *Neuron, distanceFromOutputNode int)) {
	// Start at the output node, traverse left towards the input nodes
	// The network has a counter for how many nodes has been added/removed, for quick memory allocation here
	allNodes := getAllNodes(net.OutputNode, 0)
	for _, node := range allNodes {
		f(node, node.distanceFromOutputNode)
	}
}
