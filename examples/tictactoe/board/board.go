package board

import "fmt"

type State struct {
	x     [3]uint8
	o     [3]uint8
	xturn bool
}

func StartState() *State {
	return &State{
		x:     [3]uint8{255, 255, 255},
		o:     [3]uint8{255, 255, 255},
		xturn: true,
	}
}

func (s *State) Xturn() bool {
	return s.xturn
}

func (s *State) Id() uint32 {
	id := uint32(0)

	for _, p := range s.x {
		id = 10*id + uint32(min(p, 9))
	}
	for _, p := range s.o {
		id = 10*id + uint32(min(p, 9))
	}

	id *= 2
	if s.xturn {
		id += 1
	}

	return id
}

func (s *State) Copy() *State {
	return &State{
		x:     s.x, // This makes a shallow copy
		o:     s.o, // This makes a shallow copy
		xturn: s.xturn,
	}
}

func (s *State) Moves() [9]bool {
	can := [9]bool{true, true, true, true, true, true, true, true, true}
	for _, p := range s.x {
		if p != 255 {
			can[p] = false
		}
	}
	for _, p := range s.o {
		if p != 255 {
			can[p] = false
		}
	}
	return can
}

func (s *State) Move(pos uint8) {
	var a *[3]uint8
	if s.xturn {
		a = &s.x
	} else {
		a = &s.o
	}

	*a = [3]uint8{a[1], (*a)[2], pos}
	s.xturn = !s.xturn
}

func (s *State) Won() bool {
	var a *[3]uint8
	if s.xturn {
		a = &s.o
	} else {
		a = &s.x
	}

	if a[0] == 255 {
		return false
	}

	return ((a[0]/3 == a[1]/3 && a[1]/3 == a[2]/3) || // all same row
		(a[0]%3 == a[1]%3 && a[1]%3 == a[2]%3) || // all same col
		(a[0]/3 == a[0]%3 && a[1]/3 == a[1]%3 && a[2]/3 == a[2]%3) || // positive diagonal
		(a[0]/3 == 2-a[0]%3 && a[1]/3 == 2-a[1]%3 && a[2]/3 == 2-a[2]%3)) // negative diagonal
}

func (s *State) String() string {
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"

	board := [][]string{{"0 ", "1 ", "2 "}, {"3 ", "4 ", "5 "}, {"6 ", "7 ", "8 "}}
	ts := []string{green + "x", red + "o"}

	for t, arr := range [][3]uint8{s.x, s.o} {
		for i, p := range arr {
			if p == 255 {
				continue
			}
			board[p/3][p%3] = ts[t]
			if (s.xturn == (t == 0)) && i == 0 {
				board[p/3][p%3] += "*"
			} else {
				board[p/3][p%3] += " "
			}
			board[p/3][p%3] += reset
		}
	}

	// str := fmt.Sprintf("id: %d\n", s.Id())
	str := ""
	str += "----------------\n"
	for i := range 3 {
		str += "|"
		for j := range 3 {
			str += " " + board[i][j] + " |"
		}
		str += "\n----------------\n"
	}

	return str
}

func (s *State) Verbose() string {
	return fmt.Sprintf("x=%v, o=%v, xturn=%v", s.x, s.o, s.xturn)
}
