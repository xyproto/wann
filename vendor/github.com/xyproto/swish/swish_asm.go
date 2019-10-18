//+build amd64

package swish

// SwishAssembly is the swish function, written in hand-optimized assembly
// go: noescape
func SwishAssembly(x float64) float64
