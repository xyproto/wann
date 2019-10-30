package wann

import "strings"

// generateExpression will try to return the formula for the given node
// using the "x" string if .Value is set and no input nodes are available.
// returns true if the maximum number of evaluation loops is reached
// TODO: Rewrite this function
func (neuron *Neuron) generateExpression(variableLetter *rune, maxEvaluationLoops *int, callDepth int) (string, bool) {
	if *maxEvaluationLoops <= 0 {
		return "", true
	}
	// Assume this is the Output neuron, recursively evaluating the result
	// For each input neuron, evaluate them
	combined := ""
	counter := 0

	for _, inputNeuronIndex := range neuron.InputNodes {
		// Let each input neuron do its own evauluation, using the given weight
		(*maxEvaluationLoops)--
		// TODO: Figure out exactly why this one kicks in (and if it matters)
		//       It only seems to kick in during "go test" and not in evolve/main.go
		if int(inputNeuronIndex) >= len(neuron.Net.AllNodes) {
			continue
			//panic("TOO HIGH INPUT NEURON INDEX")
		}
		// TODO: Don't add the formulas, replace "X" and place them within each other
		vl := 'A'
		if callDepth == 0 {
			vl = *variableLetter
		}
		result, stopNow := neuron.Net.AllNodes[inputNeuronIndex].generateExpression(&vl, maxEvaluationLoops, callDepth+1)
		if callDepth == 0 {
			combined += result
		} else {
			combined = strings.Replace(combined, string(vl), result, -1)
		}
		counter++
		if stopNow || (*maxEvaluationLoops < 0) {
			break
		}
	}
	// No input neurons. Invent a variable name (using the input node counter?)
	if counter == 0 && neuron.Value != nil {
		returnString := string(*variableLetter) + ", "
		(*variableLetter)++
		return returnString, false
	}
	// Return the averaged sum, or 0
	if counter == 0 {
		return "", false
	}
	if callDepth != 0 {
		// X is used as a placeholder for inserting other expressions
		return neuron.ActivationFunction.goExpression(string('A')), false
	}
	// No further expressions, this is the end of the line, just use "x"
	return neuron.ActivationFunction.goExpression(string(*variableLetter)), false
}

// GoExpression returns the source code for a Go expression that does the same thing as this network
func (net *Network) GoExpression(varLetter rune) string {
	outputNode := net.AllNodes[net.OutputNode]
	maxIterationCounter := 100 // to avoid circular connections in the graph
	result, _ := outputNode.generateExpression(&varLetter, &maxIterationCounter, 0)
	return strings.Replace(result, ", )", ")", -1)
}

// GoFunction returns the source code for a Go function that does the same thing as this network
func (net *Network) GoFunction() string {
	return "func f(x float64) float64 { return " + net.GoExpression('x') + " }"
}
