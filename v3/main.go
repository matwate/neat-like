package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	// New simulation
	s := NewPopulation(
		4.5,
		1, 1,
		ABOVE,
		fit,
		5,
	)
	s.Select()
	for i := 0; i < len(s.agents); i++ {
		fmt.Printf("Agent %v, fitness %v\n", i, s.agents[i].fitness)
	}
}

func fit(g Genome) float64 {
	// We will give it 5 random values and if the net's value is close enough to sin(x) we'll give a +1 to the fitness
	var fitness float64
	for i := 0; i < 5; i++ {
		x := rand.Float64() * 2 * math.Pi
		if math.Abs(g.Forward(x)[0]-math.Sin(x)) < 0.1 {
			fitness++
		}
	}
	return fitness
}
