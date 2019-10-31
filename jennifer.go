package wann

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
)

// render renders a *jen.Statement to a string, if possible
// if there is an error about an extra ")", then that's because anonymous functions are not supported by jen
// Do not render until statements could be placed at the top-level in a Go program.
func render(inner *jen.Statement) string {
	return inner.GoString()
}

// run will run the given statement by wrapping it in a program and using "go run"
func run(statement *jen.Statement, x float64) (float64, error) {
	// TODO: Use the standard library for aquiring a temporary filename
	filename := "/tmp/main.go"
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Id("x").Op(":=").Lit(x),
		jen.Qual("fmt", "Println").Call(statement),
	)
	err := ioutil.WriteFile(filename, []byte(f.GoString()), 0664)
	if err != nil {
		return 0.0, err
	}
	cmd := exec.Command("go", "run", filename)
	out, err := cmd.CombinedOutput()
	resultString := strings.TrimSpace(string(out))
	resultFloat, err := strconv.ParseFloat(resultString, 64)
	if err != nil {
		return 0.0, err
	}
	if err := os.Remove(filename); err != nil {
		return 0.0, err
	}
	return resultFloat, nil
}

// JenniferOutputNodeX returns a statement for the output node, using "x" for the variable
func (net *Network) JenniferOutputNodeX(functionName string) string {
	inner := net.AllNodes[net.OutputNode].ActivationFunction.Statement(jen.Id("x"))
	f := jen.Id(functionName).Op(":=").Add(inner)
	return render(f)
}
