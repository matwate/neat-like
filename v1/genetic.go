package main

import (
	"math"
	"math/rand"
	"sort"
	"sync"
)

type Agent struct {
	neatlike   Neatlike
	fitness    float64
	wasMutated bool
} // This is for sorting agents by fitness

type Population []Agent

type Condition int

const (
	Above Condition = iota
	Below
	AsCloseAs
)

// Now methods for running the simulation
// It will be like this:
/*
   First N agents will be created, with random weights and random biases.
  Then they will run through the fitness function, and the fitness will be calculated.
  Then they will be sorted by fitness.
  30% of the agents will pass to the next generation
  The other 70% will be randomly sampled from the whole population weighted by fitness
  To those 70% a mutation will be applied.
  Things will stop once the fitness is above a certain threshold, set by the user.
*/

type Simulation struct {
	population       Population
	fitnessThreshold float64
	condition        Condition
	fitnessFunction  func(Agent, int) float64 // This is so fitness functions can do whatever they want
	templateAgent    Agent
	once             sync.Once
}

func NewSim(
	population int,
	fitnessThreshold float64,
	fitnessFunction func(Agent, int) float64,
	templateAgent Agent,
	condition Condition,
) *Simulation {
	s := new(Simulation)
	s.population = make([]Agent, population)
	s.fitnessThreshold = fitnessThreshold
	s.fitnessFunction = fitnessFunction
	s.templateAgent = templateAgent
	s.condition = condition
	s.once = sync.Once{}
	return s
}

func (s *Simulation) Run() []Agent {
	// Spawn the first generation
	s.once.Do(func() {
		for i := 0; i < len(s.population); i++ {
			s.population[i].neatlike = NewNeatlike(
				s.templateAgent.neatlike.inputs,
				s.templateAgent.neatlike.outputs,
			)
		}
	})
	// Concurrently calculate the fitness of each agent
	var wg sync.WaitGroup
	for i := range s.population {
		wg.Add(1)
		go func(agent *Agent, index int) {
			defer wg.Done()
			fitness := s.fitnessFunction(*agent, index)
			agent.fitness = fitness
		}(&s.population[i], i)
	}
	wg.Wait()

	switch s.condition {
	case Above:
		sort.Slice(s.population, func(i, j int) bool {
			return s.population[i].fitness > s.population[j].fitness
		})
	case Below:
		sort.Slice(s.population, func(i, j int) bool {
			return s.population[i].fitness < s.population[j].fitness
		})
	case AsCloseAs:
		sort.Slice(s.population, func(i, j int) bool {
			diffI := math.Abs(s.population[i].fitness - s.fitnessThreshold)
			diffJ := math.Abs(s.population[j].fitness - s.fitnessThreshold)
			return diffI < diffJ
		})

	}

	// Now we pick the best 30% of the agents
	thresh := int(float64(len(s.population)) * 0.3)

	// Calculate total fitness
	totalFitness := 0.0
	for _, agent := range s.population {
		totalFitness += agent.fitness
	}

	// Initialize newPopulation with top 30%
	newPopulation := make([]Agent, 0, len(s.population))
	for i := 0; i < thresh; i++ {
		newPopulation = append(newPopulation, s.population[i])
	}

	// Sample the remaining 70% of the agents
	for i := thresh; i < len(s.population); i++ {

		selectedAgent := selectAgent(s.population, totalFitness)

		// Create a copy of the selected agent
		copiedAgent := selectedAgent

		// Mutate the copied agent
		copiedAgent.neatlike.Mutate()

		copiedAgent.wasMutated = true
		newPopulation = append(newPopulation, copiedAgent)

	}

	// Replace the old population with the new one
	s.population = newPopulation

	return newPopulation
}

// Roulette wheel selection function
func selectAgent(population []Agent, totalFitness float64) Agent {
	r := rand.Float64() * totalFitness
	cumulative := 0.0

	for _, agent := range population {

		cumulative += agent.fitness
		if cumulative >= r {
			return agent
		}
	}
	// In case of rounding errors, return the last agent
	return population[len(population)-1]
}

// This one will run until the fitness is above a certain threshold
func (s *Simulation) RunUntilThreshold(max int) Agent {
	atts := 0
	best := s.population[0].fitness
Outer:
	for {
		if atts > max {
			break
		}
		atts++
		s.Run()
		switch s.condition {
		case Above:
			if s.population[0].fitness > s.fitnessThreshold {
				if s.population[0].fitness > best {
					best = s.population[0].fitness
				}
				break Outer
			}
		case Below:
			if s.population[0].fitness < s.fitnessThreshold {
				if s.population[0].fitness < best {
					best = s.population[0].fitness
				}
				break Outer
			}
		case AsCloseAs:
			if math.Abs(s.population[0].fitness-s.fitnessThreshold) < 0.0001 {
				if math.Abs(
					s.population[0].fitness-s.fitnessThreshold,
				) < math.Abs(
					best-s.fitnessThreshold,
				) {
					best = s.population[0].fitness

					break Outer
				}
			}
		}
	}
	return s.population[0]
}

func (a *Agent) Copy() Agent {
	return Agent{
		neatlike:   a.neatlike.Copy(),
		fitness:    a.fitness,
		wasMutated: a.wasMutated,
	}
}
