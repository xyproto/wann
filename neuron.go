package wann

import (
	"errors"
	"fmt"
	"github.com/xyproto/swish"
)

// Neuron is a list of input-neurons, and an activation function.
type Neuron struct {
	id                 int
	InputNeurons       []*Neuron
	ActivationFunction func(float64) float64
}

// NewNeuron creates a new *Neuron, with an id
func NewNeuron() *Neuron {
	// Pre-allocate room for 1024 connections and use Swish as the default activation function
	return &Neuron{InputNeurons: make([]*Neuron, 0, 1024), ActivationFunction: swish.Swish}
}

// HasInput checks if the given neuron is an input neuron to this one
func (neuron *Neuron) HasInput(e *Neuron) bool {
	for _, n := range neuron.InputNeurons {
		if n == e {
			return true
		}
	}
	return false
}

// FindInputNeuron checks if the given neuron is an input neuron to this one,
// and also returns the index to InputNeurons, if found.
func (neuron *Neuron) FindInputNeuron(e *Neuron) (int, bool) {
	for i, n := range neuron.InputNeurons {
		if n == e {
			return i, true
		}
	}
	return -1, false
}

// AddInput will add an input neuron
func (neuron *Neuron) AddInput(e *Neuron) error {
	if neuron.HasInput(e) {
		return errors.New("neuron already exists")
	}
	neuron.InputNeurons = append(neuron.InputNeurons, e)
	return nil
}

// RemoveInput will remove an input neuron
func (neuron *Neuron) RemoveInput(e *Neuron) error {
	if i, found := neuron.FindInputNeuron(e); found {
		// Found it, remove the neuron at index i
		neuron.InputNeurons = append(neuron.InputNeurons[:i], neuron.InputNeurons[i+1:]...)
		return nil
	}
	return errors.New("neuron does not exist")
}

// String will return a string containing both the pointer address and the number of input neurons
func (neuron *Neuron) String() string {
	return fmt.Sprintf("NEURON[%p,%d]", neuron, len(neuron.InputNeurons))
}
