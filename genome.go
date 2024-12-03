package main

type Genome struct {
	graph  *Graph
	nodes  map[int]gNode
	input  int
	output int
	hidden int
}

type gNode struct {
	activation ActivationFunction
	edges      map[int]gEdge
}

type gEdge struct {
	weight float64
	bias   float64
}

type ActivationFunction int

const (
	None ActivationFunction = iota
	Sigmoid
	Tanh
)
