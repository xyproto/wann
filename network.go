package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// init will initialize the random number generator with the current time
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NeuronIndex is an index into the AllNodes slice
type NeuronIndex int

// Network is a collection of nodes, an output node and a shared weight.
type Network struct {
	AllNodes   []Neuron      // Storing the actual neurons
	InputNodes []NeuronIndex // Pointers to the input nodes
	OutputNode NeuronIndex   // Pointer to the output node
	Weight     float64       // Shared weight
}

// NewNetwork creates a new minimal network with n input nodes and ratio of r connections.
// Passing "nil" as an argument is supported.
func NewNetwork(cs ...*Config) Network {
	c := &Config{}
	// If a single non-nil *Config struct is given, use that
	if len(cs) == 1 && cs[0] != nil {
		c = cs[0]
	}
	n := c.Inputs
	r := c.ConnectionRatio
	w := c.SharedWeight
	// Create a new network that has one node, the output node
	outputNodeIndex := NeuronIndex(0)
	net := Network{make([]Neuron, 0, n+1), make([]NeuronIndex, n), outputNodeIndex, w}
	outputNode, outputNodeIndex := net.NewRandomNeuron()
	net.OutputNode = outputNodeIndex

	// Initialize n input nodes that all are inputs to the one output node.
	for i := 0; i < n; i++ {
		// Add a new input node

		_, nodeIndex := net.NewRandomNeuron()

		// Register the input node index in the input node NeuronIndex slice
		net.InputNodes[i] = nodeIndex

		// Make connections for all nodes where a random number between 0 and 1 are larger than r
		if r >= rand.Float64() {
			if err := outputNode.AddInput(nodeIndex); err != nil {
				panic(err)
			}
		}
	}

	// Store the modified output node
	net.AllNodes[outputNodeIndex] = *outputNode

	return net
}

// Get returns a pointer to a neuron, based on the given NeuronIndex
func (net *Network) Get(i NeuronIndex) *Neuron {
	return &(net.AllNodes[i])
}

// IsInput checks if the given node is an input node
func (net *Network) IsInput(ni NeuronIndex) bool {
	for _, inputNodeIndex := range net.InputNodes {
		if ni == inputNodeIndex {
			return true
		}
	}
	return false
}

//
// Operators for searching the space of network topologies
//

// InsertNode takes two neurons and inserts a third neuron between them
// Assumes that a is the leftmost node and the b is the rightmost node.
func (net *Network) InsertNode(a, b NeuronIndex, newNodeIndex NeuronIndex) error {
	// This is done by first checking that a is an input node to b,
	// then setting newNode to be an input node to b,
	// then setting a to be an input node to a.

	// TODO: When a neuron is inserted, the input index

	if a == b {
		return errors.New("the a and b nodes are the same")
	}
	// Sort the nodes by where they place in the diagram
	//fmt.Println("InsertNode: BEFORE LEFT RIGHT:", a, b)
	a, b = net.LeftRight(a, b)
	//fmt.Println("InsertNode: AFTER LEFT RIGHT:", a, b)
	if net.IsInput(b) {
		if net.IsInput(a) {
			return errors.New("both node a and b are special input nodes")
		}
		return errors.New("node b (but not a) is a special input node")
	}

	if b == net.OutputNode {
		// this is fine
		//fmt.Println("b is the output node")
	}
	if net.IsInput(a) {
		// this is fine
		//fmt.Println("a is an input node")
	}
	if a == net.OutputNode {
		// If now, after swapping, a is an output node, return with an error
		return errors.New("the leftmost node is an output node")
	}

	// b already has a as an input (a -> b)
	if net.AllNodes[b].HasInput(a) {
		// Remove the old connection
		if err := net.AllNodes[b].RemoveInput(a); err != nil {
			return errors.New("error in InsertNode b.RemoveInput(a): " + err.Error())
		}
	}

	// b already has newNodeIndex as an input (newIndex -> b)
	if net.AllNodes[b].HasInput(newNodeIndex) {
		// Remove the old connection
		if err := net.AllNodes[b].RemoveInput(a); err != nil {
			return errors.New("error in InsertNode b.RemoveInput(a): " + err.Error())
		}
	}

	//net.AllNodes[net.OutputNode].Net = net

	// Connect the new node to b
	if err := net.AllNodes[b].AddInput(newNodeIndex); err != nil {
		// This does not kick in, the problem must be in AddInput!
		return errors.New("error in InsertNode b.AddInput(newNode): " + err.Error())
	}

	newNode := net.Get(newNodeIndex)

	// Connect a to the new node
	if err := newNode.AddInput(a); err != nil {
		return errors.New("error in InsertNode newNode.AddInput(a): " + err.Error())
	}

	// The situation should now be: a -> newNode -> b
	return nil
}

// AddConnection adds a connection from a to b
func (net *Network) AddConnection(a, b NeuronIndex) error {
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
	if net.AllNodes[a].distanceFromOutputNode > net.AllNodes[b].distanceFromOutputNode {
		// Swap a and b
		tmp := a
		b = a
		a = tmp
	}
	if net.IsInput(b) {
		return errors.New("b is an input node")
	}
	//if b.Value != nil {
	//return errors.New("b is a value node/input node"
	//}
	return net.AllNodes[b].AddInput(a)
}

// ChangeActivationFunction changes the activation function for a given node
func (net *Network) ChangeActivationFunction(n *Neuron, f func(float64) float64) {
	n.ActivationFunction = f
}

// Evaluate will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
// A shared weight can be given.
func (net *Network) Evaluate(inputValues []float64) float64 {
	inputLength := len(inputValues)
	for i, nindex := range net.InputNodes {
		if i < inputLength {
			net.AllNodes[nindex].SetValue(inputValues[i])
		}
	}
	maxIterationCounter := len(net.AllNodes)
	result, _ := net.AllNodes[net.OutputNode].evaluate(net.Weight, &maxIterationCounter)
	return result
}

// SetWeight will set a shared weight for the entire network
func (net *Network) SetWeight(weight float64) {
	net.Weight = weight
}

// Complexity measures the network complexity
// Will return 1.0 at a minimum
func (net *Network) Complexity() float64 {
	// TODO: Score the complexity of the various activation functions
	// TODO: Score the amount of connections
	return float64(1.0 + len(net.All()))
}

// LeftRight returns two neurons, such that the first on is the one that is
// most to the left (towards the input neurons) and the second one is most to
// the right (towards the output neuron). Assumes that a and b are not equal.
func (net *Network) LeftRight(a, b NeuronIndex) (left NeuronIndex, right NeuronIndex) {
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
	if net.AllNodes[a].In(net.InputNodes) {
		left = a
		right = b
		return
	}
	if net.AllNodes[b].In(net.InputNodes) {
		left = b
		right = a
		return
	}
	if net.AllNodes[a].distanceFromOutputNode <= net.AllNodes[b].distanceFromOutputNode {
		left = b
		right = a
		return
	}
	left = a
	right = b
	return
}

// Depth returns the maximum connection distance from the output node
func (net *Network) Depth() int {
	maxDepth := 0
	for _, n := range net.All() {
		if n.distanceFromOutputNode > maxDepth {
			maxDepth = n.distanceFromOutputNode
		}
	}
	return maxDepth
}

func (net *Network) checkInputNeurons() {
	for _, n := range net.All() {
		n.checkInputNeurons()
	}
}

// Copy a Network to a new network
func (net Network) Copy() Network {
	var newNet Network
	newNet.AllNodes = make([]Neuron, len(net.AllNodes))
	for i, node := range net.AllNodes {
		newNet.AllNodes[i] = node.Copy(&newNet)
	}
	newNet.InputNodes = net.InputNodes
	newNet.OutputNode = net.OutputNode
	//newNet.checkInputNeurons()
	newNet.Weight = net.Weight
	return newNet
}

// All returns a slice with pointers to all nodes in this network
func (net *Network) All() []*Neuron {
	allNodes := make([]*Neuron, 0)
	for i := range net.AllNodes {
		allNodes = append(allNodes, &net.AllNodes[i])
	}
	// Return pointers to all nodes in this network
	return allNodes
}

// GetRandomNeuron will select a random neuron.
// This can be any node, including the output node.
func (net *Network) GetRandomNeuron() NeuronIndex {
	return NeuronIndex(rand.Intn(len(net.AllNodes)))
}

// GetRandomInputNode returns a random input node
func (net *Network) GetRandomInputNode() NeuronIndex {
	inputPosition := rand.Intn(len(net.InputNodes))
	inputNodeIndex := net.InputNodes[inputPosition]
	return inputNodeIndex
}

// In checks if this neuron is in the given collection
func (node *Neuron) In(collection []NeuronIndex) bool {
	for _, existingNodeIndex := range collection {
		if node.Is(existingNodeIndex) {
			return true
		}
	}
	return false
}

// Combine will combine two lists of indices
func Combine(a, b []NeuronIndex) []NeuronIndex {
	lena := len(a)
	lenb := len(b)
	// Allocate the exact size needed
	res := make([]NeuronIndex, lena+lenb)
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

// getAllNodes is a helper function for the recursive network traversal.
// Given the output node and the number 0, it will return a slice of all
// connected nodes, where the distance from the output node has been stored in
// node.distanceFromOutputNode.
func (net *Network) getAllNodes(nodeIndex NeuronIndex, distanceFromFirstNode int, alreadyHaveThese []NeuronIndex) []NeuronIndex {
	allNodes := make([]NeuronIndex, 0, len(net.AllNodes))
	node := net.AllNodes[nodeIndex]
	if nodeIndex != net.OutputNode {
		node.distanceFromOutputNode = distanceFromFirstNode
		net.AllNodes[nodeIndex] = node
	}
	if !node.In(alreadyHaveThese) {
		allNodes = append(allNodes, nodeIndex)
	}
	for _, inputNodeIndex := range node.InputNeurons {
		if node.Is(inputNodeIndex) {
			panic("implementation error: node is input node to self")
		}
		if int(inputNodeIndex) >= len(net.AllNodes) {
			continue
		}
		inputNode := net.AllNodes[inputNodeIndex]
		if !inputNode.In(allNodes) && !inputNode.In(alreadyHaveThese) {
			allNodes = Combine(allNodes, net.getAllNodes(inputNodeIndex, distanceFromFirstNode+1, append(allNodes, alreadyHaveThese...)))
		}
	}
	return allNodes
}

// ForEachConnected will only go through nodes that are connected to the output node (directly or indirectly)
// Unconnected input nodes are not covered.
func (net *Network) ForEachConnected(f func(n *Neuron, distanceFromOutputNode int)) {
	// Start at the output node, traverse left towards the input nodes
	// The network has a counter for how many nodes has been added/removed, for quick memory allocation here
	allNodes := net.getAllNodes(net.OutputNode, 0, []NeuronIndex{})
	for _, nodeIndex := range allNodes {
		node := net.AllNodes[nodeIndex]
		f(&node, node.distanceFromOutputNode)
	}
}

// ForEachConnectedNodeIndex will only go through nodes that are connected to the output node (directly or indirectly)
// Unconnected input nodes are not covered.
func (net *Network) ForEachConnectedNodeIndex(f func(ni NeuronIndex, distanceFromOutputNode int)) {
	// Start at the output node, traverse left towards the input nodes
	// The network has a counter for how many nodes has been added/removed, for quick memory allocation here
	allNodes := net.getAllNodes(net.OutputNode, 0, []NeuronIndex{})
	for _, nodeIndex := range allNodes {
		node := net.AllNodes[nodeIndex]
		f(nodeIndex, node.distanceFromOutputNode)
	}
}

// String creates a simple and not very useful ASCII representation of the input nodes and the output node.
// Nodes that are not input nodes are skipped.
// Input nodes that are not connected directly to the output node are drawn as non-connected,
// even if they are connected via another node.
func (net Network) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Network (%d nodes, %d input nodes, %d output node)\n", len(net.AllNodes), len(net.InputNodes), 1))
	sb.WriteString("\tInputs to output node: " + strconv.Itoa(len(net.AllNodes[net.OutputNode].InputNeurons)) + "\n")
	for _, node := range net.All() {
		sb.WriteString("\t" + node.String() + "\n")
	}
	return sb.String()
}
