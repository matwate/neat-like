package main

import (
	"fmt"
	"sort"

	"github.com/matwate/neat-like/v2/set"
)

type DAG struct {
	Nodes []Node
}

func NewDAG() DAG {
	return DAG{}
}

func (dg *DAG) createNode() {
	currentLen := len(dg.Nodes)
	Nodes := make([]Node, currentLen+1)
	copy(Nodes, dg.Nodes)
	dg.Nodes = Nodes
	// I hope go is smart and frees up old nodes.
}

// Modified createConnection to ensure stricter validation
func (dg *DAG) createConnection(from, to int) bool {
	if !dg.IsValid(from) || !dg.IsValid(to) {
		return false
	}

	// No cycles
	if from == to || dg.IsAncestor(from, to) || dg.IsParent(from, to) {
		return false
	}

	// Temporarily add the connection
	dg.Nodes[from].outgoing = append(dg.Nodes[from].outgoing, to)
	dg.Nodes[to].incoming++

	// Check if this created a cycle
	if dg.HasCycle() {
		// If it did, remove the connection and return false
		dg.Nodes[from].outgoing = dg.Nodes[from].outgoing[:len(dg.Nodes[from].outgoing)-1]
		dg.Nodes[to].incoming--
		return false
	}

	return true
}

func (dg *DAG) IsValid(node int) bool {
	return node >= 0 && node < len(dg.Nodes)
}

func (dg *DAG) IsParent(node_1, node_2 int) bool {
	if !dg.IsValid(node_1) || !dg.IsValid(node_2) {
		return false
	}
	out := dg.Nodes[node_1].outgoing
	return any(out, func(x int) bool { return x == node_2 })
}

func (dg *DAG) ComputeDepth() {
	node_count := len(dg.Nodes)
	if node_count == 0 {
		return
	}

	// Initialize incoming count array
	incoming := make([]int, node_count)

	// First pass: validate all connections and count incoming edges
	for i, node := range dg.Nodes {
		// Validate outgoing connections
		validOutgoing := make([]int, 0, len(node.outgoing))
		for _, out := range node.outgoing {
			if out >= 0 && out < node_count {
				validOutgoing = append(validOutgoing, out)
				incoming[out]++
			}
		}
		// Update node with only valid connections
		dg.Nodes[i].outgoing = validOutgoing
	}

	// Initialize no_incoming set for nodes with no incoming edges
	no_incoming := set.NewSet[int]()
	for i, node := range dg.Nodes {
		if node.incoming == 0 {
			node.depth = 0
			dg.Nodes[i] = node
			no_incoming.Add(i)
		}
	}

	// Process nodes in topological order
	processed := 0
	for no_incoming.Len() > 0 && processed < node_count {
		// Extract a node from the starting set
		node := no_incoming.Back()
		no_incoming.Pop_back()
		processed++

		n := dg.Nodes[node]
		for _, v := range n.outgoing {
			if v >= node_count {
				continue // Skip invalid connections
			}

			incoming[v]--
			connected := dg.Nodes[v]
			connected.depth = max(connected.depth, n.depth+1)
			dg.Nodes[v] = connected

			if incoming[v] == 0 {
				no_incoming.Add(v)
			}
		}
	}

	// Handle any remaining nodes (in case of cycles)
	if processed < node_count {
		for i := range dg.Nodes {
			if incoming[i] > 0 {
				dg.Nodes[i].depth = processed // Assign remaining nodes to a safe depth
			}
		}
	}
}

// Helper function to validate node connections
func (dg *DAG) validateConnections() {
	for i := range dg.Nodes {
		validOutgoing := make([]int, 0)
		for _, out := range dg.Nodes[i].outgoing {
			if dg.IsValid(out) {
				validOutgoing = append(validOutgoing, out)
			}
		}
		dg.Nodes[i].outgoing = validOutgoing
	}
}

func (dg *DAG) setOrder() []int {
	order := make([]int, len(dg.Nodes))
	for i := range dg.Nodes {
		order[i] = i
	}
	dg.ComputeDepth()
	sort.Slice(order, func(i, j int) bool {
		return dg.Nodes[i].depth < dg.Nodes[j].depth
	})

	return order
}

func (dg *DAG) hasConnection(from, to int) int {
	if !dg.IsValid(from) || !dg.IsValid(to) {
		return -1
	}
	for i, v := range dg.Nodes[from].outgoing {
		if v == to {
			return i
		}
	}
	return -1
}

func (dg *DAG) RemoveConnection(from, to int) {
	if !dg.IsValid(from) || !dg.IsValid(to) {
		return // or handle the error appropriately
	}
	connections := dg.Nodes[from].outgoing
	count := len(connections)
	found := 0
	for i, v := range connections {
		// If we find the connection, we remove it
		if v == to {
			// We swap the last element with the found element
			connections[i] = connections[count-1]
			// We remove the last element
			dg.Nodes[from].outgoing = connections[:count-1]
			found = 1
		}
	}
	if found == 0 {
		panic("Connection not found")
	}
}

func (dg *DAG) Print() {
	for i, v := range dg.Nodes {
		fmt.Printf("Node %d: %v\n", i, v)
	}
}

func (dg *DAG) ToDagVIS() []string {
	vis := make([]string, len(dg.Nodes))
	for i, node := range dg.Nodes {
		vis[i] = fmt.Sprintf("%v ", node.outgoing)
	}
	return vis
}

func any[T comparable](slice []T, f func(T) bool) bool {
	for _, v := range slice {
		if f(v) {
			return true
		}
	}
	return false
}

func (dg *DAG) IsAncestor(node_1, node_2 int) bool {
	if !dg.IsValid(node_1) || !dg.IsValid(node_2) {
		return false
	}
	// Use a visited set to prevent infinite recursion
	visited := make(map[int]bool)
	return dg.isAncestorHelper(node_1, node_2, visited)
}

func (dg *DAG) isAncestorHelper(node_1, node_2 int, visited map[int]bool) bool {
	if !dg.IsValid(node_1) || !dg.IsValid(node_2) {
		return false
	}
	// If we've already visited this node, return false to break potential cycles
	if visited[node_1] {
		return false
	}

	// Mark current node as visited
	visited[node_1] = true

	// Check if node_1 is a direct parent of node_2
	if dg.IsParent(node_1, node_2) {
		return true
	}

	// Check all outgoing connections recursively
	out := dg.Nodes[node_1].outgoing
	return any(out, func(x int) bool {
		return dg.isAncestorHelper(x, node_2, visited)
	})
}

// You might also want to add this validation function to check for cycles
func (dg *DAG) HasCycle() bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool)

	// Check from each node in case graph is not fully connected
	for i := range dg.Nodes {
		if !visited[i] {
			if dg.hasCycleHelper(i, visited, recStack) {
				return true
			}
		}
	}
	return false
}

func (dg *DAG) hasCycleHelper(node int, visited, recStack map[int]bool) bool {
	// If node is already in recursion stack, we found a cycle
	if recStack[node] {
		return true
	}

	// If node was already visited and isn't in recursion stack, no cycle here
	if visited[node] {
		return false
	}

	// Mark node as visited and add to recursion stack
	visited[node] = true
	recStack[node] = true

	// Check all neighbors
	for _, neighbor := range dg.Nodes[node].outgoing {
		if dg.hasCycleHelper(neighbor, visited, recStack) {
			return true
		}
	}

	// Remove node from recursion stack
	recStack[node] = false
	return false
}
