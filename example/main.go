package main

import (
	"fmt"
	"github.com/athlum/astar"
	tm "github.com/buger/goterm"
	"math/rand"
	"time"
)

type (
	fieldFunc func(int, int, int) string
)

const (
	Escaper = 2
	Ghost   = 3
)

var (
	Fields = map[int]fieldFunc{
		0: func(x, y, edge int) string {
			return " "
		},
		1: Border,
		2: func(x, y, edge int) string {
			return "P"
		},
		3: func(x, y, edge int) string {
			return "G"
		},
	}

	Map = [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	Edge = 20

	seps = []time.Duration{}
)

func Border(x, y, edge int) string {
	if x == 0 || x == edge {
		return "|"
	}
	return "-"
}

type Element struct {
	*astar.Point
	Type int
}

func (e *Element) At(x, y int) bool {
	return e.X == x && e.Y == y
}

type Engine struct {
	Map     [][]int
	Edge    int
	Escaper *Element
	Ghost   []*Element
	closec  chan struct{}
	astar   *astar.Engine
}

func newEngine() *Engine {
	return &Engine{
		Map:    Map,
		Edge:   Edge,
		closec: make(chan struct{}),
		astar: astar.New(func(i int) bool {
			return i == 1
		}),
	}
}

func (e *Engine) rand() (int, int) {
	return rand.Intn(e.Edge-2) + 1, rand.Intn(e.Edge-2) + 1
}

func (e *Engine) init() {
	x, y := e.rand()
	e.Escaper = &Element{
		Type: Escaper,
		Point: &astar.Point{
			X: x,
			Y: y,
		},
	}

	e.Ghost = make([]*Element, 1)
	for i := 0; i < 1; i += 1 {
		x, y := e.rand()
		e.Ghost[i] = &Element{
			Type: Ghost,
			Point: &astar.Point{
				X: x,
				Y: y,
			},
		}
	}
}

func (e *Engine) state() [][]int {
	s := make([][]int, len(e.Map))
	for y, r := range e.Map {
		s[y] = make([]int, len(r))
		for x, f := range r {
			if e.Escaper.At(x, y) {
				s[y][x] = e.Escaper.Type
			} else {
				located := false
				for _, g := range e.Ghost {
					if g.At(x, y) {
						located = true
						s[y][x] = g.Type
						break
					}
				}
				if !located {
					s[y][x] = f
				}
			}
		}
	}
	return s
}

func (e *Engine) ToMap(m [][]int) string {
	if m == nil {
		m = e.Map
	}
	s := ""
	for y, r := range m {
		for x, f := range r {
			s = fmt.Sprintf("%v %v", s, Fields[f](x, y, e.Edge))
		}
		s += "\n"
	}
	return s
}

func main() {
	e := newEngine()
	e.init()
	tm.Clear()
	for {
		select {
		case <-e.closec:
			var d time.Duration
			for _, s := range seps {
				d += s
			}
			fmt.Println("Game Over", time.Duration(int(d)/len(seps)))
			return
		case <-time.After(time.Millisecond * 100):
			e.flush()
			tm.MoveCursor(1, 1)
			tm.Println(e.ToMap(e.state()))
			tm.Flush()
		}
	}
}

func (e *Engine) flush() {
	e.Escaper.Escape(e)
	for _, g := range e.Ghost {
		g.Catch(e)
		if g.At(e.Escaper.X, e.Escaper.Y) {
			close(e.closec)
		}
	}
}

func (e *Element) Escape(s *Engine) {
}

func (e *Element) Catch(s *Engine) {
	m := s.state()
	now := time.Now()
	n := s.astar.Router(e.Point, s.Escaper.Point).Path(1, m)
	seps = append(seps, time.Now().Sub(now))
	if n != nil {
		e.Point = n[0]
	}
}
