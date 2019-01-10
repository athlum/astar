package astar

import (
	"fmt"
	"sync"
)

type Router struct {
	Src      *Point
	Dst      *Point
	engine   *Engine
	routeMap *sync.Map
}

func (r *Router) Path(step int, m [][]int) []*Point {
	ps, ok := r.process([]*Point{r.Src}, r.Src, step, m, nil)
	if !ok {
		return nil
	}
	return r.improve(append([]*Point{r.Src}, ps...), m)[1:]
}

func (r *Router) improve(ps []*Point, m [][]int) []*Point {
	cps := make([]*Point, len(ps))
	copy(cps, ps)
	cps = r.improveCircle(cps, m)
	cps = r.improveUPath(cps, m)
	return cps
}

type (
	uPath struct {
		in        []*Point
		inDir     *Direction
		corner    []*Point
		cornerDir *Direction
		out       []*Point
		outDir    *Direction
		state     int
		cache     []*Point
	}
)

const (
	uIn = iota
	uCorner
	uOut
)

func (p *uPath) reset() {
	p.in = []*Point{}
	p.inDir = nil
	p.corner = []*Point{}
	p.cornerDir = nil
	p.out = []*Point{}
	p.outDir = nil
	p.state = uIn
	p.cache = []*Point{}
}

func (r *Router) improveUPath(ps []*Point, m [][]int) []*Point {
	var (
		i = 0
		u = &uPath{
			in:     []*Point{},
			corner: []*Point{},
			out:    []*Point{},
			state:  uIn,
		}
		posIn  = -1
		posOut = -1
		lastOut = false
	)
	for {
		n := ps[i]
		next := false
		switch u.state {
		case uIn:
			if len(u.in) == 0 {
				u.in = append(u.in, n)
				next = true
			} else {
				l := u.in[len(u.in)-1]
				d := l.Dir(n)
				if u.inDir == nil || d.Equals(u.inDir) {
					u.in = append(u.in, n)
					if u.inDir == nil {
						u.inDir = d
					}
					next = true
				} else {
					u.cornerDir = d
					u.state = uCorner
				}
			}
		case uCorner:
			if len(u.corner) == 0 {
				u.corner = append(u.corner, n)
				next = true
			} else {
				l := u.corner[len(u.corner)-1]
				d := l.Dir(n)
				if d.Equals(u.cornerDir) {
					u.corner = append(u.corner, n)
					next = true
				} else {
					if d.IsReverse(u.inDir) {
						u.state = uOut
						u.outDir = d
					} else {
						u.cache = append(u.cache, u.in[:len(u.in)-1]...)
						u.in = append([]*Point{u.in[len(u.in)-1]}, u.corner...)
						u.corner = []*Point{}
						u.inDir = u.cornerDir
						u.cornerDir = d
					}
				}
			}
		case uOut:
			if posOut == -1 && posIn == -1 {
				if len(u.out) == 0 {
					u.out = append(u.out, n)
					next = true
				} else {
					l := u.out[len(u.out)-1]
					d := l.Dir(n)
					if d.Equals(u.outDir) {
						u.out = append(u.out, n)
						next = true
					} else {
						posOut = len(u.out) - 1
						posIn = len(u.in) - len(u.out) - 1
						if posIn < 0 {
							posOut += posIn
							posIn = 0
						}
					}
				}

				if i == len(ps)-1 && next {
					posOut = len(u.out) - 1
					posIn = len(u.in) - len(u.out) - 1
					if posIn < 0 {
						posOut += posIn
						posIn = 0
					}
					next = false
					lastOut = true
				}
			} else if posOut == -1 { // unimproveable
				if lastOut {
					next = true
					break
				}
				u.cache = append(u.cache, u.in...)
				u.cache = append(u.cache, u.corner[:len(u.corner)-1]...)
				u.in = append([]*Point{u.corner[len(u.corner)-1]}, u.out...)
				u.inDir = u.outDir
				u.corner = []*Point{}
				u.out = []*Point{}
				u.state = uIn
				posIn = -1
				posOut = -1
			} else {
				s := u.in[posIn]
				n := s.Fork()
				d := u.out[posOut]
				dd := n.Dir(d)
				corner := []*Point{}
				access := true
				for {
					n = dd.Move(n)
					if r.blocked(n, m) {
						posIn += 1
						posOut -= 1
						access = false
						break
					}
					if n.At(d.X, d.Y) {
						break
					} else {
						corner = append(corner, n)
					}
				}
				if access {
					ips := u.in[:posIn+1]
					if len(u.cache) > 0 {
						ips = append(u.cache, ips...)
					}
					ips = append(ips, corner...)
					ips = append(ips, u.out[posOut:]...)
					if !(i == len(ps)-1 && !next) {
						ips = append(ips, ps[i:]...)
					}
					ps = ips
					i = 0
					posIn = -1
					posOut = -1
					u.reset()
				}
			}
		}
		if next {
			i += 1
			if i == len(ps) {
				break
			}
		}
	}
	return ps
}

func (r *Router) improveCircle(ps []*Point, m [][]int) []*Point {
	i := 0
	pi := make(map[string]int)
	for {
		p := ps[i]
		if v, ok := pi[p.Str()]; ok {
			ps = append(ps[:v], ps[i:len(ps)]...)
			i = 0
			pi = make(map[string]int)
			continue
		} else {
			pi[p.Str()] = i
		}
		i += 1
		if i == len(ps) {
			break
		}
	}
	return ps
}

type router struct {
	path []*Point
	ok   bool
}

func (r *Router) process(prev []*Point, n *Point, step int, m [][]int, sd *Direction) ([]*Point, bool) {
	var ps []*Point
	for !n.At(r.Dst.X, r.Dst.Y) {
		nps, ok := r.next(append(prev, ps...), n, step, m, sd)
		if !ok {
			return nil, ok
		}
		ps = append(ps, nps...)
		n = nps[len(nps)-1]
	}
	return ps, true
}

func (r *Router) next(prev []*Point, n *Point, step int, m [][]int, sd *Direction) ([]*Point, bool) {
	d := n.Dir(r.Dst).Mul(step)
	np := d.Move(n)
	if r.blocked(np, m) || (sd != nil && sd.IsReverse(d)) {
		rp, ok := r.routing(prev, n, step, d, m)
		if !ok {
			return nil, ok
		}
		return rp, true
	}
	return []*Point{np}, true
}

func (r *Router) routing(prev []*Point, n *Point, step int, d *Direction, m [][]int) ([]*Point, bool) {
	var (
		rp    *router
		valid = 0
		pc    = make(chan router, 2)
		min   = -1
	)

	for _, vd := range d.Vertical() {
		rd := &routerDirection{
			p: n,
			d: vd,
		}
		if _, loaded := r.routeMap.LoadOrStore(rd.Str(), rd); loaded {
			continue
		}
		valid += 1
		go r.route(prev, pc, n, step, d, vd.Mul(step), m)
	}
	for i := 0; i < valid; i += 1 {
		p := <-pc
		if !p.ok {
			continue
		}
		nps := r.improve(append(prev, p.path...), m)[1:]
		if min == -1 || len(nps) < min {
			min = len(nps)
			rp = &p
		}
	}
	if rp == nil {
		return nil, false
	}
	return rp.path, rp.ok
}

type routerDirection struct {
	p *Point
	d *Direction
}

func (r *routerDirection) Str() string {
	return fmt.Sprintf("%v.%v-%v.%v", r.p.X, r.p.Y, r.d.X, r.d.Y)
}

func (r *Router) route(prev []*Point, pc chan router, n *Point, step int, d, vd *Direction, m [][]int) {
	var (
		rr = router{
			path: []*Point{},
		}
		cn = n.Fork()
		np *Point
	)
	for {
		np = vd.Move(cn)
		if r.blocked(np, m) {
			pp, ok := r.routing(append(prev, rr.path...), cn, step, vd, m)
			rr.path = append(rr.path, pp...)
			rr.ok = ok
			break
		}
		rr.path = append(rr.path, np)
		cn = np
		bn := d.Move(cn)
		if r.Dst.At(bn.X, bn.Y) {
			rr.path = append(rr.path, bn)
			rr.ok = true
			break
		}
		if !r.blocked(bn, m) {
			bbn := d.Move(bn)
			if r.Dst.At(bbn.X, bbn.Y) {
				rr.path = append(rr.path, bn, bbn)
				rr.ok = true
				break
			}
			if r.blocked(bbn, m) {
				rr.path = append(rr.path, bn)
				pp, ok := r.routing(append(prev, rr.path...), bn, step, d, m)
				rr.path = append(rr.path, pp...)
				rr.ok = ok
				break
			} else {
				rr.path = append(rr.path, bn, bbn)
				pp, ok := r.process(append(prev, rr.path...), bbn, step, m, d)
				rr.path = append(rr.path, pp...)
				rr.ok = ok
				break
			}
		}
	}
	pc <- rr
}

func (r *Router) blocked(n *Point, m [][]int) bool {
	return n.OutOfArea(m) || r.engine.isBlock(m[n.Y][n.X])
}
