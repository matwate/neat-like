package main

import "fmt"

type Graph struct {
	adjacencyList map[int]Node
}

type Node struct {
	adjacents []int
	incoming  int
}

func NewGraph() *Graph {
	return &Graph{adjacencyList: make(map[int]Node)}
}

func (g *Graph) AddNode() {
	g.adjacencyList[len(g.adjacencyList)] = Node{adjacents: []int{}, incoming: 0}
}

func (g *Graph) AddEdge(from, to int) {
	if from >= len(g.adjacencyList) || to >= len(g.adjacencyList) || from < 0 || to < 0 {
		panic(fmt.Errorf("invalid connection from %d to %d", from, to))
	}
	fromNode := g.adjacencyList[from]
	toNode := g.adjacencyList[to]

	fromNode.adjacents = append(fromNode.adjacents, to)
	toNode.incoming++

	g.adjacencyList[from] = fromNode
	g.adjacencyList[to] = toNode
}

func (g *Graph) RemoveEdge(from, to int) {
	if from >= len(g.adjacencyList) || to >= len(g.adjacencyList) || from < 0 || to < 0 {
		panic(fmt.Errorf("invalid connection from %d to %d", from, to))
	}
	fromNode := g.adjacencyList[from]
	toNode := g.adjacencyList[to]

	for i, v := range fromNode.adjacents {
		if v == to {
			fromNode.adjacents = append(fromNode.adjacents[:i], fromNode.adjacents[i+1:]...)
			toNode.incoming--
			break
		}
	}

	g.adjacencyList[from] = fromNode
	g.adjacencyList[to] = toNode
}

func (g *Graph) TopoSort() []int {
	count := len(g.adjacencyList)
	noIncoming := []int{}
	incoming := make([]int, count)

	for i := 0; i < count; i++ {
		incoming[i] = g.adjacencyList[i].incoming
		if incoming[i] == 0 {
			noIncoming = append(noIncoming, i)
		}
	}

	var sorted []int
	for len(noIncoming) > 0 {
		node := noIncoming[0]
		noIncoming = noIncoming[1:]
		sorted = append(sorted, node)

		for _, adj := range g.adjacencyList[node].adjacents {
			incoming[adj]--
			if incoming[adj] == 0 {
				noIncoming = append(noIncoming, adj)
			}
		}
	}

	if len(sorted) != count {
		panic("graph has at least one cycle")
	}

	return sorted
}

func main() {
	g := NewGraph()
	g.AddNode()
	g.AddNode()
	g.AddNode()
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	fmt.Println(g.TopoSort())
}
