package main

type Activation int

const (
	Sigmoid Activation = iota
	Tanh
	ReLU
	None // Linear
)
