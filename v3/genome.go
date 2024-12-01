package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Genome struct {
	dag    Dag
	nodes  map[int]*GenomeNode // This maps node numbers to node information (activation function and connection weights and biases
	input  int
	output int
	hidden int
}

type GenomeNode struct {
	activationFunction activationFun
	connections        map[int]GenomeConnection
	Type               NodeType
}

func (gn *GenomeNode) applyActivation(x float64) float64 {
	switch activation := gn.activationFunction; activation {
	case None:
		return x
	case ReLU:
		return math.Max(0, x)
	case Sigmoid:
		return 1 / (1 + math.Exp(-x))
	}
	// unreachable
	panic("Unknown activation function")
}

type NodeType int

const (
	Input NodeType = iota
	Output
	Hidden
)

type GenomeConnection struct {
	weight, bias float64
}

func NewGenome() Genome {
	return Genome{
		dag:    NewDag(),
		nodes:  make(map[int]*GenomeNode),
		input:  0,
		output: 0,
		hidden: 0,
	}
}

func (g *Genome) Init(input, output int) {
	// This will make a genome with input input nodes and output output nodes
	for range input {
		g.AddNode(Input)
	}

	for range output {
		g.AddNode(Output)
	}

	g.input = input
	g.output = output

	for i := range input {
		for j := range output {
			g.AddConnection(i, j)
		}
	}
}

func (g *Genome) AddNode(nt NodeType) {
	g.dag.AddNode()
	var activation activationFun
	switch nt {
	case Input:
		g.input++
		activation = None
	case Output:
		g.output++
		activation = ReLU
	case Hidden:
		g.hidden++
		activation = ReLU
	}
	newNode := GenomeNode{
		activationFunction: activation,
		connections:        make(map[int]GenomeConnection),
		Type:               nt,
	}
	g.nodes[len(g.dag.nodes)-1] = &newNode
}

func (g *Genome) AddConnection(from, to int, values ...float64) {
	switch g.dag.AddEdge(from, to) {
	case false:
		// panic("Connection already exists")
		return
	default:
		if len(values) != 2 || values[0] != 0 && values[1] != 0 {
			g.nodes[from].connections[to] = GenomeConnection{
				weight: rand.Float64(),
				bias:   rand.Float64(),
			}
		} else {
			g.nodes[from].connections[to] = GenomeConnection{
				weight: values[0],
				bias:   values[1],
			}
		}
	}
}

func (g *Genome) RemoveConnection(from, to int) {
	g.dag.RemoveConnection(from, to)
	delete(g.nodes[from].connections, to)
}

func (g *Genome) Forward(inputs ...float64) []float64 {
	assert_equal(len(inputs), g.input)
	// Get the evaluation order
	order := g.dag.getOrder()
	fmt.Println(order)
	// Make a values list to store the values of each node
	var values []float64 = make([]float64, len(g.dag.nodes))
	for i := range inputs {
		values[i] = inputs[i]
	}
	fmt.Println(values)

	// Iterate over the order and calculate the values of each nodes
	for _, node := range order {
		n := g.nodes[node]
		for to, conn := range n.connections {
			values[to] += n.applyActivation(values[node]*conn.weight + conn.bias)
		}
	}

	// Apply the activation function to the output nodes
	for i := g.input; i < g.input+g.output; i++ {
		values[i] = g.nodes[i].applyActivation(values[i])
	}
	// Return the output values
	output := make([]float64, g.output)
	for i := range output {
		output[i] = values[g.input+i]
	}
	return output
}

func (g *Genome) SetWeight(from, to int, weight float64) {
	if !g.dag.hasConnection(from, to) {
		panic("Connection does not exist")
	}
	g.nodes[from].connections[to] = GenomeConnection{
		weight: weight,
		bias:   g.nodes[from].connections[to].bias,
	}
}

func (g *Genome) SetBias(from, to int, bias float64) {
	if !g.dag.hasConnection(from, to) {
		panic("Connection does not exist")
	}
	g.nodes[from].connections[to] = GenomeConnection{
		weight: g.nodes[from].connections[to].weight,
		bias:   bias,
	}
}

func assert_equal[i comparable](a, b i) {
	if a != b {
		panic(fmt.Sprintf("Expected %v, got %v", a, b))
	}
}
