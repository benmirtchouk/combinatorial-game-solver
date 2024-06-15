package main

import (
	"combsolve/examples/tictactoe/board"
	"combsolve/examples/tictactoe/play"
	"combsolve/solver"
	"fmt"
)

func main() {
	start := board.StartState()
	startid := start.Id()

	states := map[uint32]*board.State{startid: start}
	children := make(map[uint32][]uint32)
	status := make(map[uint32]int)
	q := []uint32{startid}

	for len(q) > 0 {
		id := q[0]
		q = q[1:]

		state := states[id]
		if state.Won() {
			status[id] = -1
		}

		for i, can := range state.Moves() {
			if !can {
				continue
			}

			newstate := state.Copy()
			newstate.Move(uint8(i))
			newstate = newstate.MinimizeId()
			newid := newstate.Id()

			children[id] = append(children[id], newid)

			if _, ok := states[newid]; !ok {
				states[newid] = newstate
				q = append(q, newid)
			}
		}
	}

	fmt.Printf("Found a total of %d states\n", len(states))

	solver.Solve(status, children)
	wins := 0
	loses := 0
	for _, s := range status {
		if s > 0 {
			wins += 1
		} else if s < 0 {
			loses += 1
		} else {
			panic("draws arent stored in status")
		}
	}

	fmt.Printf("Found %d wins, %d loses, %d draws\n", wins, loses, len(states)-wins-loses)

	play.Play(status)
}
