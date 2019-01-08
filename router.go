package astar

type Router struct {
	Src     *Point
	Dst     *Point
	current *Point
	closed  map[string]*Point
	engine  *Engine
}

func (r *Router) Pos() *Point {
	return r.current
}

func (r *Router) Next(step int, m [][]int) *Point {
	min := -1
	for _, d := range []*Direction{} {
		n := d.Move(r.current)
		tag := n.Str()
		if _, ok := r.closed[tag]; ok {
			continue
		}
		r.closed[tag] = n
		if f := r.F(n, d, step, m); min < 0 || f < min {
			min = f
			r.current = n
		}
	}
	return r.current
}

func (r *Router) baseF(n *Point, step int) int {
	return n.Distance(r.Dst) + step
}

func (r *Router) F(n *Point, d *Direction, step int, m [][]int) int {
	f := r.baseF(n, step)
	np := n.Fork()
	for np.Inline(r.Dst) {
		np = d.Move(np)
		if r.engine.isBlock(m[np.Y][np.X]) {
			f += 2*step + 2/step + 2%step
			min := -1
			for _, v := range d.Vertical() {
				tnp := v.Move(np.Fork())
				tnp = d.Move(tnp)
				cf := 0
				if r.engine.isBlock(m[tnp.Y][tnp.X]) {
					cf = 2 * step
				}
				if min < 0 || cf < min {
					min = cf
				}
			}
			f += min
			break
		}
	}
	return f
}

func (r *Router) Closed() map[string]*Point {
	return r.closed
}

func (r *Router) MoreClosed(c map[string]*Point) {
	for k, p := range c {
		r.closed[k] = p
	}
}
