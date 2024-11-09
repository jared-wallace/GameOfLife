package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
	"sync"
)

// Game implements the ebiten.Game interface.
type Game struct {
	width, height int
	cells         [][]bool
	colors        [][]color.RGBA
	nextCells     [][]bool
	configIndex   int
	generation    int

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
		width:      width,
		height:     height,
		cells:      cells,
		colors:     colors,
		nextCells:  nextCells,
		generation: 0,
		cellSize:   4, // Default cell size
	}

	g.RandomConfig()

	return g
}

// RandomConfig initializes the grid with a random configuration.
func (g *Game) RandomConfig() {
	g.generation = 0
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			alive := rand.Float64() < 0.2 // 20% chance to be alive
			g.cells[y][x] = alive
			if alive {
				g.colors[y][x] = color.RGBA{
					R: uint8(rand.Intn(256)),
					G: uint8(rand.Intn(256)),
					B: uint8(rand.Intn(256)),
					A: 255,
				}
			} else {
				// Reset color if cell is dead
				g.colors[y][x] = color.RGBA{A: 255}
			}
		}
	}
}

// GliderConfig initializes the grid with a glider pattern at the center.
func (g *Game) GliderConfig() {
	// Clear the cells and reset colors
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.cells[y][x] = false
			g.colors[y][x] = color.RGBA{A: 255}
		}
	}
	g.generation = 0
	// Place a glider in the center
	midX := g.width / 2
	midY := g.height / 2
	glider := [5][2]int{
		{0, 1},
		{1, 2},
		{2, 0},
		{2, 1},
		{2, 2},
	}
	for _, pos := range glider {
		x := midX + pos[0]
		y := midY + pos[1]
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

// GunConfig initializes the grid with a Gosper Glider Gun pattern at the center.
func (g *Game) GunConfig() {
	// Clear the cells and reset colors
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.cells[y][x] = false
			g.colors[y][x] = color.RGBA{A: 255}
		}
	}
	g.generation = 0
	// Coordinates for the Gosper Glider Gun
	gunPattern := [36][2]int{
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

// countAliveNeighbors returns the number of alive neighbors for a cell at (x, y).
func (g *Game) countAliveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		ny := y + dy
		if ny < 0 || ny >= g.height {
			continue
		}
		for dx := -1; dx <= 1; dx++ {
			nx := x + dx
			if nx < 0 || nx >= g.width {
				continue
			}
			if dx == 0 && dy == 0 {
				continue
			}
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
		g.configIndex = (g.configIndex + 1) % 3
		switch g.configIndex {
		case 0:
			g.RandomConfig()
		case 1:
			g.GliderConfig()
		case 2:
			g.GunConfig()
		}
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
				//ebitenutil.DrawRect(screen, float64(rectX), float64(rectY), float64(cellSize), float64(cellSize), col)
				vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(cellSize), float32(cellSize), col, true)
			}
		}
	}

	// Display FPS and current configuration
	info := fmt.Sprintf(
		"FPS: %.2f\nConfig: %d/3\nCell Size: %d\nGeneration: %d\nPress SPACE to change config\nPress '+'/'-' to adjust cell size",
		ebiten.ActualFPS(),
		g.configIndex+1,
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
}

// main initializes and runs the game.
func main() {
	// Initial grid size
	initialGridWidth, initialGridHeight := 320, 256 // Adjusted for better performance

	// Create a new game instance
	game := NewGame(initialGridWidth, initialGridHeight)

	// Configure Ebiten window
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Set initial window size based on grid size and cell size
	ebiten.SetWindowSize(initialGridWidth*game.cellSize, initialGridHeight*game.cellSize)
	ebiten.SetWindowTitle("Conway's Game of Life")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
