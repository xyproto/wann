package main

import (
	"fmt"
	"github.com/xyproto/wann"
)

func main() {
	net := wann.NewNetwork(&wann.Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println(net)
}
