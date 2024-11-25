package main

import (
	"math"
	"math/rand"
	"sort"
	"sync"
)

type TypeOfEval int

const (
	ABOVE TypeOfEval = iota
	BELOW
	CLOSEST
)

type Agent struct {
	dna     Genome
	fitness float64
}

type Population struct {
	agents        []Agent
	threshold     float64
	mutationCount int
	eval          TypeOfEval
	evaluate      func(Genome) float64
}

func NewPopulation(
	threshold float64,
	x, y int,
	eval TypeOfEval,
	evaluate func(Genome) float64,
	mutationCount int,
) *Population {
	var p Population
	p.threshold = threshold
	p.eval = eval
	p.evaluate = evaluate
	p.mutationCount = mutationCount
	p.Init(x, y)
	return &p
}

func (p *Population) Select() {
	// Evaluate the fitness of each agents
	var wg sync.WaitGroup
	for i := range p.agents {
		wg.Add(1)
		go func(i int) {
			p.agents[i].fitness = p.evaluate(p.agents[i].dna)
			wg.Done()
		}(i)
	}
	wg.Wait()
	// Sort the population by fitness
	switch p.eval {
	case ABOVE:
		sort.Slice(p.agents, func(i, j int) bool {
			return p.agents[i].fitness > p.agents[j].fitness
		})
	case BELOW:
		sort.Slice(p.agents, func(i, j int) bool {
			return p.agents[i].fitness < p.agents[j].fitness
		})
	case CLOSEST:
		sort.Slice(p.agents, func(i, j int) bool {
			return math.Abs(
				p.agents[i].fitness-p.threshold,
			) < math.Abs(
				p.agents[j].fitness-p.threshold,
			)
		})
	}
	// Pick the top 30% of the population
	newPopulation := make([]Agent, len(p.agents))
	for i := 0; i < len(p.agents)/3; i++ {
		newPopulation[i] = p.agents[i]
	}
	// Sample the rest of the population weighted by fitness
	totalFitness := 0.0
	for _, agent := range p.agents {
		totalFitness += agent.fitness
	}
	for i := len(p.agents) / 3; i < len(p.agents); i++ {
		selectedAgent := selectAgent(p.agents, totalFitness)
		selectedAgent.dna.Mutate()
		newPopulation[i] = selectedAgent
	}
	copy(p.agents, newPopulation)
}

func (p *Population) SelectUntilTreshold(escape int) {
	for i := 0; i < escape; i++ {
		p.Select()
		switch p.eval {
		case ABOVE:
			if p.agents[0].fitness > p.threshold {
				return
			}
		case BELOW:
			if p.agents[0].fitness < p.threshold {
				return
			}
		case CLOSEST:
			if math.Abs(p.agents[0].fitness-p.threshold) < 0.01 {
				return
			}

		}
	}
}

func (p *Population) Init(x, y int) {
	p.agents = make([]Agent, 100)
	for i := range p.agents {
		p.agents[i].dna = NewGenome()
		p.agents[i].dna.Init(x, y)
	}
}

func selectAgent(agents []Agent, totalFitness float64) Agent {
	r := rand.Float64() * totalFitness
	for _, agent := range agents {
		r -= agent.fitness
		if r <= 0 {
			return agent
		}
	}
	return agents[0]
}
