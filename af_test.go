package wann

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func ExampleActivationFunctionIndex_Call() {
	fmt.Println(Gauss.Call(2.0))
	// Output:
	// 0.13427659965015956
}

func ExampleGauss_Statement() {
	statement := Gauss.Statement(jen.Id("x"))
	fmt.Println(statement.GoString())
	result, err := RunStatementX(statement, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// math.Exp(-(math.Pow(x, 2.0)) / 2.0)
	// 0.8824969025845955
}

func ExampleActivationFunctionIndex_GoRun() {
	// Run the Gauss function directly
	fmt.Println(ActivationFunctions[Gauss](0.5))
	// Use Jennifer to generate a source file just for running the Gauss function, then use "go run" and fetch the result
	if result, err := Gauss.GoRun(0.5); err == nil { // no error
		fmt.Println(result)
	}
	// Output:
	// 0.8824699625576026
	// 0.8824969025845955
}
