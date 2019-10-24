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
	net := &Network{make([]*Neuron, n), NewRandomNeuron(), w}

	// Initialize n input nodes that all are inputs to the one output node.
	for i := 0; i < n; i++ {
		net.InputNodes[i] = NewRandomNeuron()
		// Make connections for all nodes where a random number between 0 and 1 are larger than r
		if rand.Float64() > r {
			if err := net.OutputNode.AddInput(net.InputNodes[i]); err != nil {
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
// Assumes that a is the leftmost node and the b is the rightmost node.
func (net *Network) InsertNode(a, b, newNode *Neuron) error {
	// This is done by first checking that a is an input node to b,
	// then setting newNode to be an input node to b,
	// then setting a to be an input node to a.
	if a == b {
		return errors.New("the a and b nodes are the same")
	}
	// Sort the nodes by where they place in the diagram
	a, b = net.LeftRight(a, b)
	if b.In(net.InputNodes) {
		return errors.New("node b is a special input node")
	}
	if b == net.OutputNode {
		// this is fine
	}
	if a.In(net.InputNodes) {
		// this is fine
	}
	if a == net.OutputNode {
		// If now, after swapping, a is an output node, return with an error
		return errors.New("the leftmost node is an output node")
	}
	// b already has a as an input (a -> b)
	if b.HasInput(a) {
		// Remove the old connection
		if err := b.RemoveInput(a); err != nil {
			return errors.New("error in InsertNode b.RemoveInput(a): " + err.Error())
		}
	}
	// Connect the new node to b
	if err := b.AddInput(newNode); err != nil {
		return errors.New("error in InsertNode b.AddInput(newNode): " + err.Error())
	}
	// Connect a to the new node
	if err := newNode.AddInput(a); err != nil {
		return errors.New("error in InsertNode newNode.AddInput(a): " + err.Error())
	}
	// The situation should now be: a -> newNode -> b
	return nil
}

// AddConnection adds a connection from a to b
func (net *Network) AddConnection(a, b *Neuron) error {
	if a == b {
		return errors.New("can't connect to self")
	}
	// Sort the nodes by where they place in the diagram
	a, b = net.LeftRight(a, b)
	if a == net.OutputNode {
		// Swap a and b
		tmp := a
		b = a
		a = tmp
	}
	if a == net.OutputNode {
		// If now, after swapping, a is an output node, return with an error
		return errors.New("will not insert a node between the output node and another node")
	}
	if a.distanceFromOutputNode > b.distanceFromOutputNode {
		// Swap a and b
		tmp := a
		b = a
		a = tmp
	}
	if b.In(net.InputNodes) {
		return errors.New("b is an input node")
	}
	//if b.Value != nil {
	//return errors.New("b is a value node/input node"
	//}
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
	inputLength := len(inputValues)
	for i, n := range net.InputNodes {
		if i < inputLength {
			n.SetValue(inputValues[i])
		}
	}
	return net.OutputNode.evaluate(net.Weight)
}

// Evaluate2 will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
// A shared weight can be given. An error might be returned.
func (net *Network) Evaluate2(inputValues []float64) (float64, error) {
	inputLength := len(inputValues)
	if inputLength > len(net.InputNodes) {
		return 0.0, errors.New("Too many input values compared to input nodes")
	}
	for i, n := range net.InputNodes {
		if i < inputLength {
			n.SetValue(inputValues[i])
		}
	}
	return net.OutputNode.evaluate(net.Weight), nil
}

// SetWeight will set a shared weight for the entire network
func (net *Network) SetWeight(weight float64) {
	net.Weight = weight
}

// Complexity measures the network complexity
func (net *Network) Complexity() float64 {
	// Just return the node count, for now
	// TODO: Score the complexity of the various activation functions
	// TODO: Add complexity for each connected node
	return float64(len(net.InputNodes))
}

// // Copy will take a deep copy of the network struct
// func (net *Network) Copy() *Network {
// 	var newNet Network
// 	for _, neuron := range net.InputNodes {
// 		newNet.InputNodes = append(newNet.InputNodes, neuron.Copy())
// 	}
// 	newOutputNeuron := *(net.OutputNode)
// 	newNet.OutputNode = &newOutputNeuron
// 	newNet.Weight = net.Weight
// 	return &newNet
// }

// LeftRight returns two neurons, such that the first on is the one that is
// most to the left (towards the input neurons) and the second one is most to
// the right (towards the output neuron). Assumes that a and b are not equal.
func (net *Network) LeftRight(a, b *Neuron) (left *Neuron, right *Neuron) {
	if a.In(net.InputNodes) {
		left = a
		right = b
		return
	}
	if b.In(net.InputNodes) {
		left = b
		right = a
		return
	}
	if a == net.OutputNode {
		left = b
		right = a
		return
	}
	if b == net.OutputNode {
		left = a
		right = b
		return
	}
	if a.distanceFromOutputNode <= b.distanceFromOutputNode {
		left = b
		right = a
		return
	}
	left = a
	right = b
	return
}

type neuronList []*Neuron

func (neurons neuronList) Copy() []*Neuron {
	newList := make([]Neuron, len(neurons))
	for i, neuron := range neurons {
		newList[i] = *neuron
	}
	newList2 := make([]*Neuron, len(neurons))
	for i, neuron := range newList2 {
		newList2[i] = neuron
	}
	return newList2
}

// All returns a slice with pointers to all nodes in this network
func (net *Network) All() []*Neuron {
	allNodes := make([]*Neuron, 0)
	// For each node that is connected to the output node
	net.ForEachConnected(func(node *Neuron, _ int) {
		if !node.In(allNodes) {
			allNodes = append(allNodes, node)
		}
	})
	// Return all nodes in this network
	return allNodes
}

// GetRandomNeuron will select a random neuron.
// This can be any node, including the output node.
func (net *Network) GetRandomNeuron() *Neuron {
	allNeurons := net.All()
	chosenIndex := rand.Intn(len(allNeurons))
	if chosenIndex < 0 || chosenIndex >= len(allNeurons) {
		panic("implementation error: the chosen Index is invalid")
	}
	chosenNeuron := allNeurons[chosenIndex]
	if chosenNeuron == nil {
		panic("implementation error: the chosen neuron is nil")
	}
	return chosenNeuron
}

// Modify this network a bit
func (net *Network) Modify() {
	// Use method 0, 1 or 2
	method := rand.Intn(3) // up to and not including 3
	// TODO: Perform a modfification, using one of the three methods outlined in the paper
	switch method {
	case 0:
		//fmt.Println("Modifying the network using method 1 - insert node")
		nodeA, nodeB, newNode := net.GetRandomNeuron(), net.GetRandomNeuron(), NewRandomNeuron()
		// A bit risky, time-wise, but continue finding random neurons until they work out
		// Insert a new node with a random activation function
		for net.InsertNode(nodeA, nodeB, newNode) != nil {
			nodeA, nodeB, newNode = net.GetRandomNeuron(), net.GetRandomNeuron(), NewRandomNeuron()
		}
	case 1:
		//fmt.Println("Modifying the network using method 2 - add connection")

		nodeA, nodeB := net.GetRandomNeuron(), net.GetRandomNeuron()
		// A bit risky, time-wise, but continue finding random neurons until they work out
		// Create a new connection
		for net.AddConnection(nodeA, nodeB) != nil {
			nodeA, nodeB = net.GetRandomNeuron(), net.GetRandomNeuron()
		}
	case 2:
		//fmt.Println("Modifying the network using method 3 - change activation")
		// Change the activation function
		net.GetRandomNeuron().RandomizeActivationFunction()
	default:
		panic("implementation error: invalid method number: " + strconv.Itoa(method))
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
		if inputNode == node {
			panic("implementation error: node is input node to self")
		}
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
