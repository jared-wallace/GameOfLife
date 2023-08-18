package models

import (
	"github.expedia.biz/jarwallace/gol/internal/display"
)

func InjectRPentomino() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	// set our bounds
	minX := init.X
	minY := init.Y
	// create the R-pentamino
	pop := map[Point]Cell{}
	pop[Point{X: minX + 2, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 3, Y: minY + 2}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 3}] = Cell{Alive: true, Color: display.RandomColor()}
	pop[Point{X: minX + 3, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}

	pop[Point{X: minX + 4, Y: minY + 4}] = Cell{Alive: true, Color: display.RandomColor()}
	return pop

}

func InjectAcorn() map[Point]Cell {
	init := Point{
		X: 80,
		Y: 50,
	}
	// set our bounds
	minX := init.X
	minY := init.Y
	// create the acorn
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
	// "gosper" glider gun
	// grid required is 11 x 38
	// pick a random starting point towards the top left
	init := Point{
		X: 10,
		Y: 80,
	}
	// set our bounds
	minX := init.X
	minY := init.Y
	// create our glider gun
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

func getRandomPoint() Point {
	return Point{
		X: 40, // chosen by fair dice roll
		Y: 40, // guaranteed to be random
	}
}
