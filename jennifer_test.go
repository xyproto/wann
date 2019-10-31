package wann

import (
	"fmt"
	"math/rand"
)

func ExampleNetwork_JenniferOutputNodeX_first() {
	// First create a network with only one output node, that has a step function
	net := NewNetwork()
	net.AllNodes[net.OutputNode].ActivationFunction = Step

	fmt.Println(net.JenniferOutputNodeX("score"))

	// Output:
	// score := func(s float64) float64 {
	// 	if s >= 0 {
	// 		return 1
	// 	} else {
	// 		return 0
	// 	}
	// }(x)
}

func ExampleNetwork_JenniferOutputNodeX_second() {
	// Then create a network with an input node that has a sigmoid function and an output node that has an invert function
	net := NewNetwork()
	net.NewInputNode(Sigmoid, true)
	net.AllNodes[net.OutputNode].ActivationFunction = Inv

	// Output a Go expression for this network, using the given input variable names
	fmt.Println(net.JenniferOutputNodeX("score"))

	// Output:
	// score := -(x)

}

func ExampleNetwork_JenniferOutputNodeX_third() {
	rand.Seed(999)
	net := NewNetwork(&Config{
		inputs:                 1,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.JenniferOutputNodeX("score"))

	// Output:
	// score := math.Exp(-(math.Pow(x, 2.0)) / 2.0)
}
func ExampleNetwork_JenniferOutputNodeX_fourth() {

	rand.Seed(1111113)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.JenniferOutputNodeX("score"))

	// Output:
	// score := func(r float64) float64 {
	//	if r >= 0 {
	//		return r
	//	} else {
	//		return 0
	//	}
	// }(x)
}

func ExampleNetwork_JenniferOutputNodeX() {
	rand.Seed(1)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	fmt.Println(net.JenniferOutputNodeX("f"))
	//fmt.Println(net.Score())

	// Output:
	// f := math.Pow(x, 2.0)
}
