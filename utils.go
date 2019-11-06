package wann

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
)

// Pair is used for sorting dictionaries by value.
// Thanks https://stackoverflow.com/a/18695740/131264
type Pair struct {
	Key   int
	Value float64
}

// PairList is a slice of Pair
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SortByValue sorts a map[int]float64 by value
func SortByValue(m map[int]float64) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// checkInputNeurons was used for debugging
func (net *Network) checkInputNeurons() {
	for neuronIndex, neuron := range net.AllNodes {
		if len(net.AllNodes) != len(neuron.Net.AllNodes) {
			panic("net.AllNodes and neuron.Net.AllNodes have different length")
		}
		if net != neuron.Net {
			//panic("neuron Net pointer is out of sync")
			net.AllNodes[neuronIndex].Net = net
		}
		neuron = net.AllNodes[neuronIndex]
		if net != neuron.Net {
			panic("neuron Net pointer is out of sync")
		}
		for _, inputNeuronIndex := range neuron.InputNodes {
			if int(inputNeuronIndex) >= len(net.AllNodes) {
				fmt.Println("Network:", net.String())
				panic(fmt.Sprintf("indexNeuronIndex out of range: %d\n", inputNeuronIndex))
			}
			if int(inputNeuronIndex) >= len(neuron.Net.AllNodes) {
				panic(fmt.Sprintf("indexNeuronIndex out of range: %d\n", inputNeuronIndex))
			}
		}
	}
}

// runStatement will run the given statement by wrapping it in a program and using "go run"
func runStatement(statement *jen.Statement, x float64) (float64, error) {
	file, err := ioutil.TempFile("/tmp", "af_*.go")
	if err != nil {
		return 0.0, err
	}
	filename := file.Name()
	defer os.Remove(filename)
	// Build the contents of the source file using jennifer
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Id("x").Op(":=").Lit(x),
		jen.Qual("fmt", "Println").Call(statement),
	)
	// Save the file
	if ioutil.WriteFile(filename, []byte(f.GoString()), 0664) != nil {
		return 0.0, err
	}
	// Run the file
	cmd := exec.Command("go", "run", filename)
	out, err := cmd.CombinedOutput()
	// Return the outputted float string as a float64
	resultString := strings.TrimSpace(string(out))
	resultFloat, err := strconv.ParseFloat(resultString, 64)
	if err != nil {
		return 0.0, err
	}
	return resultFloat, nil
}
