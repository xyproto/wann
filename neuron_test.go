package wann

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestNeuron(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()
	n, _ := net.NewNeuron()
	n.ActivationFunctionIndex = Swish
	result := n.ActivationFunction()(0.5)
	diff := math.Abs(result - 0.311287)
	if diff > 0.00001 { // 0.0000001 {
		t.Errorf("default swish activation function, expected a number close to 0.311287, got %f:", result)
	}

	fmt.Printf("Neurons in network: %d\n", len(net.AllNodes))
}

func TestString(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()
	n, _ := net.NewNeuron()
	_ = n.String()
}

func TestHasInput(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()     // 0
	a, _ := net.NewNeuron() // 1
	b, _ := net.NewNeuron() // 2
	fmt.Println("a is 1?", a)
	fmt.Println("b is 2?", b)
	a.AddInput(0)
	if !a.HasInput(0) {
		t.Errorf("a should have b as an input")
	}
	if b.HasInput(0) {
		t.Errorf("b should not have a as an input")
	}
}

func TestFindInput(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()

	a, _ := net.NewNeuron()  // a, 1
	_, bi := net.NewNeuron() // b, 2
	c, ci := net.NewNeuron() // c, 3
	_, di := net.NewNeuron() //  d, 4

	a.AddInput(bi)      // b
	a.AddInputNeuron(c) // c

	if _, found := a.FindInput(di); found {
		t.Errorf("a should not have d as an input")
	}
	if pos, found := a.FindInput(bi); !found {
		t.Errorf("a should have b as an input")
	} else if found && pos != 0 {
		t.Errorf("a should have b as an input at position 0")
	}
	if pos, found := a.FindInput(ci); !found {
		t.Errorf("a should have c as an input")
	} else if found && pos != 1 {
		t.Errorf("a should have c as an input at position 1")
	}
}

func TestRemoveInput(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})

	a, _ := net.NewNeuron() // 0
	a.AddInput(1)
	a.AddInput(2)
	if a.RemoveInput(1) != nil {
		t.Errorf("could not remove input b from a")
	}
	if a.RemoveInput(2) != nil {
		t.Errorf("could not remove input c from a")
	}
	if a.HasInput(1) {
		t.Errorf("a should not have b as an input")
	}
	if a.HasInput(2) {
		t.Errorf("a should not have c as an input")
	}
}

// func (neuron *Neuron) RemoveInput(e *Neuron) error {

func TestEvaluate(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          7,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})

	// Set a few activation functions
	net.AllNodes[net.InputNodes[0]].ActivationFunctionIndex = Linear
	net.AllNodes[net.InputNodes[1]].ActivationFunctionIndex = Swish
	net.AllNodes[net.InputNodes[2]].ActivationFunctionIndex = Gauss
	net.AllNodes[net.InputNodes[3]].ActivationFunctionIndex = Sigmoid
	net.AllNodes[net.InputNodes[4]].ActivationFunctionIndex = ReLU
	net.AllNodes[net.InputNodes[5]].ActivationFunctionIndex = Step
	net.AllNodes[net.InputNodes[6]].ActivationFunctionIndex = Inv

	result := net.Evaluate([]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5})
	fmt.Println(result)
}
