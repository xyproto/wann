package wann

import (
	"errors"
	"fmt"

	"github.com/xyproto/af"
)

// Neuron is a list of input-neurons, and an activation function.
type Neuron struct {
	InputNeurons       []*Neuron
	ActivationFunction func(float64) float64
	Value              *float64
}

// NewNeuron creates a new *Neuron, with an id
func NewNeuron() *Neuron {
	// Pre-allocate room for 64 connections and use Linear as the default activation function
	return &Neuron{InputNeurons: make([]*Neuron, 0, 64), ActivationFunction: af.Linear}
}

// SetValue can be used for setting a value for this neuron instead of using input neutrons.
// This changes how the Evaluation function behaves.
func (neuron *Neuron) SetValue(x float64) {
	neuron.Value = &x
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

// FindInput checks if the given neuron is an input neuron to this one,
// and also returns the index to InputNeurons, if found.
func (neuron *Neuron) FindInput(e *Neuron) (int, bool) {
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
	if i, found := neuron.FindInput(e); found {
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

// Evaluate will return a weighted sum of the input nodes,
// using the .Value field if it is set and no input nodes are available.
func (neuron *Neuron) Evaluate(weight float64) float64 {
	// Assume this is the Output neuron, recursively evaluating the result
	// For each input neuron, evaluate them
	summed := 0.0
	counter := 0
	for _, inputNeuron := range neuron.InputNeurons {
		// Let each input neuron do its own evauluation, using the given weight
		summed += inputNeuron.Evaluate(weight) * weight
		counter++
	}
	// No input neurons. Use the .Value field if it's not nil.
	if counter == 0 && neuron.Value != nil {
		return neuron.ActivationFunction(*(neuron.Value))
	}
	// Return the averaged sum
	return neuron.ActivationFunction(summed / float64(counter))
}
