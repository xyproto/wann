package wann

import (
	"fmt"
	"strconv"
	"strings"
)

// CombineInputs will create an expression that takes an activation function and a bunch of inputs then uses the averaged sum to combine them
func (net *Network) CombineInputs(activationFunction ActivationFunctionIndex, inputNodes []NeuronIndex, inputVariableNames []string) string {
	var inputExpressions []string
	for i, ni := range inputNodes {
		variableName := "?"
		if i < len(inputVariableNames) {
			variableName = inputVariableNames[i]
		}
		multiplyWeight := ""
		if net.Weight != 1.0 {
			multiplyWeight = " * " + strconv.FormatFloat(net.Weight, 'f', -1, 64)
		}
		inner := variableName
		inputExpressions = append(inputExpressions, net.AllNodes[ni].ActivationFunction.goExpression(inner)+multiplyWeight)
	}
	inner := ""
	if len(inputNodes) == 0 {
		variableName := "?"
		if len(inputVariableNames) >= 1 {
			variableName = inputVariableNames[0]
		}
		// multiplyWeight := ""
		// if net.Weight != 1.0 {
		// 	multiplyWeight = " * " + strconv.FormatFloat(net.Weight, 'f', -1, 64)
		// }
		// inner := variableName + multiplyWeight
		inner = variableName
		//return activationFunction.goExpression(variableName)
	} else if len(inputNodes) == 1 {
		inner = inputExpressions[0]
		//return activationFunction.goExpression(inner)
	} else {
		inner = fmt.Sprintf("((%s) / %d.0)", strings.Join(inputExpressions, " + "), len(inputNodes))
	}
	return activationFunction.goExpression(inner)
}

// WrapInFunction wraps a Go expression in a function signature, then places the expression in the body
func WrapInFunction(expression string, variableNames ...string) string {
	usedNames := []string{}
	for _, vn := range variableNames {
		if strings.Contains(expression, vn) {
			usedNames = append(usedNames, vn)
		}
	}
	return fmt.Sprintf("func f(%s float64) float64 { return %s }", strings.Join(usedNames, ", "), expression)
}

// GoExpression returns the source code for a Go expression that does the same thing as this network
func (net *Network) GoExpression(variableNames ...string) string {
	// Make sure the slice of input variable names contains at least one string
	if len(variableNames) == 0 {
		variableNames = []string{"x"}
	}
	// A special case, no input nodes
	if len(net.InputNodes) == 0 {
		return net.AllNodes[net.OutputNode].ActivationFunction.goExpression(variableNames[0])
	}
	// A special case, one input node that is directly connected to the output node, and that's it
	outputNodeInputNodes := net.AllNodes[net.OutputNode].InputNodes
	if len(net.InputNodes) == 1 && len(outputNodeInputNodes) == 1 && outputNodeInputNodes[0] == net.InputNodes[0] {
		inner := net.AllNodes[net.InputNodes[0]].ActivationFunction.goExpression(variableNames[0])
		return net.AllNodes[net.OutputNode].ActivationFunction.goExpression(inner)
	}
	// A special case, one input node that has many inputs, and that's it
	if len(net.Connected()) == (1 + len(outputNodeInputNodes)) {
		return net.CombineInputs(net.AllNodes[net.OutputNode].ActivationFunction, outputNodeInputNodes, variableNames)
	}
	// Not a special case, traverse from the output node to the input nodes, while gathering expression strings
	return "NOT IMPLEMENTED YET"
}

// GoFunction returns the source code for a Go function that does the same thing as this network
func (net *Network) GoFunction(variableNames ...string) string {
	// Make sure the slice of input variable names contains at least one string
	if len(variableNames) == 0 {
		variableNames = []string{"x"}
	}
	return WrapInFunction(net.GoExpression(variableNames...), variableNames...)
}
