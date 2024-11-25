package main

// Structural node
type Node struct {
	activation activationFun // activation function
	outgoing   []int         // outgoing edges
	incoming   int           // Amount of incoming edges
	depth      int           //
}
