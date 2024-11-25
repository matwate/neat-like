package main

import "math/rand"

func (g *Genome) MutateAddNode() {
	// Adds a node, by splitting a connection
	splitNode := rand.Intn(g.input + g.hidden)
	if splitNode >= g.input {
		splitNode += g.output
	}

	rand_conn := g.dag.nodes[splitNode].outgoing[rand.Intn(len(g.dag.nodes[splitNode].outgoing))]
	// Now we change splitNode -> rand_conn to splitNode -> newNode -> rand_conn
	g.AddNode(Hidden)
	g.AddConnection(splitNode, len(g.nodes)-1, 1, 0)
	g.AddConnection(
		len(g.nodes)-1,
		rand_conn,
		g.nodes[splitNode].connections[rand_conn].weight,
		g.nodes[splitNode].connections[rand_conn].bias,
	)
	g.RemoveConnection(splitNode, rand_conn)
}

func (g *Genome) MutateAddConnection() {
	// Pick a random not output nodes
	from := rand.Intn(g.input + g.hidden)
	if from >= g.input {
		from += g.output
	}
	// Pick a random non input nodes
	to := rand.Intn(g.hidden+g.output) + g.input
	// Check if connection already exists
	if g.dag.hasConnection(from, to) {
		return
	}
	g.AddConnection(from, to)
}

func (g *Genome) MutateChangeWeight() {
	// Pick a random connection
	from := rand.Intn(g.input + g.hidden)
	if from >= g.input {
		from += g.output
	}
	to := g.dag.nodes[from].outgoing[rand.Intn(len(g.dag.nodes[from].outgoing))]
	// Change the weight
	g.nodes[from].connections[to] = GenomeConnection{
		weight: g.nodes[from].connections[to].weight + rand.NormFloat64(),
		bias:   g.nodes[from].connections[to].bias,
	}
}

func (g *Genome) MutateChangeBias() {
	// Pick a random connection
	from := rand.Intn(g.input + g.hidden)
	if from >= g.input {
		from += g.output
	}
	to := g.dag.nodes[from].outgoing[rand.Intn(len(g.dag.nodes[from].outgoing))]
	// Change the bias
	g.nodes[from].connections[to] = GenomeConnection{
		weight: g.nodes[from].connections[to].weight,
		bias:   g.nodes[from].connections[to].bias + rand.NormFloat64(),
	}
}

func (g *Genome) Mutate() {
	// Equal chance of each mutation
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
