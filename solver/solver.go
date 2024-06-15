package solver

var inf int = (1 << 30)

func invert[T comparable](m map[T][]T) map[T][]T {
	inv := make(map[T][]T)
	for k, vals := range m {
		for _, val := range vals {
			inv[val] = append(inv[val], k)
		}
	}
	return inv
}

func Solve(status map[uint32]int, children map[uint32][]uint32) {
	parents := invert(children)

	unknown_children := make(map[uint32]uint8)
	for id, ch := range children {
		if len(ch) > 255 {
			panic("too many children")
		}
		unknown_children[id] = uint8(len(ch))
	}

	q := make([]uint32, 0)
	for id, s := range status {
		if s >= 0 {
			panic("expected status to be pre-populated with only losing states")
		}
		q = append(q, parents[id]...)
	}

	for len(q) > 0 {
		id := q[0]
		q = q[1:]
		if _, ok := status[id]; ok {
			continue
		}

		lose := 1     // the longest lose path
		win := inf    // the shortest win path
		draw := false // whether we can draw
		for _, child := range children[id] {
			cstatus, ok := status[child]

			if !ok {
				draw = true
			} else if cstatus < 0 {
				win = min(win, -cstatus+1)
			} else {
				lose = max(lose, cstatus+1)
			}
		}

		if win != inf {
			status[id] = win

			// If we know the results for all children of a node,
			// we can determine the result for the node.
			for _, p := range parents[id] {
				unknown_children[p] -= 1
				if unknown_children[p] == 0 {
					q = append(q, p)
				}
			}
		} else if draw {
			panic("unknown child status, but not a win state")
		} else {
			status[id] = -lose

			// Since this is a losing state, parent states are winning.
			q = append(q, parents[id]...)
		}
	}
}
