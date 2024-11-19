package main

import (
	"math"
	"math/rand"
)

// Here we will implement a variation of neat algorithm

// First the definitions and the accessable stuff.
type NNPair struct {
	weight, bias float64
}

type Neatlike struct {
	// It will have a DAG and some extra methods.
	dag             DAG[NNPair]
	inputs, outputs int
}

// Neatlike constructor
func NewNeatlike(inputs, outputs int) Neatlike {
	// The DAG will have n inputs and m outputs, and the values will be NNPair
	// The null value will be a weight of 0 and a bias of 0
	// The apply function will be tanh
	dag := NewDAG(
		inputs, outputs,
		NNPair{0, 0},
		NNPair{1, 0},
		func(x float64, y NNPair) float64 {
			return math.Tanh(x*y.weight + y.bias)
		},
		func(x float64, y NNPair) float64 {
			return x*y.weight + y.bias
		},
	)
	// Connect the input nodes to the output nodes, with random values normally distributed
	for i := 0; i < inputs; i++ {
		for j := inputs; j < inputs+outputs; j++ {
			dag.AddEdge(i, j, NNPair{rand.NormFloat64(), rand.NormFloat64()})
		}
	}

	// Give random values to the input nodes
	for i := 0; i < inputs; i++ {
		dag.nodeValues[i] = rand.NormFloat64()
	}
	return Neatlike{dag, inputs, outputs}
}

func (n *Neatlike) Copy() Neatlike {
	return Neatlike{
		n.dag.Copy(),
		n.inputs,
		n.outputs,
	}
}

// Neatlike functions will be:
// Fit,  Mutate.
// Fit will be the same as in the neural network, but it will be using the DAG, it has an input of len n and an output of len
// Mutate is more complicated, it will be either structural or parametric, structural will add a node or an edge to the DAG, parametric will tweak the weights and biases of the DAG

func (n *Neatlike) Fit(inputs []float64) []float64 {
	if len(inputs) != n.inputs {
		panic("Invalid input size")
	}
	// Reset the values of the nodes
	for i := 0; i < n.dag.size; i++ {
		n.dag.nodeValues[i] = 0
	}

	// Set the input nodes
	for i := 0; i < n.inputs; i++ {
		n.dag.nodeValues[i] = inputs[i]
	}

	result := n.dag.Fit()

	return result[n.inputs : n.inputs+n.outputs] // This returns the output nodes
}

func (n *Neatlike) Mutate() {
	const (
		structural   = 0
		parametric   = 1
		mutationRate = 0.1
	)
	switch rand.Intn(2) {
	case structural:
		// Add a node or an edge
		switch rand.Intn(2) {
		case 0:
			// Add a node
			node1 := rand.Intn(n.dag.size)
			attempts := 0
			for node1 >= n.inputs && node1 < n.inputs+n.outputs {
				node1 = rand.Intn(n.dag.size)
				attempts++
				if attempts > n.dag.size {
					return
				}
			}
			selectedNode := node1
			conn := rand.Intn(n.dag.size)
			attempts = 0
			for n.dag.adjacencyMatrix[selectedNode][conn] == n.dag.nullValue {
				conn = rand.Intn(n.dag.size)
				attempts++
				if attempts > n.dag.size {
					return
				}
			}

			newWeight := math.Sqrt(n.dag.adjacencyMatrix[selectedNode][conn].weight)
			newBias1 := (rand.Float64()*2 - 1) * mutationRate
			newBias2 := (n.dag.adjacencyMatrix[selectedNode][conn].bias - newBias1) / newWeight

			n.dag.SetEdgeValue(selectedNode, conn, n.dag.nullValue)
			n.dag.AddNode()
			n.dag.AddEdge(
				selectedNode,
				n.dag.size-1,
				NNPair{newWeight, newBias1},
			)
			n.dag.AddEdge(n.dag.size-1, conn, NNPair{newWeight, newBias2})

		case 1:
			// Add an edge
			node1 := rand.Intn(n.dag.size)
			attempts := 0
			for node1 >= n.inputs && node1 < n.inputs+n.outputs {
				node1 = rand.Intn(n.dag.size)
				attempts++
				if attempts > n.dag.size {
					return
				}
			}
			selectedNode := node1
			conn := rand.Intn(n.dag.size)
			attempts = 0
			for n.dag.adjacencyMatrix[selectedNode][conn] != n.dag.nullValue || conn < n.inputs {
				conn = rand.Intn(n.dag.size)
				attempts++
				if attempts > n.dag.size {
					return
				}
			}
			if n.dag.HasPath(conn, selectedNode) || selectedNode == conn {
				return
			}
			n.dag.AddEdge(
				selectedNode,
				conn,
				NNPair{rand.NormFloat64(), rand.NormFloat64()},
			)
		}
	case parametric:
		// Select a random node
		node := rand.Intn(n.dag.size)
		attempts := 0
		for node >= n.inputs && node < n.inputs+n.outputs {
			node = rand.Intn(n.dag.size)
			attempts++
			if attempts > n.dag.size {
				return
			}
		}
		conn := rand.Intn(n.dag.size)
		attempts = 0
		for n.dag.adjacencyMatrix[node][conn] == n.dag.nullValue {
			conn = rand.Intn(n.dag.size)
			attempts++
			if attempts > n.dag.size {
				return
			}
		}
		change := (rand.Float64()*2 - 1) * mutationRate
		switch rand.Intn(2) {
		case 0:
			n.dag.SetEdgeValue(
				node,
				conn,
				NNPair{
					n.dag.adjacencyMatrix[node][conn].weight + change,
					n.dag.adjacencyMatrix[node][conn].bias,
				},
			)
		case 1:
			n.dag.SetEdgeValue(
				node,
				conn,
				NNPair{
					n.dag.adjacencyMatrix[node][conn].weight,
					n.dag.adjacencyMatrix[node][conn].bias + change,
				},
			)
		}
	}
}

// In the DAG implementation, add the HasPath function to check for cycles
func (dag *DAG[NNPair]) HasPath(from, to int) bool {
	visited := make([]bool, dag.size)
	var dfs func(int) bool
	dfs = func(u int) bool {
		if u == to {
			return true
		}
		visited[u] = true
		for v, value := range dag.adjacencyMatrix[u] {
			if value != dag.nullValue && !visited[v] {
				if dfs(v) {
					return true
				}
			}
		}
		return false
	}
	return dfs(from)
}
