package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Genome struct {
	dag    DAG
	nodes  map[int]*GenomeNode // This maps node numbers to node information (activation function and connection weights and biases
	input  int
	output int
	hidden int
}

type GenomeNode struct {
	activationFunction Activation
	connections        map[int]GenomeConnection
	Type               NodeType
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
		dag:    NewDAG(),
		nodes:  make(map[int]*GenomeNode),
		input:  0,
		output: 0,
		hidden: 0,
	}
}

func (g *Genome) Init(input, output int) {
	// This will get a new genoms with n input nodes and m output nodes
	for i := 0; i < input; i++ {
		g.AddNode(Input)
	}

	for i := 0; i < output; i++ {
		g.AddNode(Output)
	}

	g.input = input
	g.output = output
	// For the first generation, we have to connect every input node to every output node
	for i := 0; i < input; i++ {
		for j := input; j < input+output; j++ {
			g.AddConnection(i, j)
		}
	}
}

func (g *Genome) AddNode(node_type NodeType) {
	g.dag.createNode()
	var activation Activation
	switch node_type {
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
	// Generate the new node that we're going to add
	newNode := GenomeNode{
		activationFunction: activation,
		connections:        make(map[int]GenomeConnection),
		Type:               node_type,
	}
	// Add the new node to the genome
	g.nodes[len(g.dag.Nodes)-1] = &newNode
}

func (g *Genome) AddConnection(from, to int, value ...float64) {
	if !g.dag.IsValid(from) || !g.dag.IsValid(to) {
		return
	}
	g.dag.createConnection(from, to)
	var newWeight, newBias float64
	if len(value) != 2 || value[0] == 0 || value[1] == 0 {
		newWeight = rand.Float64()
		newBias = rand.Float64()
	} else {
		newWeight = value[0]
		newBias = value[1]
	}
	// Generate the new connection that we're going to add
	newConnection := GenomeConnection{
		weight: newWeight,
		bias:   newBias,
	}
	// Add the new connection to the node
	g.nodes[from].connections[to] = newConnection
}

func (g *Genome) RemoveConnection(from, to int) {
	if !g.dag.IsValid(from) || !g.dag.IsValid(to) {
		return
	}
	g.dag.RemoveConnection(from, to)
	delete(g.nodes[from].connections, to)
}

func (g *Genome) Forward(inputs ...float64) []float64 {
	// This function will take a set of values and return what it outputs
	if len(inputs) != g.input {
		panic("Input size mismatch")
	}
	// Get the evaluation order
	order := g.dag.setOrder()
	// Set the input values
	values := make([]float64, len(g.nodes))
	for i, v := range inputs {
		values[i] = v
	}
	// Evaluate the nodes
	for _, node := range order {
		n := g.nodes[node]
		for to, conn := range n.connections {
			values[to] += n.ApplyActivation(values[node]*conn.weight + conn.bias)
		}
	}
	// The code above doesnt apply the activation function to the output nodes
	for i := g.input; i < g.input+g.output; i++ {
		n := g.nodes[i]
		values[i] = n.ApplyActivation(values[i])
	}
	// Get the output values
	outputs := make([]float64, g.output)
	for i := 0; i < g.output; i++ {
		outputs[i] = values[g.input+i]
	}
	return outputs
}

func (g *Genome) Print() {
	g.dag.Print()
	for k, v := range g.nodes {
		fmt.Printf("Node %d: %v\n", k, v)
	}
}

func (g *Genome) SetWeight(from, to int, weight float64) {
	if !g.dag.IsValid(from) || !g.dag.IsValid(to) {
		return
	}
	hasConnection := g.dag.hasConnection(from, to)
	if hasConnection == -1 {
		panic("Connection not found")
	}
	newEntry := g.nodes[from].connections[to]
	newEntry.weight = weight
	g.nodes[from].connections[to] = newEntry
}

func (g *Genome) SetBias(from, to int, bias float64) {
	if !g.dag.IsValid(from) || !g.dag.IsValid(to) {
		return
	}
	hasConnection := g.dag.hasConnection(from, to)
	if hasConnection == -1 {
		panic("Connection not found")
	}
	newEntry := g.nodes[from].connections[to]
	newEntry.bias = bias
	g.nodes[from].connections[to] = newEntry
}

func (gn *GenomeNode) ApplyActivation(x float64) float64 {
	switch gn.activationFunction {
	case None:
		return x
	case Sigmoid:
		return 1 / (1 + math.Exp(-x))
	case Tanh:
		return math.Tanh(x)
	case ReLU:
		return math.Max(0, x)
	}
	panic("Unknown activation function")
}

func (g *Genome) HasConnection(from, to int) bool {
	if !g.dag.IsValid(from) || !g.dag.IsValid(to) {
		return false
	}
	return g.dag.hasConnection(from, to) != -1
}
