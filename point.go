package astar

import (
	"fmt"
	"math"
)

var (
	TOP = &Direction{
		Point: Point{
			X: 0,
			Y: -1,
		},
	}

	DOWN = &Direction{
		Point: Point{
			X: 0,
			Y: 1,
		},
	}

	LEFT = &Direction{
		Point: Point{
			X: -1,
			Y: 0,
		},
	}

	RIGHT = &Direction{
		Point: Point{
			X: 1,
			Y: 0,
		},
	}
)

type Point struct {
	X int
	Y int
}

func (p *Point) At(x, y int) bool {
	return p.X == x && p.Y == y
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

func (p *Point) OutOfArea(m [][]int) bool {
	return p.Y < 0 || p.X < 0 || p.Y >= len(m) || p.X >= len(m[p.Y])
}

func (p *Point) Dir(t *Point) *Direction {
	if p.At(t.X, t.Y) {
		return nil
	}
	sx := math.Abs(float64(p.X - t.X))
	sy := math.Abs(float64(p.Y - t.Y))
	if sx >= sy {
		x := 1
		if p.X > t.X {
			x = -1
		}
		return &Direction{
			Point: Point{
				X: x,
				Y: 0,
			},
		}
	}
	y := 1
	if p.Y > t.Y {
		y = -1
	}
	return &Direction{
		Point: Point{
			X: 0,
			Y: y,
		},
	}
}

type Direction struct {
	Point
}

func (d *Direction) Mul(step int) *Direction {
	d.X = d.X * step
	d.Y = d.Y * step
	return d
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

func (d *Direction) IsReverse(t *Direction) bool {
	return d.X == 0-t.X && d.Y == 0-t.Y
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

func (d *Direction) Equals(t *Direction) bool {
	if (d.X == 0 && t.X == 0) || (d.Y == 0 && t.Y == 0) {
		return true
	}
	if (d.X == 0 && t.X != 0) || (d.Y == 0 && t.Y != 0) {
		return false
	}
	return d.X/t.X == d.Y/t.Y
}
