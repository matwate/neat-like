package main

import (
	"fmt"
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
	dna     *Genome
	fitness float64
}

type Population struct {
	agents        []*Agent
	threshold     float64
	mutationCount int
	eval          TypeOfEval
	evaluate      func(*Genome) float64
}

func NewPopulation(
	threshold float64,
	x, y int,
	eval TypeOfEval,
	evaluate func(*Genome) float64,
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

func (p *Population) SelectUntilTreshold(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("Iteration:", i)
		p.Select()
		switch p.eval {
		case ABOVE:
			if p.agents[0].fitness > p.threshold {
				fmt.Println("Found a fitness above the threshold")
				fmt.Println("Best fitness:", p.agents[0].fitness)
				return
			}
		case BELOW:
			if p.agents[0].fitness < p.threshold {
				fmt.Println("Found a fitness below the threshold")
				fmt.Println("Best fitness:", p.agents[0].fitness)
				return
			}
		case CLOSEST:
			if p.agents[0].fitness-p.threshold < 0.01 {
				fmt.Println("Found a fitness close to the threshold")
				fmt.Println("Best fitness:", p.agents[0].fitness)
				return
			}
		}
	}
}

func (p *Population) Select() {
	// First, evaluate all agents concurrently
	var wg sync.WaitGroup
	for i := range p.agents {
		wg.Add(1)
		go func(i int) {
			p.agents[i].fitness = p.evaluate(p.agents[i].dna)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Sort agents by fitness
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

	// Select the best 35% of the agents
	elite := make([]*Agent, int(0.35*float64(len(p.agents))))
	copy(elite, p.agents[:len(elite)])
	// Sample the rest of the agents, weighted by fitness
	totalFitness := 0.0
	for _, agent := range p.agents {
		totalFitness += agent.fitness
	}
	for i := len(p.agents); i < len(elite); i++ {
		p.agents = append(p.agents, selectAgent(elite, totalFitness))
	}
	// Mutate the agents
	for i := len(elite); i < len(p.agents); i++ {
		for j := 0; j < p.mutationCount; j++ {
			p.agents[i].dna.Mutate()
		}
	}
	assert_equal(len(p.agents), 1000)
}

func (p *Population) Init(x, y int) {
	p.agents = make([]*Agent, 1000)
	for i := range p.agents {
		gen := NewGenome()
		gen.Init(x, y)
		p.agents[i] = &Agent{dna: &gen}
	}
}

func selectAgent(agents []*Agent, totalFitness float64) *Agent {
	r := rand.Float64() * totalFitness
	for _, agent := range agents {
		r -= agent.fitness
		if r <= 0 {
			return agent
		}
	}
	return agents[0]
}
