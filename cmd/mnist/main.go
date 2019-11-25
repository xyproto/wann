package main

import (
	"fmt"
	"log"

	"github.com/JoshVarga/GoMNIST"
	"github.com/xyproto/wann"
)

func main() {

	// WORK IN PROGRESS

	// Read label file
	ll, err := GoMNIST.ReadLabelFile("data/t10k-labels-idx1-ubyte.gz")
	if err != nil {
		log.Fatalf("read (%s)", err)
	}
	if len(ll) != 10000 {
		err := fmt.Errorf("unexpected count %d", len(ll))
		panic(err)
	}

	// Read image file
	nrow, ncol, imgs, err := GoMNIST.ReadImageFile("data/t10k-images-idx3-ubyte.gz")
	if err != nil {
		log.Fatalf("read (%s)", err)
	}
	if len(imgs) != 10000 {
		err := fmt.Errorf("unexpected count %d", len(imgs))
		panic(err)
	}
	fmt.Printf("%d images, %dx%d format\n", len(imgs), nrow, ncol)

	// Read data?
	train, test, err := GoMNIST.Load("./data")
	if err != nil {
		log.Fatalf("load (%s)", err)
	}
	println(train.Count(), test.Count())

	// Prepare a neural network configuration struct
	config := &wann.Config{
		InitialConnectionRatio: 0.2,
		Generations:            2000,
		PopulationSize:         500,
		Verbose:                true,
	}

	fmt.Println("CONFIG", config)

	// Sweeper next
	train, _, err = GoMNIST.Load("./data")
	if err != nil {
		log.Fatalf("load (%s)", err)
	}
	sweeper := train.Sweep()
	sweeper.Next()
}
