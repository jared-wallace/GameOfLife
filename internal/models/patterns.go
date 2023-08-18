package models

import (
	"github.expedia.biz/jarwallace/gol/internal/display"
	"math/rand"
	"time"
)

const NumPatterns = 7

func GetPattern(num int, xMax, yMax int) (map[Point]Cell, string) {
	switch num {
	case 0:
		return InjectRPentomino(), "r-pentomino"
	case 1:
		return InjectOscillator(), "oscillator"
	case 2:
		return InjectOscillator2(), "oscillator2"
	case 3:
		return InjectBlockLayingSwitchEngine(), "block-laying switch engine"
	case 4:
		return InjectAcorn(), "acorn"
	case 5:
		return InjectGlider(), "gosper glider"
	case 6:
		return InjectRandom(xMax, yMax), "random"
	default:
		return InjectRPentomino(), "r-pentomino"
	}
}

func InjectRPentomino() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 3, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 4, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop

}

func InjectOscillator2() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 4, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 4, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 5, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 6, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 7, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 8, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 9, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 9, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 10, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 11, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop
}

func InjectOscillator() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 2, Y: minY + 9}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 3, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 8}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 10}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 11}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 4, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 4, Y: minY + 9}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 10, Y: minY + 10}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 10, Y: minY + 11}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 11, Y: minY + 8}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 11, Y: minY + 10}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 11, Y: minY + 11}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 12, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 13, Y: minY + 10}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 14, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 14, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 14, Y: minY + 9}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 15, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 15, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop

}

func InjectBlockLayingSwitchEngine() map[Point]Cell {
	init := Point{
		X: 50,
		Y: 50,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 4, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 5, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 6, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 7, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 8, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 9, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 11, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 12, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 13, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 14, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 15, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 19, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 20, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 21, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 28, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 29, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 30, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 31, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 32, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 33, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 34, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 36, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 37, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 38, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 39, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 40, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop

}

func InjectAcorn() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 3, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 5, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 6, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 7, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 8, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop
}

func InjectGlider() map[Point]Cell {
	init := Point{
		X: 10,
		Y: 80,
	}
	minX := init.X
	minY := init.Y
	pop := map[Point]Cell{}
	pop[Point{X: minX + 1, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 1, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 2, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 2, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 11, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 11, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 11, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 12, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 12, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 13, Y: minY + 1}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 13, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 14, Y: minY + 1}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 14, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 15, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 16, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 16, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 17, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 17, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 17, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 18, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 21, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 21, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 21, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 22, Y: minY + 5}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 22, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 22, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 23, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 23, Y: minY + 8}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 25, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 25, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 25, Y: minY + 8}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 25, Y: minY + 9}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 35, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 35, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 36, Y: minY + 6}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 36, Y: minY + 7}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop
}

func InjectRandom(winXMax, winYMax int) map[Point]Cell {
	res := map[Point]Cell{}
	min := 0
	maxX := winXMax / 10
	maxY := winYMax / 10
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20000; i++ {
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		pt := getPoint(xVal, yVal)
		res[pt] = Cell{Point: pt, Alive: true, Color: display.RandomColor()}
	}
	return res
}

func getPoint(x int, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}
