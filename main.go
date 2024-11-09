package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
	"sync"
)

type Game struct {
	width, height int
	cells         [][]bool
	colors        [][]color.RGBA
	nextCells     [][]bool
	configIndex   int
}

func NewGame(width, height int) *Game {
	cells := make([][]bool, height)
	colors := make([][]color.RGBA, height)
	nextCells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
		colors[i] = make([]color.RGBA, width)
		nextCells[i] = make([]bool, width)
	}

	g := &Game{
		width:     width,
		height:    height,
		cells:     cells,
		colors:    colors,
		nextCells: nextCells,
	}

	g.RandomConfig()

	return g
}

func (g *Game) RandomConfig() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			alive := rand.Float64() < 0.5
			g.cells[y][x] = alive
			if alive {
				g.colors[y][x] = color.RGBA{
					R: uint8(rand.Intn(256)),
					G: uint8(rand.Intn(256)),
					B: uint8(rand.Intn(256)),
					A: 255,
				}
			}
		}
	}
}

func (g *Game) GliderConfig() {
	// Clear the cells
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.cells[y][x] = false
		}
	}
	// Place a glider in the center
	midX := g.width / 2
	midY := g.height / 2
	g.cells[midY][midX+1] = true
	g.cells[midY+1][midX+2] = true
	g.cells[midY+2][midX] = true
	g.cells[midY+2][midX+1] = true
	g.cells[midY+2][midX+2] = true
	// Assign colors
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.cells[y][x] {
				g.colors[y][x] = color.RGBA{
					R: uint8(rand.Intn(256)),
					G: uint8(rand.Intn(256)),
					B: uint8(rand.Intn(256)),
					A: 255,
				}
			}
		}
	}
}

func (g *Game) GunConfig() {
	// Clear the cells
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.cells[y][x] = false
		}
	}
	// Coordinates for the Gosper Glider Gun
	// Relative positions
	gunPattern := [][2]int{
		{5, 1}, {5, 2}, {6, 1}, {6, 2},
		{5, 11}, {6, 11}, {7, 11}, {4, 12}, {8, 12}, {3, 13}, {9, 13},
		{3, 14}, {9, 14}, {6, 15}, {4, 16}, {8, 16}, {5, 17}, {6, 17},
		{7, 17}, {6, 18}, {3, 21}, {4, 21}, {5, 21}, {3, 22}, {4, 22},
		{5, 22}, {2, 23}, {6, 23}, {1, 25}, {2, 25}, {6, 25}, {7, 25},
		{3, 35}, {4, 35}, {3, 36}, {4, 36},
	}
	offsetX := g.width/2 - 18
	offsetY := g.height/2 - 9
	for _, p := range gunPattern {
		x := p[0] + offsetX
		y := p[1] + offsetY
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			g.cells[y][x] = true
			g.colors[y][x] = color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 255,
			}
		}
	}
}

func (g *Game) countAliveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := (x + dx + g.width) % g.width
			ny := (y + dy + g.height) % g.height
			if g.cells[ny][nx] {
				count++
			}
		}
	}
	return count
}

func (g *Game) Update() error {
	// Handle input
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// Switch configurations
		g.configIndex = (g.configIndex + 1) % 3
		switch g.configIndex {
		case 0:
			g.RandomConfig()
		case 1:
			g.GliderConfig()
		case 2:
			g.GunConfig()
		}
		return nil
	}

	// Create a wait group for concurrency
	var wg sync.WaitGroup
	numWorkers := 8
	rowsPerWorker := g.height / numWorkers

	for w := 0; w < numWorkers; w++ {
		startY := w * rowsPerWorker
		endY := startY + rowsPerWorker
		if w == numWorkers-1 {
			endY = g.height
		}
		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := 0; x < g.width; x++ {
					aliveNeighbors := g.countAliveNeighbors(x, y)
					if g.cells[y][x] {
						// Cell is alive
						if aliveNeighbors < 2 || aliveNeighbors > 3 {
							g.nextCells[y][x] = false
						} else {
							g.nextCells[y][x] = true
						}
					} else {
						// Cell is dead
						if aliveNeighbors == 3 {
							g.nextCells[y][x] = true
							// Assign a color to the new cell
							g.colors[y][x] = color.RGBA{
								R: uint8(rand.Intn(256)),
								G: uint8(rand.Intn(256)),
								B: uint8(rand.Intn(256)),
								A: 255,
							}
						} else {
							g.nextCells[y][x] = false
						}
					}
				}
			}
		}(startY, endY)
	}

	wg.Wait()

	// Swap cells and nextCells
	g.cells, g.nextCells = g.nextCells, g.cells

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.cells[y][x] {
				screen.Set(x, y, g.colors[y][x])
			} else {
				screen.Set(x, y, color.Black)
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nConfig: %d", ebiten.ActualFPS(), g.configIndex))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Adjust the grid size based on the window size
	if outsideWidth != g.width || outsideHeight != g.height {
		g.width = outsideWidth
		g.height = outsideHeight
		g.resizeGrid(g.width, g.height)
	}
	return g.width, g.height
}

func (g *Game) resizeGrid(newWidth, newHeight int) {
	newCells := make([][]bool, newHeight)
	newColors := make([][]color.RGBA, newHeight)
	newNextCells := make([][]bool, newHeight)
	for y := 0; y < newHeight; y++ {
		newCells[y] = make([]bool, newWidth)
		newColors[y] = make([]color.RGBA, newWidth)
		newNextCells[y] = make([]bool, newWidth)
		if y < len(g.cells) {
			copy(newCells[y], g.cells[y])
			copy(newColors[y], g.colors[y])
			copy(newNextCells[y], g.nextCells[y])
		}
	}
	g.cells = newCells
	g.colors = newColors
	g.nextCells = newNextCells
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Conway's Game of Life")
	width, height := 1280, 1024
	game := NewGame(width, height)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
