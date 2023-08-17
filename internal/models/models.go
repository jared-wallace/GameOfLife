package models

import "github.com/faiface/pixel"

type Point struct {
	X int
	Y int
}

type Cell struct {
	Point
	Alive bool
	Color pixel.RGBA
}
