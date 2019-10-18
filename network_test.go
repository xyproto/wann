package wann

import (
	"fmt"
	"testing"
)

func TestNetwork(t *testing.T) {
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println(net)
}

// func NewNetwork(c *Config) *Network {
// func (net *Network) InsertNode(a, b, newNode *Neuron) error {
// func (net *Network) AddConnection(a, b *Neuron) error {
// func (net *Network) ChangeActivationFunction(n *Neuron, f func(float64) float64) {
// func (net *Network) String() string {
