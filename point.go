package astar

import (
	"fmt"
	"math"
)

type Point struct {
	X int
	Y int
}

func (p *Point) Distance(t *Point) int {
	return int(math.Abs(float64(t.X-p.X)) + math.Abs(float64(t.Y-p.Y)))
}

func (p *Point) Str() string {
	return fmt.Sprintf("%v-%v", p.X, p.Y)
}

func (p *Point) Fork() *Point {
	return &Point{
		X: p.X,
		Y: p.Y,
	}
}

func (p *Point) Inline(t *Point) bool {
	return p.X == t.X || p.Y == t.Y
}

type Direction struct {
	Point
}

func (d *Direction) Move(p *Point) *Point {
	return &Point{
		X: p.X + d.X,
		Y: p.Y + d.Y,
	}
}

func (d *Direction) Reverse() *Direction {
	return &Direction{
		Point: Point{
			X: 0 - d.X,
			Y: 0 - d.Y,
		},
	}
}

func (d *Direction) Vertical() []*Direction {
	return []*Direction{
		&Direction{
			Point: Point{
				X: 0 - d.Y,
				Y: d.X,
			},
		},
		&Direction{
			Point: Point{
				X: d.Y,
				Y: 0 - d.X,
			},
		},
	}
}

func TOP(step int) *Direction {
	return &Direction{
		Point: Point{
			Y: (0 - step),
		},
	}
}

func BOTTOM(step int) *Direction {
	return &Direction{
		Point: Point{
			Y: step,
		},
	}
}

func LEFT(step int) *Direction {
	return &Direction{
		Point: Point{
			X: (0 - step),
		},
	}
}

func RIGHT(step int) *Direction {
	return &Direction{
		Point: Point{
			X: step,
		},
	}
}
