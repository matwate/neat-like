package main

import (
	"slices"
)

// This will be an implementation of directed acyclic graph, of dynamic size, using an adjacency matrix

type DAG[T comparable] struct {
	adjacencyMatrix [][]T
	nodeValues      []float64
	nullValue       T
	apply           func(float64, T) float64 // For using as a neural network
	apply_no_act    func(float64, T) float64 // For using without activation function
	size            int
	inputs, outputs int
}

// This will be here so i dont lose my sanity, this will be de NN pair, (weight, bias) that will be used in the neural network

// DAG constructor

func NewDAG[T comparable](
	inputs, outputs int,
	nullValue, initialValue T,
	apply func(float64, T) float64,
	apply_no_act func(float64, T) float64,
) DAG[T] {
	// The first N nodes will be the input nodes, the next M nodes will be the output nodes
	// The adjacency matrix will be of size N+M x N+M, so if we have 3 inputs and 2 outputs, the matrix will be 5x5
	// And nodes 0, 1, 2 will be the input nodes, and nodes 3, 4 will be the output nodes
	// Base is the base value of the adjacency matrix, it will be used to initialize the matrix, (i should change it to a function that returns somewhat random values)
	adjacencyMatrix := make([][]T, inputs+outputs)
	for i := range adjacencyMatrix {
		adjacencyMatrix[i] = make([]T, inputs+outputs)
		for j := 0; j < inputs+outputs; j++ {
			adjacencyMatrix[i][j] = nullValue
		}
	}
	dag := DAG[T]{
		adjacencyMatrix: adjacencyMatrix,
		nodeValues:      make([]float64, inputs+outputs),
		nullValue:       nullValue,
		apply:           apply,
		apply_no_act:    apply_no_act,
		size:            inputs + outputs,
		inputs:          inputs,
		outputs:         outputs,
	}
	return dag
}

// DAG methods
// AddNode will add a node, this will be treated as hidden layer for neat so no connections will be add
// AddEdge(T) will add a connection between two nodes, the value of the connection will be T
// TopoSort will return a list of integers that are the indices of the nodes in topological order
// We won't be removing nodes
// Helper function resize
//

func (dg *DAG[T]) AddNode() {
	// Resize the adjaceny Matrix to new size
	// Add the new node to the nodeValues
	// Add the new node to the adjacencyMatrix
	dg.resize(len(dg.adjacencyMatrix) + 1)
	// Ngl this is evertyhing i need to do
}

func (dg *DAG[T]) resize(newSize int) {
	if newSize < len(dg.adjacencyMatrix) {
		panic("Don't do that, you can't resize to a smaller size")
	}

	// Create a new adjacency matrix of size newSize
	newAdjacencyMatrix := make([][]T, newSize)
	for i := range newAdjacencyMatrix {
		newAdjacencyMatrix[i] = make([]T, newSize)
		for j := 0; j < newSize; j++ {
			newAdjacencyMatrix[i][j] = dg.nullValue
		}
	}
	// Copy all the data
	for i := 0; i < len(dg.adjacencyMatrix); i++ {
		for j := 0; j < len(dg.adjacencyMatrix); j++ {
			newAdjacencyMatrix[i][j] = dg.adjacencyMatrix[i][j]
		}
	}
	// Set the new adjacency Matrix
	dg.adjacencyMatrix = newAdjacencyMatrix
	// Resize the nodeValues
	newNodes := make([]float64, newSize)
	_ = copy(newNodes, dg.nodeValues)
	dg.nodeValues = newNodes
	dg.size = newSize
}

func (dg *DAG[T]) AddEdge(from, to int, value T) {
	// Add the edge to the adjacency matrix
	dg.adjacencyMatrix[from][to] = value
}

func (dg *DAG[T]) TopoSort() []int {
	adjCopy := make([][]T, len(dg.adjacencyMatrix))
	for i := range adjCopy {
		adjCopy[i] = make([]T, len(dg.adjacencyMatrix))
		copy(adjCopy[i], dg.adjacencyMatrix[i])
	}
	var topoOrder []int
	var marked []bool = make([]bool, len(dg.adjacencyMatrix))
	for len(topoOrder) < len(dg.adjacencyMatrix) {
		noOutIdx := dg.no_outgoing(adjCopy, marked)
		if noOutIdx == -1 {
			break
		}
		for i := 0; i < len(adjCopy); i++ {
			adjCopy[i][noOutIdx] = dg.nullValue
		}
		topoOrder = append(topoOrder, noOutIdx)
		marked[noOutIdx] = true
	}
	slices.Reverse(topoOrder)
	return topoOrder
}

func (dg *DAG[T]) no_outgoing(adj [][]T, marked []bool) int {
	for i := 0; i < len(adj); i++ {
		noOut := true
		if marked[i] {
			continue
		}
		for j := 0; j < len(adj); j++ {
			if adj[i][j] != dg.nullValue {
				noOut = false
				break
			}
		}
		if noOut {
			return i
		}
	}
	return -1
}

func (dg *DAG[T]) Fit() []float64 {
	// This will be adding the values to the nodes, in the topological order
	// This will NOT be modifying the nodeValues, it will be returning a new slice, since this is like a neural network, we can't modify the values of the nodes outside of training
	topoOrder := dg.TopoSort()
	newNodes := make([]float64, len(dg.nodeValues))
	_ = copy(newNodes, dg.nodeValues)
	for i := 0; i < len(topoOrder); i++ {
		next := topoOrder[i]
		// Now we're going to check for its connections in the adjacency matrix, use the apply function to get the value and place it at the node the connection is going to +=
		for conn := range dg.adjacencyMatrix[next] {
			if next < dg.inputs && dg.adjacencyMatrix[next][conn] != dg.nullValue {
				newNodes[conn] += dg.apply_no_act(
					dg.nodeValues[next],
					dg.adjacencyMatrix[next][conn],
				)
			}

			if dg.adjacencyMatrix[next][conn] != dg.nullValue {
				newNodes[conn] += dg.apply(dg.nodeValues[next], dg.adjacencyMatrix[next][conn])
			}
		}
	}
	return newNodes
}

func (dg *DAG[T]) SetNodeValue(node int, value float64) {
	dg.nodeValues[node] = value
}

func (dg *DAG[T]) SetEdgeValue(from, to int, value T) {
	dg.adjacencyMatrix[from][to] = value
}

func (dg *DAG[T]) Copy() DAG[T] {
	return DAG[T]{
		adjacencyMatrix: dg.adjacencyMatrix,
		nodeValues:      dg.nodeValues,
		nullValue:       dg.nullValue,
		apply:           dg.apply,
		apply_no_act:    dg.apply_no_act,
		size:            dg.size,
		inputs:          dg.inputs,
		outputs:         dg.outputs,
	}
}
