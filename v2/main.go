package main

import (
	"fmt"
	"math"
)

func main() {
	s := NewPopulation(9, 1, 1, ABOVE, fitnessF, 10)
	s.Select()
	for i := 0; i < len(s.agents); i++ {
		fmt.Printf("Agent %v, fitness %v\n", i, s.agents[i].fitness)
	}
	var bestAgent Agent
	s.SelectUntilTreshold(2136871236126798)
	bestAgent = s.agents[0]
	for i := 0.0; i < 1; i += 0.1 {
		y := bestAgent.dna.Forward(i)[0]
		fmt.Printf(
			"Sin(x) = %v, fitted %v, passed %v\n",
			math.Sin(i),
			y,
			math.Abs(y-math.Sin(i)) < 0.01,
		)
	}
}

func fitnessF(g Genome) float64 {
	// We will sample 10 points between 0 and 1, and if the result is close enough to sin(x) we'll +1 to the fitnessa
	var score float64
	for i := 0.0; i < 1; i += 0.1 {
		x := i
		y := g.Forward(x)[0]
		if math.Abs(y-math.Sin(x)) < 0.01 {
			score++
		}
	}
	return score
}
