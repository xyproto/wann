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
	result := net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
	if result != 0.5415839586477849 {
		t.Fail()
	}
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
	result := net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
	if result != 0.5 {
		t.Fail()
	}
}

// 	func (net *Network) AddConnection(a, b NeuronIndex) error {
// 	func (net *Network) ChangeActivationFunction(n *Neuron, f func(float64) float64) {
// 	func (net *Network) String() string {
// 	func (net *Network) SetWeight(weight float64) {
// 	func (net *Network) Complexity() float64 {
// 	func (net *Network) LeftRight(a, b NeuronIndex) (left NeuronIndex, right NeuronIndex) {
// 	func (net *Network) Depth() int {
// 	func (net *Network) checkInputNeurons() {
// 	func (net Network) Copy() Network {
// 	func (net *Network) GetRandomNeuron() NeuronIndex {
// 	func (net *Network) GetRandomInputNode() NeuronIndex {
// 	func (node *Neuron) In(collection []NeuronIndex) bool {
// 	func Combine(a, b []NeuronIndex) []NeuronIndex {
// 	func (net *Network) getAllNodes(nodeIndex NeuronIndex, distanceFromFirstNode int, alreadyHaveThese []NeuronIndex) []NeuronIndex {
// 	func (net *Network) ForEachConnectedNodeIndex(f func(ni NeuronIndex, distanceFromOutputNode int)) {
