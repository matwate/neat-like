package main

import (
	"fmt"
	"math"
)

func main() {
	// New simulation
	s := NewPopulation(
		0,
		1, 1,
		CLOSEST,
		fit,
		5,
	)

	s.SelectUntilTreshold(20)
	fmt.Println("Best fitness:", s.agents[0].fitness)
}

func fit(g *Genome) float64 {
	// We will take 100 evenly spaced samples between 0 and 2π and calculate the sum of the squares of the differences between the sine of the sample and the sample itself.
	// The fitness of the genome will be the negative of this sum.
	sum := 0.0
	for i := 0; i < 100; i++ {
		x := 2 * math.Pi * float64(i) / 100
		sum += math.Pow(math.Sin(x)-g.Forward(x)[0], 2)
	}
	return sum
}
