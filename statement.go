package wann

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
)

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
			//for _, inputStatement := range inputStatements {
			//	fmt.Printf("\t%#v\n", inputStatement)
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

// Render renders a *jen.Statement to a string, if possible
// if there is an error about an extra ")", then that's because anonymous functions are not supported by jen
// Do not Render until statements could be placed at the top-level in a Go program.
func Render(inner *jen.Statement) string {
	return inner.GoString()
}

// OutputNodeStatementX returns a statement for the output node, using "x" for the variable
func (net *Network) OutputNodeStatementX(functionName string) string {
	inner := net.AllNodes[net.OutputNode].ActivationFunction.Statement(jen.Id("x"))
	f := jen.Id(functionName).Op(":=").Add(inner)
	return Render(f)
}

// RunStatementX will run the given statement by wrapping it in a program and using "go run"
func RunStatementX(statement *jen.Statement, x float64) (float64, error) {
	file, err := ioutil.TempFile("", "af_*.go")
	if err != nil {
		return 0.0, err
	}
	filename := file.Name()
	defer os.Remove(filename)
	// Build the contents of the source file using jennifer
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Id("x").Op(":=").Lit(x),
		jen.Qual("fmt", "Println").Call(statement),
	)
	// Save the file
	if ioutil.WriteFile(filename, []byte(f.GoString()), 0664) != nil {
		return 0.0, err
	}
	// Run the file
	cmd := exec.Command("go", "run", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0.0, err
	}
	// Return the outputted float string as a float64
	resultString := strings.TrimSpace(string(out))
	resultFloat, err := strconv.ParseFloat(resultString, 64)
	if err != nil {
		return 0.0, err
	}
	return resultFloat, nil
}

// RunStatementInputData will run the given statement by wrapping it in a program and using "go run"
func RunStatementInputData(statement *jen.Statement, inputData []float64) (float64, error) {
	file, err := ioutil.TempFile("", "af_*.go")
	if err != nil {
		return 0.0, err
	}
	filename := file.Name()
	defer os.Remove(filename)
	// Build the contents of the source file using jennifer
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		// Build a statement that declares and initializes "inputData" based on the contents of inputData
		jen.Id("inputData").Op(":=").Index().Float64().ValuesFunc(func(g *jen.Group) {
			for i := 0; i < len(inputData); i++ {
				g.Lit(inputData[i])
			}
		}),
		jen.Qual("fmt", "Println").Call(statement),
	)
	// Save the file
	if ioutil.WriteFile(filename, []byte(f.GoString()), 0664) != nil {
		return 0.0, err
	}
	// Run the file
	cmd := exec.Command("go", "run", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0.0, err
	}
	// Return the outputted float string as a float64
	resultString := strings.TrimSpace(string(out))
	resultFloat, err := strconv.ParseFloat(resultString, 64)
	if err != nil {
		return 0.0, err
	}
	return resultFloat, nil
}
