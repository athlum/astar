package astar

import (
	"sync"
)

type (
	BlockFunc func(int) bool

	Engine struct {
		isBlock BlockFunc
	}
)

func New(f BlockFunc) *Engine {
	return &Engine{
		isBlock: f,
	}
}

func (e *Engine) Router(src, dst *Point) *Router {
	return &Router{
		Src:      src,
		Dst:      dst,
		engine:   e,
		routeMap: &sync.Map{},
	}
}
