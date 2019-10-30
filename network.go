package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

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
	n := c.inputs
	r := c.InitialConnectionRatio
	w := c.sharedWeight
	// Create a new network that has one node, the output node
	outputNodeIndex := NeuronIndex(0)
	net := Network{make([]Neuron, 0, n+1), make([]NeuronIndex, n), outputNodeIndex, w}
	outputNode, outputNodeIndex := net.NewNeuron()
	net.OutputNode = outputNodeIndex

	// Initialize n input nodes that all are inputs to the one output node.
	for i := 0; i < n; i++ {
		// Add a new input node
		_, nodeIndex := net.NewNeuron()

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

	// Sort the nodes by where they place in the diagram
	a, b, arbitrary := net.LeftRight(a, b)
	if arbitrary {
		if a == b {
			return errors.New("insert node: the a and b nodes are the same")
		}
		if net.IsInput(a) && net.IsInput(b) {
			return errors.New("insert node: both node a and b are input nodes")
		}
		return errors.New("insert node: aribtrary ordering when inserting a node")
	}

	// This should never happen
	if a == net.OutputNode {
		panic("implementation error: the leftmost node is an output node and this was not cought earlier")
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

	// Connect the new node to b
	if err := net.AllNodes[b].AddInput(newNodeIndex); err != nil {
		// This does not kick in, the problem must be in AddInput!
		return errors.New("error in InsertNode b.AddInput(newNode): " + err.Error())
	}

	// Connect a to the new node
	if err := net.AllNodes[newNodeIndex].AddInput(a); err != nil {
		return errors.New("error in InsertNode newNode.AddInput(a): " + err.Error())
	}

	// The situation should now be: a -> newNode -> b
	return nil
}

// AddConnection adds a connection from a to b
func (net *Network) AddConnection(a, b NeuronIndex) error {
	lastIndex := NeuronIndex(len(net.AllNodes) - 1)
	if a < 0 || a > lastIndex || b < 0 || b > lastIndex {
		return errors.New("index out of range")
	}
	// Sort the nodes by where they place in the diagram
	arbitrary := false
	a, b, arbitrary = net.LeftRight(a, b)
	if arbitrary {
		if a == b {
			return errors.New("can't connect to self")
		}
		return errors.New("error: arbitrary ordering when adding a connection")
	}
	// a should not be an output node
	if a == net.OutputNode {
		return errors.New("error: will not insert a node between the output node and another node")
	}
	// b should not be an input node
	if net.IsInput(b) {
		return errors.New("error: b is an input node")
	}
	// same thing
	if net.AllNodes[b].Value != nil {
		return errors.New("error: b is an input node")
	}
	if net.AllNodes[b].HasInput(a) {
		return errors.New("error: input already exists")
	}
	return net.AllNodes[b].AddInput(a)
}

// RandomizeActivationFunctionForRandomNeuron randomizes the activation function for a randomly selected neuron
func (net *Network) RandomizeActivationFunctionForRandomNeuron() {
	chosenNeuronIndex := net.GetRandomNode()
	chosenActivationFunctionIndex := rand.Intn(len(ActivationFunctions))
	net.AllNodes[chosenNeuronIndex].ActivationFunctionIndex = chosenActivationFunctionIndex
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
	outputNode := net.AllNodes[net.OutputNode]
	maxIterationCounter := inputLength
	result, _ := outputNode.evaluate(net.Weight, &maxIterationCounter)
	return result
}

// SetWeight will set a shared weight for the entire network
func (net *Network) SetWeight(weight float64) {
	net.Weight = weight
}

// Complexity measures the network complexity
// Will return 1.0 at a minimum
func (net *Network) Complexity() float64 {

	// TODO: These two constants really affect the results. Place them in the Config struct instead.

	// How much should the function complexity matter in relation to the number of connected nodes?
	const functionComplexityMultiplier = 10.0

	// How much should the complexity score matter in relation to the network results, when scoring the network?
	const complexityMultiplier = 10.0

	sum := 0.0
	// Sum the complexity of all activation functions.
	// This penalizes both slow activation functions and
	// unconnected nodes.
	for _, n := range net.AllNodes {
		if n.Value == nil {
			sum += ComplexityEstimate[n.ActivationFunctionIndex] * functionComplexityMultiplier
		}
	}
	// The number of connected nodes should also carry some weight
	connectedNodes := float64(len(net.Connected()))
	// This must always be larger than 0, to avoid divide by zero later
	return (connectedNodes + sum) * complexityMultiplier
}

// LeftRight returns two neurons, such that the first on is the one that is
// most to the left (towards the input neurons) and the second one is most to
// the right (towards the output neuron). Assumes that a and b are not equal.
// The returned bool is true if there is no order (if the nodes are equal, both are output nodes or both are input nodes)
func (net *Network) LeftRight(a, b NeuronIndex) (NeuronIndex, NeuronIndex, bool) {
	// First check if they are equal
	if a == b {
		return a, b, true // Arbitrary order
	}
	// First check the network output nodes
	if a == net.OutputNode && b == net.OutputNode {
		return a, b, true // Arbitrary order
	}
	if a == net.OutputNode && b != net.OutputNode {
		return b, a, false // Swap order
	}
	if a != net.OutputNode && b == net.OutputNode {
		return a, b, false // Same order
	}
	// Then check if the nodes are already connected
	if net.AllNodes[a].In(net.AllNodes[b].InputNodes) {
		return a, b, false // Same order
	}
	if net.AllNodes[b].In(net.AllNodes[a].InputNodes) {
		return b, a, false // Swap order
	}
	// Then check the input nodes of the network
	aIsNetworkInputNode := net.AllNodes[a].In(net.InputNodes)
	bIsNetworkInputNode := net.AllNodes[b].In(net.InputNodes)
	if aIsNetworkInputNode && !bIsNetworkInputNode {
		return a, b, false // Same order
	}
	if !aIsNetworkInputNode && bIsNetworkInputNode {
		return b, a, false // Swap order
	}
	if aIsNetworkInputNode && bIsNetworkInputNode {
		return a, b, true // Arbitrary order
	}
	// Then check the distance from the output node, in steps
	aDistance := net.AllNodes[a].distanceFromOutputNode
	bDistance := net.AllNodes[b].distanceFromOutputNode
	if bDistance > aDistance {
		return b, a, false // Swap order, b is further away from the output node, which (usually) means further left in the graph
	}
	// Everything else
	return a, b, false
}

// Depth returns the maximum connection distance from the output node
func (net *Network) Depth() int {
	maxDepth := 0
	net.ForEachConnected(func(n *Neuron) {
		if n.distanceFromOutputNode > maxDepth {
			maxDepth = n.distanceFromOutputNode
		}
	})
	return maxDepth
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

// GetRandomNode will select a random neuron.
// This can be any node, including the output node.
func (net *Network) GetRandomNode() NeuronIndex {
	return NeuronIndex(rand.Intn(len(net.AllNodes)))
}

// GetRandomInputNode returns a random input node
func (net *Network) GetRandomInputNode() NeuronIndex {
	inputPosition := rand.Intn(len(net.InputNodes))
	inputNodeIndex := net.InputNodes[inputPosition]
	return inputNodeIndex
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
func (net *Network) getAllConnectedNodes(nodeIndex NeuronIndex, distanceFromFirstNode int, alreadyHaveThese []NeuronIndex) []NeuronIndex {
	allNodes := make([]NeuronIndex, 0, len(net.AllNodes))
	node := net.AllNodes[nodeIndex]
	if nodeIndex != net.OutputNode {
		node.distanceFromOutputNode = distanceFromFirstNode
		net.AllNodes[nodeIndex] = node
	}
	if !node.In(alreadyHaveThese) {
		allNodes = append(allNodes, nodeIndex)
	}
	for _, inputNodeIndex := range node.InputNodes {
		if node.Is(inputNodeIndex) {
			panic("implementation error: node is input node to self")
		}
		if int(inputNodeIndex) >= len(net.AllNodes) {
			continue
		}
		inputNode := net.AllNodes[inputNodeIndex]
		if !inputNode.In(allNodes) && !inputNode.In(alreadyHaveThese) {
			allNodes = Combine(allNodes, net.getAllConnectedNodes(inputNodeIndex, distanceFromFirstNode+1, append(allNodes, alreadyHaveThese...)))
		}
	}
	return allNodes
}

// ForEachConnected will only go through nodes that are connected to the output node (directly or indirectly)
// Unconnected input nodes are not covered.
func (net *Network) ForEachConnected(f func(n *Neuron)) {
	// Start at the output node, traverse left towards the input nodes
	// The network has a counter for how many nodes has been added/removed, for quick memory allocation here
	// the final slice is to avoid circular connections
	for _, nodeIndex := range net.getAllConnectedNodes(net.OutputNode, 0, []NeuronIndex{}) {
		f(&(net.AllNodes[nodeIndex]))
	}
}

// Connected returns a slice of neuron indexes, that are all connected to the output node (directly or indirectly)
func (net *Network) Connected() []NeuronIndex {
	allConnected := make([]NeuronIndex, 0, len(net.AllNodes)) // Use a bit more memory, but don't allocate at every iteration
	net.ForEachConnectedNodeIndex(func(ni NeuronIndex) {
		allConnected = append(allConnected, ni)
	})
	return allConnected
}

// Unconnected returns a slice of all unconnected neurons
func (net *Network) Unconnected() []NeuronIndex {
	connected := net.Connected()
	// TODO: Benchmark if using len(net.AllNodes) here is faster or not
	unconnected := make([]NeuronIndex, 0, len(net.AllNodes))
	for i, node := range net.AllNodes {
		if !node.In(connected) {
			unconnected = append(unconnected, NeuronIndex(i))
		}
	}
	return unconnected
}

// ForEachConnectedNodeIndex will only go through nodes that are connected to the output node (directly or indirectly)
// Unconnected input nodes are not covered.
func (net *Network) ForEachConnectedNodeIndex(f func(ni NeuronIndex)) {
	net.ForEachConnected(func(n *Neuron) {
		f(n.neuronIndex)
	})
}

// InsertRandomNode will try the given number of times to insert a node in a random location,
// replacing an existing connection between two nodes.
// `a -> b` will then become `a -> newNode -> b`
// Returns true if one was inserted or false if the randomly chosen location wasn't fruitful
func (net *Network) InsertRandomNode() bool {

	// Find a random node among the nodes that are connected to the output node (directly or indirectly)
	connectedNodes := net.Connected()
	randomNodeIndexThatIsConnected := connectedNodes[rand.Intn(len(connectedNodes))]

	// If this is one of the network input nodes, return
	if net.IsInput(randomNodeIndexThatIsConnected) {
		// Nothing to do here, the input nodes get their input from the input numbers
		return false
	}

	// If this is the output node, and there are no inputs to the output node, return
	if randomNodeIndexThatIsConnected == net.OutputNode && len(net.AllNodes[net.OutputNode].InputNodes) == 0 {
		// Nothing to do here, no connections to the output node
		return false
	}

	// If we arrived here, this node must have input nodes. Choose one at random.
	rightIndex := randomNodeIndexThatIsConnected
	inputNodes := net.AllNodes[rightIndex].InputNodes

	leftIndex := inputNodes[rand.Intn(len(inputNodes))]

	// We now have a left and right node index, that we know are connected, replace this connection with
	// one that goes through an entirely new node.

	// Create a new node and connect it with the left node
	newNode, newNodeIndex := net.NewNeuron()
	err := newNode.AddInput(leftIndex)
	if err != nil {
		panic(err)
	}

	// Remove the connection to the left node and then connect the right node to the new node
	err = net.AllNodes[rightIndex].RemoveInput(leftIndex)
	if err != nil {
		panic(err)
	}
	err = net.AllNodes[rightIndex].AddInput(newNodeIndex)
	if err != nil {
		panic(err)
	}

	return true
}

// UpdateNetworkPointers will update all the node.Net pointers to point to this network
func (net *Network) UpdateNetworkPointers() {
	for nodeIndex := range net.AllNodes {
		net.AllNodes[nodeIndex].Net = net
	}
}

// Copy a Network to a new network
func (net Network) Copy() *Network {
	var newNet Network
	newNet.AllNodes = make([]Neuron, len(net.AllNodes))
	for nodeIndex := range net.AllNodes {
		// This copies the node and also sets the .Net pointer correctly to this network
		newNet.AllNodes[nodeIndex] = newNet.AllNodes[nodeIndex].Copy(&newNet)
	}
	newNet.InputNodes = net.InputNodes
	newNet.OutputNode = net.OutputNode
	newNet.Weight = net.Weight

	// NOTE: It's important that a pointer to a Network is returned,
	//       instead of an entire Network struct, so that the .Net pointers in the nodes point correctly.
	return &newNet
}

// String creates a simple and not very useful ASCII representation of the input nodes and the output node.
// Nodes that are not input nodes are skipped.
// Input nodes that are not connected directly to the output node are drawn as non-connected,
// even if they are connected via another node.
func (net Network) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Network (%d nodes, %d input nodes, %d output node)\n", len(net.AllNodes), len(net.InputNodes), 1))
	sb.WriteString("\tConnected inputs to output node: " + strconv.Itoa(len(net.AllNodes[net.OutputNode].InputNodes)) + "\n")
	for _, node := range net.All() {
		sb.WriteString("\t" + node.String() + "\n")
	}
	return sb.String()
}
