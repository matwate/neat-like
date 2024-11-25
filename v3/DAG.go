package main

import "sort"

type Dag struct {
	nodes []*Node
}

func NewDag() *Dag {
	return &Dag{}
}

func (d *Dag) AddNode() {
	d.nodes = append(d.nodes, &Node{})
}

func (d *Dag) AddEdge(from, to int) bool {
	if !d.IsValid(from) || !d.IsValid(to) {
		return false
	}

	// Ensure no cycles
	if from == to || d.IsAncestor(from, to) || d.IsParent(from, to) {
		return false
	}

	d.nodes[from].outgoing = append(d.nodes[from].outgoing, to)
	d.nodes[to].incoming++

	return true
}

func (d *Dag) IsValid(node int) bool {
	return node < len(d.nodes) && node >= 0
}

func (d *Dag) IsParent(parent, child int) bool {
	if !d.IsValid(parent) || !d.IsValid(child) {
		return false
	}
	out := d.nodes[parent].outgoing
	return any(out, func(v int) bool {
		return v == child
	})
}

func (d *Dag) IsAncestor(ancestor, descendant int) bool {
	if !d.IsValid(ancestor) || !d.IsValid(descendant) {
		return false
	}
	out := d.nodes[ancestor].outgoing
	return any(out, func(v int) bool {
		return d.IsAncestor(v, descendant)
	})
}

func (d *Dag) ComputeNodeDepths() {
	// This function sets each node's depth variable to its correct value
	var count int = len(d.nodes)
	// Nodes with no incoming edge
	var startNodes Set[int] = NewSet[int]()
	// Current incoming edge state for all nodes
	var incoming []int = make([]int, count)

	for i, n := range d.nodes {
		incoming[i] = n.depth
	}

	for i, n := range d.nodes {
		if n.incoming == 0 {
			n.depth = 0
			startNodes.Add(i)
		}
	}
	for startNodes.Len() > 0 {
		// Extract a node from the starting set
		idx := startNodes.Back()
		startNodes.Pop_back()

		n := d.nodes[idx]
		for _, o := range n.outgoing {
			incoming[o]--
			connected := d.nodes[o]
			connected.depth = max(connected.depth, n.depth+1)
			if incoming[o] == 0 {
				startNodes.Add(o)
			}
		}
	}
}

func (d *Dag) getOrder() []int {
	d.ComputeNodeDepths()
	order := make([]int, len(d.nodes))
	sort.Slice(d.nodes, func(i, j int) bool {
		return d.nodes[i].depth < d.nodes[j].depth
	})
	return order
}

func (d *Dag) RemoveConnection(from, to int) bool {
	if !d.IsValid(from) || !d.IsValid(to) {
		return false
	}
	out := d.nodes[from].outgoing
	var found int = 0
	var i int = 0
	var count int = len(out)
	for i < count-found {
		if out[i] == to {
			out = append(out[:i], out[i+1:]...)
			found++
			d.nodes[to].incoming--
		}
	}
	if found == 0 {
		return false
	}
	return true
}

func any[T comparable](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}
