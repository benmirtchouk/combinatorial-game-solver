package board

import "fmt"

type State struct {
	x [3]uint8
	o [3]uint8
}

func StartState() *State {
	return &State{
		x: [3]uint8{255, 255, 255},
		o: [3]uint8{255, 255, 255},
	}
}

func (s *State) Id() uint32 {
	id := uint32(0)

	for _, p := range s.x {
		id = 10*id + uint32(min(p, 9))
	}
	for _, p := range s.o {
		id = 10*id + uint32(min(p, 9))
	}

	return id
}

func (s *State) Copy() *State {
	return &State{
		x: s.x, // This makes a shallow copy
		o: s.o, // This makes a shallow copy
	}
}

var r90map [9]uint8 = [9]uint8{6, 3, 0, 7, 4, 1, 8, 5, 2}
var flipmap [9]uint8 = [9]uint8{2, 1, 0, 5, 4, 3, 8, 7, 6}

func applyMap(a [3]uint8, m [9]uint8) [3]uint8 {
	na := [3]uint8{}
	for i := range 3 {
		if a[i] == 255 {
			na[i] = 255
		} else {
			na[i] = m[a[i]]
		}
	}
	return na
}

func (s *State) Flip() *State {
	return &State{
		x: applyMap(s.x, flipmap),
		o: applyMap(s.o, flipmap),
	}
}

func (s *State) Rotate90() *State {
	return &State{
		x: applyMap(s.x, r90map),
		o: applyMap(s.o, r90map),
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

func (s *State) MinimizeId() *State {
	var best *State = nil
	mn := uint32(1 << 30)

	for f := range 2 {
		cur := s
		if f != 0 {
			cur = cur.Flip()
		}
		for r := range 4 {
			if r != 0 {
				cur = cur.Rotate90()
			}
			if id := cur.Id(); id < mn {
				best = cur
				mn = id
			}
		}
	}

	return best
}

func (s *State) Move(pos uint8) {
	s.x, s.o = s.o, [3]uint8{s.x[1], s.x[2], pos}
}

func (s *State) Won() bool {
	if s.o[0] == 255 {
		return false
	}

	return ((s.o[0]/3 == s.o[1]/3 && s.o[1]/3 == s.o[2]/3) || // all same row
		(s.o[0]%3 == s.o[1]%3 && s.o[1]%3 == s.o[2]%3) || // all same col
		(s.o[0]/3 == s.o[0]%3 && s.o[1]/3 == s.o[1]%3 && s.o[2]/3 == s.o[2]%3) || // positive diagonal
		(s.o[0]/3 == 2-s.o[0]%3 && s.o[1]/3 == 2-s.o[1]%3 && s.o[2]/3 == 2-s.o[2]%3)) // negative diagonal
}

func (s *State) String(oturn bool) string {
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"

	board := [][]string{{"0 ", "1 ", "2 "}, {"3 ", "4 ", "5 "}, {"6 ", "7 ", "8 "}}
	ts := []string{green + "x", red + "o"}
	if oturn {
		ts[0], ts[1] = ts[1], ts[0]
	}

	for t, arr := range [][3]uint8{s.x, s.o} {
		for i, p := range arr {
			if p == 255 {
				continue
			}
			board[p/3][p%3] = ts[t]
			if t == 0 && i == 0 {
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
	return fmt.Sprintf("x=%v, o=%v", s.x, s.o)
}
