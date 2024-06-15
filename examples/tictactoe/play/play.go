package play

import (
	"bufio"
	"combsolve/examples/tictactoe/board"
	"fmt"
	"os"
	"strconv"
)

func statStr(s int) string {
	if s == -1 {
		return "Mate"
	} else if s < 0 {
		return fmt.Sprintf("-M%d", -(s+1)/2)
	} else if s > 0 {
		return fmt.Sprintf("M%d", s/2)
	} else {
		return "0"
	}
}

func Play(status map[uint32]int) {
	scanner := bufio.NewScanner(os.Stdin)
	state := board.StartState()

	for {
		id := state.Id()

		var stat int
		if s, ok := status[id]; !ok {
			stat = 0
		} else {
			stat = s
		}

		fmt.Printf("Current state (%s):\n%s\n", statStr(stat), state)

		if !state.Xturn() {
			if state.Won() {
				fmt.Printf("You lose.\n")
				return
			}

			for {
				fmt.Printf("Enter a move: ")
				scanner.Scan()
				move, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Printf("Invalid input %v\n", err)
					continue
				}

				state.Move(uint8(move))
				break
			}
		} else {
			if state.Won() {
				fmt.Printf("You win.\n")
				return
			}

			var score int
			var choice *board.State = nil
			for i, can := range state.Moves() {
				if !can {
					continue
				}

				newstate := state.Copy()
				newstate.Move(uint8(i))
				newid := newstate.Id()

				var sc int
				if s, ok := status[newid]; !ok {
					sc = 0
				} else {
					sc = s
				}

				if choice == nil ||
					(sc < 0 && (score >= 0 || sc > score)) ||
					(sc == 0 && score > 0) ||
					(sc > 0 && score > 0 && sc > score) {
					score = sc
					choice = newstate
				}
			}

			state = choice
		}
	}
}
