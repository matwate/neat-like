package main

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Error("NewGraph() returned nil")
	}
}

func TestAddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode()
	if len(g.adjacencyList) != 1 {
		t.Error("AddNode() failed")
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	g.AddNode()
	g.AddNode()
	g.AddEdge(0, 1)
	if len(g.adjacencyList[0].adjacents) != 1 {
		t.Error("AddEdge() failed")
	}
}
