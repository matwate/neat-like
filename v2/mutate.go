package main

import (
	"math/rand"
)

// this will be genome functions to mutate it

func (g *Genome) MutateAddNode() {
	// Pick a random node to split, it has to be an Input or Hidden node
	// Since nodes are layed out like this: [input, .. ouput, .. hidden]
	// We can just pick a random number between 0 and input+hidden and if its more than input, we add output to it to skip the output nodes
	splitNode := rand.Intn(g.input + g.hidden)
	if splitNode >= g.input {
		splitNode += g.output
	}

	// Now we pick one of the outgoing connections
	rand_conn := g.dag.Nodes[splitNode].outgoing[rand.Intn(len(g.dag.Nodes[splitNode].outgoing))]
	// So, we have to convert splitNode -> rand_conn connection to splitNode -> newNode -> rand_conn
	// So we add a new node
	if !g.dag.IsValid(splitNode) || !g.dag.IsValid(rand_conn) {
		return
	}

	if g.dag.IsParent(rand_conn, splitNode) || g.dag.IsAncestor(splitNode, rand_conn) {
		return
	}

	g.AddNode(Hidden)
	// We add a connection from the splitNode to the new Node, since we just added it, it is the last one.
	if !g.dag.IsValid(splitNode) || !g.dag.IsValid(len(g.dag.Nodes)-1) {
		return
	}
	g.AddConnection(splitNode, len(g.dag.Nodes)-1)
	// We add a connection from the new Node to the rand_conn node
	if !g.dag.IsValid(len(g.dag.Nodes)-1) || !g.dag.IsValid(rand_conn) {
		return
	}
	g.AddConnection(len(g.dag.Nodes)-1, rand_conn)
	// And we remove the connection from splitNode to rand_conn
	g.SetWeight(splitNode, len(g.dag.Nodes)-1, g.nodes[splitNode].connections[rand_conn].weight)
	g.SetWeight(len(g.dag.Nodes)-1, rand_conn, 1)
	// Set new conn bias to 0 the other one to the one before
	g.SetBias(len(g.dag.Nodes)-1, rand_conn, 0)
	g.SetBias(splitNode, len(g.dag.Nodes)-1, g.nodes[splitNode].connections[rand_conn].bias)
	g.RemoveConnection(splitNode, rand_conn)
}

func (g *Genome) MutateAddConnection() {
	// Pick a random node to connect from
	fromNode := rand.Intn(g.input + g.hidden)
	if fromNode >= g.input && fromNode <= g.input+g.output {
		fromNode += g.output
	}
	// Pick a random node to connect to
	toNode := rand.Intn(g.output+g.hidden) + g.input
	// If the connection already exists, we do nothing
	if g.HasConnection(fromNode, toNode) {
		return
	}
	if g.dag.IsParent(toNode, fromNode) || g.dag.IsAncestor(fromNode, toNode) {
		return
	}

	if !g.dag.IsValid(fromNode) || !g.dag.IsValid(toNode) {
		return
	}
	// We add the connection
	g.AddConnection(fromNode, toNode)
}

func (g *Genome) MutateChangeWeight() {
	// Pick a random node to connect from
	fromNode := rand.Intn(g.input + g.hidden)
	if fromNode >= g.input && fromNode <= g.input+g.output {
		fromNode += g.output
	}
	// Pick a random node to connect to
	toNode := rand.Intn(g.output+g.hidden) + g.input
	// If the connection does not exist, we do nothing
	if !g.HasConnection(fromNode, toNode) {
		return
	}
	if !g.dag.IsValid(fromNode) || !g.dag.IsValid(toNode) {
		return
	}
	// We change the weight, 75% chance to modify it by a random value, 25% chance to set it to a random value
	if rand.Float64() < 0.75 {
		// Modify the weight by a random value between -0.1 and 0.1
		g.SetWeight(
			fromNode,
			toNode,
			g.nodes[fromNode].connections[toNode].weight+(rand.Float64()-0.5)*0.2,
		)
	} else {
		// Set the weight to a random value between -1 and 1
		g.SetWeight(fromNode, toNode, rand.Float64()*2-1)
	}
}

func (g *Genome) MutateChangeBias() {
	// Pick a random node to connect from
	fromNode := rand.Intn(g.input + g.hidden)
	if fromNode >= g.input && fromNode <= g.input+g.output {
		fromNode += g.output
	}
	// Pick a random node to connect to
	toNode := rand.Intn(g.output+g.hidden) + g.input
	// If the connection does not exist, we do nothing
	if !g.HasConnection(fromNode, toNode) {
		return
	}
	if !g.dag.IsValid(fromNode) || !g.dag.IsValid(toNode) {
		return
	}
	// We change the weight, 75% chance to modify it by a random value, 25% chance to set it to a random value
	if rand.Float64() < 0.75 {
		// Modify the weight by a random value between -0.1 and 0.1
		g.SetBias(
			fromNode,
			toNode,
			g.nodes[fromNode].connections[toNode].bias+(rand.Float64()-0.5)*0.2,
		)
	} else {
		// Set the weight to a random value between -1 and 1
		g.SetBias(fromNode, toNode, rand.Float64()*2-1)
	}
}

func (g *Genome) Mutate() {
	// Equal chance for each mutation
	switch rand.Intn(4) {
	case 0:
		g.MutateAddNode()
	case 1:
		g.MutateAddConnection()
	case 2:
		g.MutateChangeWeight()
	case 3:
		g.MutateChangeBias()
	}
}
