package engine

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jared-wallace/gol/patterns"
	"image/color"
	"math/rand"
	"sync"
)

// Game implements the ebiten.Game interface.
type Game struct {
	width, height    int
	cells            [][]bool
	colors           [][]color.RGBA
	nextCells        [][]bool
	configIndex      int
	generation       int
	name             string
	patternGenerator *patterns.PatternGenerator

	// Fields for cell size management
	cellSize      int
	cellSizeMutex sync.Mutex

	// Fields for key state tracking
	prevSpacePressed bool
	prevPlusPressed  bool
	prevMinusPressed bool
}

// NewGame initializes a new Game instance.
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
		width:            width,
		height:           height,
		cells:            cells,
		colors:           colors,
		nextCells:        nextCells,
		generation:       0,
		patternGenerator: patterns.NewPatternGenerator(height, width),
		cellSize:         4, // Default cell size
	}

	g.cells, g.colors, g.name = patterns.RandomConfig(g.height, g.width)

	return g
}

func (g *Game) GetCellSize() int {
	return g.cellSize
}

// countAliveNeighbors returns the number of alive neighbors for a cell at (x, y).
func (g *Game) countAliveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue // Skip the cell itself
			}

			// Wrap around the edges using modulo arithmetic
			nx := (x + dx + g.width) % g.width
			ny := (y + dy + g.height) % g.height

			if g.cells[ny][nx] {
				count++
			}
		}
	}
	return count
}

// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Handle input: spacebar to switch configurations
	currentSpacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if currentSpacePressed && !g.prevSpacePressed {
		g.configIndex = (g.configIndex + 1) % g.patternGenerator.GetPatternCount()
		g.cells, g.colors, g.name = g.patternGenerator.GetConfig(g.configIndex)
		g.generation = 0
	}
	g.prevSpacePressed = currentSpacePressed

	// Handle input: '+' to increase cell size
	currentPlusPressed := ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd)
	if currentPlusPressed && !g.prevPlusPressed {
		g.cellSizeMutex.Lock()
		if g.cellSize < 20 { // Maximum cell size limit
			g.cellSize++
			// No direct call to resizeGrid here
		}
		g.cellSizeMutex.Unlock()
	}
	g.prevPlusPressed = currentPlusPressed

	// Handle input: '-' to decrease cell size
	currentMinusPressed := ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract)
	if currentMinusPressed && !g.prevMinusPressed {
		g.cellSizeMutex.Lock()
		if g.cellSize > 1 { // Minimum cell size limit
			g.cellSize--
			// No direct call to resizeGrid here
		}
		g.cellSizeMutex.Unlock()
	}
	g.prevMinusPressed = currentMinusPressed

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
	// increment generation
	g.generation++

	return nil
}

// Draw renders the current state to the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.cellSizeMutex.Lock()
	cellSize := g.cellSize
	g.cellSizeMutex.Unlock()

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.cells[y][x] {
				col := g.colors[y][x]
				rectX := x * cellSize
				rectY := y * cellSize
				// Draw a filled rectangle for the cell
				vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(cellSize), float32(cellSize), col, true)
			}
		}
	}

	// Display FPS and current configuration
	info := fmt.Sprintf(
		"FPS: %.2f\nConfig: %s\nCell Size: %d\nGeneration: %d\nPress SPACE to change config\nPress '+'/'-' to adjust cell size",
		ebiten.ActualFPS(),
		g.name,
		g.cellSize,
		g.generation,
	)
	ebitenutil.DebugPrint(screen, info)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.cellSizeMutex.Lock()
	cellSize := g.cellSize
	g.cellSizeMutex.Unlock()

	// Calculate grid dimensions based on cell size
	gridWidth := outsideWidth / cellSize
	gridHeight := outsideHeight / cellSize

	// Ensure grid dimensions are at least 1x1
	if gridWidth < 1 {
		gridWidth = 1
	}
	if gridHeight < 1 {
		gridHeight = 1
	}

	// If grid dimensions have changed, resize the grid
	if gridWidth != g.width || gridHeight != g.height {
		g.resizeGrid(gridWidth, gridHeight)
	}

	return outsideWidth, outsideHeight
}

// resizeGrid adjusts the grid size based on the new grid dimensions and current cell size.
func (g *Game) resizeGrid(newWidth, newHeight int) {
	// Create new slices with updated dimensions
	newCells := make([][]bool, newHeight)
	newColors := make([][]color.RGBA, newHeight)
	newNextCells := make([][]bool, newHeight)
	for y := 0; y < newHeight; y++ {
		newCells[y] = make([]bool, newWidth)
		newColors[y] = make([]color.RGBA, newWidth)
		newNextCells[y] = make([]bool, newWidth)
		// Copy existing data if within old bounds
		if y < g.height {
			for x := 0; x < newWidth && x < g.width; x++ {
				newCells[y][x] = g.cells[y][x]
				newColors[y][x] = g.colors[y][x]
				newNextCells[y][x] = g.nextCells[y][x]
			}
		}
	}

	// Replace old slices with new ones
	g.cells = newCells
	g.colors = newColors
	g.nextCells = newNextCells
	g.width = newWidth
	g.height = newHeight
	g.patternGenerator.SetHW(newHeight, newWidth)
}