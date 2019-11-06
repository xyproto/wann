package wann

import (
	"github.com/dave/jennifer/jen"
)

// render renders a *jen.Statement to a string, if possible
// if there is an error about an extra ")", then that's because anonymous functions are not supported by jen
// Do not render until statements could be placed at the top-level in a Go program.
func render(inner *jen.Statement) string {
	return inner.GoString()
}

// JenniferOutputNodeX returns a statement for the output node, using "x" for the variable
func (net *Network) JenniferOutputNodeX(functionName string) string {
	inner := net.AllNodes[net.OutputNode].ActivationFunction.Statement(jen.Id("x"))
	f := jen.Id(functionName).Op(":=").Add(inner)
	return render(f)
}
