package main

import "math/rand"

func (g *Genome) MutateAddNode() {
	// Adds a node, by splitting a connection
	splitNode := rand.Intn(g.input + g.hidden)
	if splitNode >= g.input {
		splitNode += g.output
	}
	if len(g.dag.nodes[splitNode].outgoing) == 0 {
		return
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
	if len(g.dag.nodes[from].outgoing) == 0 {
		return
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
	if len(g.dag.nodes[from].outgoing) == 0 {
		return
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
	if len(g.dag.nodes[from].outgoing) == 0 {
		return
	}
	to := g.dag.nodes[from].outgoing[rand.Intn(len(g.dag.nodes[from].outgoing))]
	// Change the bias
	g.nodes[from].connections[to] = GenomeConnection{
		weight: g.nodes[from].connections[to].weight,
		bias:   g.nodes[from].connections[to].bias + rand.NormFloat64(),
	}
}

func (g *Genome) Mutate() {
	switch rand.Intn(4) {
	case 0:
		switch rand.Intn(2) {
		case 0:
			g.MutateChangeWeight()
		case 1:
			g.MutateChangeBias()
		}
	default:
		// Do nothing
	}
	// Random number between 0 and 1
	random := rand.Float64()
	if random < 0.05 {
		g.MutateAddNode()
	}
	if random < 0.8 {
		g.MutateAddConnection()
	}
}
