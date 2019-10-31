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
	result, err := run(statement, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// math.Exp(-(math.Pow(x, 2.0)) / 2.0)
	// 0.8824969025845955
}
