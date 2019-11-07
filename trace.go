package wann

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
)

// In returns true if this NeuronIndex is in the given *[]NeuronIndex slice
func (ni NeuronIndex) In(nodes *[]NeuronIndex) bool {
	for _, ni2 := range *nodes {
		if ni2 == ni {
			return true
		}
	}
	return false
}

var errIgnore = errors.New("ignore this node")

// SetInputValues will assign the given values to the network input nodes
func (net Network) SetInputValues(inputValues []float64) {
	if len(net.InputNodes) > len(inputValues) {
		fmt.Println("warning: more input nodes than input values")
	} else if len(net.InputNodes) < len(inputValues) {
		fmt.Println("warning: fewer input nodes than input values")
	}
	// Assign the values
	for i, ni := range net.InputNodes {
		if i < len(inputValues) {
			v := inputValues[i]
			net.AllNodes[ni].Value = &v
		} else {
			break
		}
	}
}

// NetworkStatementWithInputValues will print out a trace of visiting all nodes from output and to the left
func (neuron Neuron) NetworkStatementWithInputValues(visited *[]NeuronIndex) (*jen.Statement, error) {
	// First guard against re-visits
	if neuron.neuronIndex.In(visited) {
		return jen.Empty(), errors.New("already visited: " + strconv.Itoa(int(neuron.neuronIndex)))
	}
	*visited = append(*visited, neuron.neuronIndex)

	// Toggle bits in neuronType to signify what type of neuron it is
	neuronType := 0
	if neuron.IsOutput() {
		neuronType ^= 1
	}
	if neuron.IsInput() {
		neuronType ^= 2
	}

	// Switch on the neuronType bits
	switch neuronType {
	case 0: // not network output and not network input, may have input nodes
		//fmt.Println("* Middle node")
		switch len(neuron.InputNodes) {
		case 0:
			//fmt.Println("No input nodes to this node, and not a network input node.")
			return jen.Empty(), errIgnore
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputValues(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputValues(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 1: // network output and not network input, may have input nodes
		//fmt.Println("* Network output node and not network input node")
		switch len(neuron.InputNodes) {
		case 0:
			//fmt.Println("No input nodes to this node, and not a network input node.")
			return jen.Empty(), errIgnore
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputValues(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputValues(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 2: // not network output, but network input, may not have input nodes
		//fmt.Println("* Network input node and not network output node")
		switch len(neuron.InputNodes) {
		case 0:
			if neuron.Value == nil {
				panic("implementation error: network input Value is nil")
			}
			// No inputs to this network input node, return the Value
			//fmt.Println("literal", *neuron.Value, "activation function", neuron.ActivationFunction.Name())
			inner := jen.Lit(*neuron.Value)
			return neuron.ActivationFunction.Statement(inner), nil
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputValues(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputValues(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 3: // network output and network input, may have input nodes
		//fmt.Println("* Network input and output node")
		switch len(neuron.InputNodes) {
		case 0:
			if neuron.Value == nil {
				panic("implementation error: network input Value is nil")
			}
			// No inputs to this network input node, return the Value
			//fmt.Println("literal", *neuron.Value)
			return jen.Lit(*neuron.Value), nil
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputValues(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputValues(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	}
	panic("implementation error")
}

// StatementWithInputValues traces the entire network
func (net *Network) StatementWithInputValues() (*jen.Statement, error) {
	//fmt.Println("=== Trace ===")
	visited := make([]NeuronIndex, 0)
	outputNode := net.AllNodes[net.OutputNode]
	statement, err := outputNode.NetworkStatementWithInputValues(&visited)
	if err != nil {
		return jen.Empty(), err
	}
	return statement, nil
}

// NetworkStatementWithInputDataVariables will print out a trace of visiting all nodes from output and to the left,
// but with the given slice of statements instead of using the input values
func (neuron Neuron) NetworkStatementWithInputDataVariables(visited *[]NeuronIndex) (*jen.Statement, error) {
	// First guard against re-visits
	if neuron.neuronIndex.In(visited) {
		return jen.Empty(), errors.New("already visited: " + strconv.Itoa(int(neuron.neuronIndex)))
	}
	*visited = append(*visited, neuron.neuronIndex)

	// Toggle bits in neuronType to signify what type of neuron it is
	neuronType := 0
	if neuron.IsOutput() {
		neuronType ^= 1
	}
	if neuron.IsInput() {
		neuronType ^= 2
	}

	// Switch on the neuronType bits
	switch neuronType {
	case 0: // not network output and not network input, may have input nodes
		//fmt.Println("* Middle node")
		switch len(neuron.InputNodes) {
		case 0:
			//fmt.Println("No input nodes to this node, and not a network input node.")
			return jen.Empty(), errIgnore
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputDataVariables(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputDataVariables(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			for _, inputStatement := range inputStatements {
				fmt.Printf("\t%s\n", inputStatement)
			}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 1: // network output and not network input, may have input nodes
		//fmt.Println("* Network output node and not network input node")
		switch len(neuron.InputNodes) {
		case 0:
			//fmt.Println("No input nodes to this node, and not a network input node.")
			return jen.Empty(), errIgnore
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputDataVariables(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputDataVariables(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 2: // not network output, but network input, may not have input nodes
		//fmt.Println("* Network input node and not network output node")
		switch len(neuron.InputNodes) {
		case 0:
			inputStatement, err := neuron.InputStatement()
			if err != nil {
				panic("implementation error: " + err.Error())
			}
			// No inputs to this network input node, return the Value
			//fmt.Println("activation function", neuron.ActivationFunction.Name())
			inner := inputStatement
			return neuron.ActivationFunction.Statement(inner), nil
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputDataVariables(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputDataVariables(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	case 3: // network output and network input, may have input nodes
		//fmt.Println("* Network input and output node")
		switch len(neuron.InputNodes) {
		case 0:
			inputStatement, err := neuron.InputStatement()
			if err != nil {
				panic("implementation error: " + err.Error())
			}
			// No inputs to this network input node, return the statement
			inner := inputStatement
			return inner, nil
		case 1:
			//fmt.Println("One input node to this node.")
			inputNode := neuron.Net.AllNodes[neuron.InputNodes[0]]
			statement, err := inputNode.NetworkStatementWithInputDataVariables(visited)
			if err != nil {
				return jen.Empty(), err
			}
			//fmt.Println("** Statement: ", render(statement))
			return statement, nil
		default:
			var inputStatements []*jen.Statement
			for _, inputNodeIndex := range neuron.InputNodes {
				statement, err := neuron.Net.AllNodes[inputNodeIndex].NetworkStatementWithInputDataVariables(visited)
				if err != nil {
					continue
				}
				inputStatements = append(inputStatements, statement)
			}
			activationStatement := ActivationStatement(neuron.ActivationFunction, neuron.Net.Weight, inputStatements)
			//fmt.Printf("** Statements to combine with %s:\n", neuron.ActivationFunction.Name())
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%s\n", inputStatement)
			//}
			//fmt.Println("** Activation statement: ", render(activationStatement))
			return activationStatement, nil
		}
	}
	panic("implementation error")
}

// StatementWithInputDataVariables traces the entire network, using statements for the input numbers
func (net *Network) StatementWithInputDataVariables() (*jen.Statement, error) {
	//fmt.Println("=== Trace2 ===")
	visited := make([]NeuronIndex, 0)
	outputNode := net.AllNodes[net.OutputNode]
	statement, err := outputNode.NetworkStatementWithInputDataVariables(&visited)
	if err != nil {
		return jen.Empty(), err
	}
	return statement, nil
}
