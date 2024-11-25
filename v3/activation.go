package main

type activationFun int

const (
	Sigmoid activationFun = iota
	Tanh
	ReLU
	None // Linear
)
