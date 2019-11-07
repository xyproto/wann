package wann

import (
	"github.com/dave/jennifer/jen"
)

// ActivationStatement creates an activation function statment, given a weight and input statements
// returns: activationFunction(input0 * w + input1 * w + ...)
// The function calling this function is responsible for inserting network input values into the network input nodes.
func ActivationStatement(af ActivationFunctionIndex, w float64, inputStatements []*jen.Statement) *jen.Statement {
	// activationFunction(input0 * w + input1 * w + ...)
	weightedSum := jen.Empty()
	for i, inputStatement := range inputStatements {
		if i == 0 {
			// first
			weightedSum.Add(inputStatement).Op("*").Lit(w)
		} else {
			// the rest, same as above, but with a leading "+"
			weightedSum.Op("+").Add(inputStatement).Op("*").Lit(w)
		}
	}
	return af.Statement(weightedSum)
}
