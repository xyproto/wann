package wann

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/xyproto/af"
	"github.com/xyproto/swish"
)

func TestNeuron(t *testing.T) {
	n := NewNeuron()
	n.ActivationFunction = swish.Swish
	result := n.ActivationFunction(0.5)
	diff := math.Abs(result - 0.311287)
	if diff > 0.00001 { // 0.0000001 {
		t.Errorf("default swish activation function, expected a number close to 0.311287, got %f:", result)
	}
}

func TestString(t *testing.T) {
	n := NewNeuron()
	s := n.String()
	if !strings.HasPrefix(s, "NEURON[") || !strings.HasSuffix(s, "]") {
		t.Errorf("could not convert neuron to a string")
	}
}

func TestAddInput(t *testing.T) {
	a := NewNeuron()
	b := NewNeuron()
	a.AddInput(b)
	if !strings.HasSuffix(a.String(), ",1]") {
		t.Errorf("a should have one connection")
	}
	if !strings.HasSuffix(b.String(), ",0]") {
		t.Errorf("b should have zero connections")
	}
}

func TestHasInput(t *testing.T) {
	a := NewNeuron()
	b := NewNeuron()
	a.AddInput(b)
	if !a.HasInput(b) {
		t.Errorf("a should have b as an input")
	}
	if b.HasInput(a) {
		t.Errorf("b should not have a as an input")
	}
}

func TestFindInput(t *testing.T) {
	a := NewNeuron()
	b := NewNeuron()
	c := NewNeuron()
	d := NewNeuron()
	a.AddInput(b)
	a.AddInput(c)
	if _, found := a.FindInput(d); found {
		t.Errorf("a should have d as an input")
	}
	if pos, found := a.FindInput(b); !found {
		t.Errorf("a should have b as an input")
	} else if found && pos != 0 {
		t.Errorf("a should have b as an input at position 0")
	}
	if pos, found := a.FindInput(c); !found {
		t.Errorf("a should have c as an input")
	} else if found && pos != 1 {
		t.Errorf("a should have c as an input at position 1")
	}
}

func TestRemoveInput(t *testing.T) {
	a := NewNeuron()
	b := NewNeuron()
	c := NewNeuron()
	a.AddInput(b)
	a.AddInput(c)
	if a.RemoveInput(b) != nil {
		t.Errorf("could not remove input b from a")
	}
	if a.RemoveInput(c) != nil {
		t.Errorf("could not remove input c from a")
	}
	if a.HasInput(b) {
		t.Errorf("a should not have b as an input")
	}
	if a.HasInput(c) {
		t.Errorf("a should not have c as an input")
	}
}

// func (neuron *Neuron) RemoveInput(e *Neuron) error {

func TestEvaluate(t *testing.T) {
	rand.Seed(currentTime)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})

	// Set a few activation functions
	net.Nodes[0].ActivationFunction = af.Linear
	net.Nodes[0].SetValue(0.5)
	net.Nodes[1].ActivationFunction = af.Swish
	net.Nodes[1].SetValue(0.5)
	net.Nodes[2].ActivationFunction = af.Gaussian01
	net.Nodes[2].SetValue(0.5)
	net.Nodes[3].ActivationFunction = af.Sigmoid
	net.Nodes[3].SetValue(0.5)
	net.Nodes[4].ActivationFunction = af.ReLU
	net.Nodes[4].SetValue(0.5)

	fmt.Println(net.OutputNode.Evaluate(0.5))
}
