package wann

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var currentTime = time.Now().UTC().UnixNano()

func TestNetwork(t *testing.T) {
	rand.Seed(currentTime)
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

func TestForEachConnected(t *testing.T) {
	rand.Seed(currentTime)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println("<connected_nodes>")
	net.ForEachConnected(func(n *Neuron, _ int) {
		fmt.Println(n)
	})
	fmt.Println("</connected_nodes>")
}

func TestAll(t *testing.T) {
	rand.Seed(currentTime)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println("<all_nodes>")
	for _, node := range net.All() {
		fmt.Println(node)
	}
	fmt.Println("</all_nodes>")
}
