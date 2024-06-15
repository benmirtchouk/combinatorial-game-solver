package board

var (
	figL    int   = 0
	figB    int   = 1
	figR    int   = 2
	figP    int   = 3
	figH    int   = 4
	pos     uint8 = 0xf
	isH     uint8 = 0x20
	isW     uint8 = 0x10
	off     uint8 = 12
	offsets       = [][]int{
		{-4, -3, -2, -1, 1, 2, 3, 4}, // L
		{-4, -2, 2, 4},               // B
		{-3, -1, 1, 3},               // R
		{-3},                         // P
		{-4, -3, -2, -1, 1, 3},       // H
	}
)

type State struct {
	pieces [8]uint8
}

func (s *State) Copy() *State {
	return &State{
		pieces: s.pieces, // This makes a shallow copy
	}
}

func (s *State) flip() {
	for f := range s.pieces {
		p := 11 - (s.pieces[f] & pos)
		s.pieces[f] = (s.pieces[f] & ^pos | p) ^ isW
	}
}

func (s *State) TryMove(f int, p uint8) *State {
	figType := f / 2
	ns := s.Copy()

	// check for captures
	for f2, cp := range s.pieces {
		if (cp & pos) == p {
			if (cp & isW) != 0 {
				ns.pieces[f2] = off
			} else {
				// can't capture own piece
				return nil
			}
		}
	}

	// check for pawn promotion
	hen := s.pieces[f] & isH
	if figType == figP && p/4 == 0 && (s.pieces[f]&pos) != off {
		hen = isH
	}

	ns.pieces[f] = p | hen

	ns.flip()
	return ns
}

func (s *State) Moves() []*State {
	ret := make([]*State, 0)

	// try move a piece
	for f, offs := range offsets {
		if (f == figP && (s.pieces[figP]&isH != 0)) ||
			(f == figH && (s.pieces[figP]&isH == 0)) {
			continue
		}

		var (
			cur  uint8
			fidx int
		)
		if f == figH {
			fidx = figP
			cur = s.pieces[figP]
		} else {
			fidx = f
			cur = s.pieces[f]
		}

		for _, o := range offs {
			if (o < 0 && cur < uint8(-o)) || cur >= uint8(12-o) {
				continue
			}
			if m := s.TryMove(fidx, cur+uint8(o)); m != nil {
				ret = append(ret, m)
			}
		}
	}

	full := [12]bool{} // all false by default
	for _, p := range s.pieces {
		if p != off {
			full[p] = true
		}
	}

	// try place a piece
	// TODO

	return ret
}
