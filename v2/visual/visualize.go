package visual

import (
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This will a specific representation of the nodes and connections in a DAG (i dont wanna do weights and biases yet)

// The representation will be basically a list of strings.
// Each string will represent a node, the connections are in the string, separated by spaces

// So that way:
/*
["1 2", "2", ""]

Would represent a DAG with 3 nodes, the first node has two connections, the second node has one connection, and the third node has no connections, i will NOT BE PERFORMING CHECKS.
*/

type DAGVisual []string

type Visualizer struct {
	pos map[int]rl.Vector2
	DAG DAGVisual
}

func (v *Visualizer) initPos() {
	for i := 0; i < len(v.DAG); i++ {
		v.pos[i] = rl.NewVector2(
			float32(rl.GetRandomValue(0, 800)),
			float32(rl.GetRandomValue(0, 600)),
		)
	}
}

func (v *Visualizer) drawDAG() {
	for i := 0; i < len(v.DAG); i++ {
		rl.DrawCircleV(v.pos[i], 10, rl.Red)
	}

	for i, conns := range v.DAG {
		conns := strings.Split(conns, " ")
		for _, conn := range conns {
			connInt, _ := strconv.Atoi(conn)
			rl.DrawLineEx(v.pos[i], v.pos[connInt], 2, rl.White)
		}
	}
}

func Vis(DAG DAGVisual) Visualizer {
	return Visualizer{
		pos: make(map[int]rl.Vector2),
		DAG: DAG,
	}
}

func Visualize(DAG DAGVisual) {
	visualizer := Vis(DAG)
	rl.InitWindow(800, 600, "DAG Visualizer")
	visualizer.initPos()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("DAG Visualizer", 10, 10, 20, rl.RayWhite)
		// Draw DAG
		visualizer.drawDAG()
		rl.EndDrawing()
	}
}
