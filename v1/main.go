package main

import (
	"fmt"
	"math"
)

func fitness(ag Agent, _ int) float64 {
	// We will sample 10 points between 0 to 1, fit the function and add 1 to the fitness if the value is within 0.1 of sin(x)
	fit := 0
	for i := 0.0; i < 1; i += 0.1 {
		pred := ag.neatlike.Fit([]float64{i})
		fmt.Println(pred)
		if math.Abs(pred[0]-math.Sin(i)) < 0.001 {
			fit += 1
		}
	}
	return float64(fit)
}

func main() {
	template := Agent{NewNeatlike(1, 1), 0, false}
	s := NewSim(100, 9, fitness, template, Above)
	agent := s.RunUntilThreshold(10)
	print_matrix(agent.neatlike.dag.adjacencyMatrix)
	fmt.Println(
		agent.neatlike.inputs,
		agent.neatlike.outputs,
		agent.wasMutated,
		agent.neatlike.dag.size,
		agent.neatlike.dag.nodeValues,
	)
	for i := 0.0; i < 1; i += 0.1 {
		fmt.Printf(
			"sin(%f) = %f, pred = %f\n",
			i,
			math.Sin(i),
			agent.neatlike.Fit([]float64{i})[0],
		)
	}
}

func print_matrix[T any](m [][]T) {
	for i := range m {
		for j := range m[i] {
			fmt.Print(m[i][j], " ")
		}
		fmt.Println()
	}
}
