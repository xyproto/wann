package wann

import (
	"fmt"
	"math/rand"
)

func ExampleNetwork_GoExpression() {
	rand.Seed(999)
	net := NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	formula := net.GoExpression('z')
	fmt.Println(formula)

	rand.Seed(1111113)
	net = NewNetwork(&Config{
		inputs:                 5,
		InitialConnectionRatio: 0.7,
		sharedWeight:           0.5,
	})
	formula = net.GoExpression('x')
	fmt.Println(formula)

	// Output:
	// math.Exp(-(z * z) / 2.0)
	// func(r float64) float64 { if r >= 0 { return r } else { return 0 } }(x)
}
