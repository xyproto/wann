package wann

import (
	"fmt"
	"math/rand"
	"testing"
)

// Use a pseudo-random seed
//var commonSeed = time.Now().UTC().UnixNano()

// Use a specific seed
var commonSeed int64 = 1571917826405889420

func TestNewNetwork(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println(net)
}

func TestGet(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()
	fmt.Println(net)
	fmt.Println(net.Get(0))
	if net.OutputNode != 0 {
		t.Fail()
	}
}

func TestIsInput(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	if !net.IsInput(1) {
		t.Fail()
	}
	if net.IsInput(0) {
		t.Fail()
	}
}

func TestForEachConnected(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	net.ForEachConnected(func(n *Neuron) {
		fmt.Printf("%d: %s, distance from output node: %d\n", n.neuronIndex, n, n.distanceFromOutputNode)
	})
}

func TestAll(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.7,
		SharedWeight:    0.5,
	})
	for _, node := range net.All() {
		fmt.Println(node)
	}
}

func TestEvaluate2(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.7,
		SharedWeight:    0.5,
	})
	_ = net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
}

func TestInsertNode(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	_, newNeuronIndex := net.NewRandomNeuron()
	if err := net.InsertNode(0, 2, newNeuronIndex); err != nil {
		t.Error(err)
	}
	//fmt.Println(net)
	_ = net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
}

func TestAddConnection(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	_, newNeuronIndex := net.NewRandomNeuron()
	if err := net.InsertNode(net.OutputNode, 2, newNeuronIndex); err != nil {
		t.Error(err)
	}
	// Add a connection from 1 to the new neuron.
	// This is the same as making the new neuron have an additional input neuron: index 1
	if err := net.AddConnection(1, newNeuronIndex); err != nil {
		t.Error(err)
	}
	// Add a connection from the output node to the output node. Should fail.
	if err := net.AddConnection(net.OutputNode, net.OutputNode); err == nil {
		t.Fail()
	}
	// Adding a made-up index should fail as well
	if err := net.AddConnection(net.OutputNode, 999); err == nil {
		t.Fail()
	}
}

func TestRandomizeActivationFunctionForRandomNeuron(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	net.RandomizeActivationFunctionForRandomNeuron()
}

func TestNetworkString(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	//fmt.Println(net.String())
	_ = net.String()
}

func TestSetWeight(t *testing.T) {
	net := NewNetwork()
	net.SetWeight(0.1234)
	if net.Weight != 0.1234 {
		t.Fail()
	}
}

func TestComplexity(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.0,
	})
	// The complexity will vary, because the performance varies when
	// estimating the complexity of each function.
	// But the complexity compared between networks should still hold.
	firstComplexity := net.Complexity()
	//fmt.Println("First network complexity:", firstComplexity)
	// Adding a connection increases the complexity
	net.AddConnection(0, 1)
	secondComplexity := net.Complexity()
	//fmt.Println("Second network complexity:", secondComplexity)
	if firstComplexity >= secondComplexity {
		t.Fail()
	}
}

func TestLeftRight(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          3,
		ConnectionRatio: 1.0,
	})
	net.AllNodes[1].ActivationFunctionIndex = Swish
	a, b := net.LeftRight(0, 1)
	// output node to the right
	if a != 1 || b != 0 {
		t.Fail()
	}
	// output node to the right
	a, b = net.LeftRight(1, 0)
	if a != 1 || b != 0 {
		t.Fail()
	}
	net.WriteSVG("before.svg")
	_, nodeIndex := net.NewRandomNeuron()
	err := net.InsertNode(0, 1, nodeIndex)
	if err != nil {
		t.Error(err)
	}
	net.WriteSVG("after.svg")
	fmt.Println("A")
	fmt.Println("NEW NODE INDEX IS", nodeIndex)
	a, b = net.LeftRight(0, nodeIndex)
	// output node to the right
	if a != nodeIndex || b != 0 {
		t.Fail()
	}
	fmt.Println("B")
	a, b = net.LeftRight(nodeIndex, 0)
	// output node to the right
	if a != nodeIndex || b != 0 {
		t.Fail()
	}
	fmt.Println("C")
	a, b = net.LeftRight(1, nodeIndex)
	// Here, the new node should be to the right, since it's between node 1 and the output node
	if a != 1 || b != nodeIndex {
		t.Fail()
	}
	fmt.Println("D")
	//net.WriteSVG("c.svg")
	fmt.Println(net)
	a, b = net.LeftRight(nodeIndex, 1)
	fmt.Println("nodeIndex:", nodeIndex)
	fmt.Println("1:", 1)
	fmt.Println("a:", a)
	fmt.Println("b:", b)

	if a != 1 || b != nodeIndex {
		t.Fail()
	}
}

func TestDepth(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          3,
		ConnectionRatio: 1.0,
	})
	fmt.Println(net.Depth())
	_, nodeIndex := net.NewNeuron()
	_ = net.InsertNode(0, 1, nodeIndex)
	fmt.Println(net.Depth())
}

// 	func (net *Network) checkInputNeurons() {
// 	func (net Network) Copy() Network {
// 	func (net *Network) GetRandomNeuron() NeuronIndex {
// 	func (net *Network) GetRandomInputNode() NeuronIndex {
// 	func (node *Neuron) In(collection []NeuronIndex) bool {
// 	func Combine(a, b []NeuronIndex) []NeuronIndex {
// 	func (net *Network) getAllNodes(nodeIndex NeuronIndex, distanceFromFirstNode int, alreadyHaveThese []NeuronIndex) []NeuronIndex {
// 	func (net *Network) ForEachConnectedNodeIndex(f func(ni NeuronIndex, distanceFromOutputNode int)) {
