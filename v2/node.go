package main

type Node struct {
	activation Activation
	outgoing   []int
	incoming   int
	depth      int
}

func (n Node) GetConnectionCount() int {
	return len(n.outgoing)
}
