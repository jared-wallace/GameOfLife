package patterns

import (
	"image/color"
	"math/rand"
)

type PatternGenerator struct {
	patterns []func(int, int) ([][]bool, [][]color.RGBA, string)
	height   int
	width    int
}

func NewPatternGenerator(height, width int) *PatternGenerator {
	pg := PatternGenerator{}
	pg.patterns = append(pg.patterns, RandomConfig)
	pg.patterns = append(pg.patterns, GliderConfig)
	pg.patterns = append(pg.patterns, GunConfig)
	pg.patterns = append(pg.patterns, BlockConfig)
	pg.patterns = append(pg.patterns, BlinkerConfig)
	pg.patterns = append(pg.patterns, ToadConfig)
	pg.patterns = append(pg.patterns, BeaconConfig)
	pg.patterns = append(pg.patterns, LWSSConfig)
	pg.patterns = append(pg.patterns, PulsarConfig)

	pg.height = height
	pg.width = width
	return &pg
}

func (pg *PatternGenerator) GetConfig(idx int) ([][]bool, [][]color.RGBA, string) {
	return pg.patterns[idx](pg.height, pg.width)
}

func (pg *PatternGenerator) GetPatternCount() int {
	return len(pg.patterns)
}

func (pg *PatternGenerator) SetHW(height int, width int) {
	pg.height = height
	pg.width = width
}

// initializeBoard creates a new board with all cells dead and colors set to default.
func initializeBoard(height, width int) ([][]bool, [][]color.RGBA) {
	cells := make([][]bool, height)
	colors := make([][]color.RGBA, height)
	defaultColor := color.RGBA{A: 255}

	for i := range cells {
		cells[i] = make([]bool, width)
		colors[i] = make([]color.RGBA, width)
		for j := range colors[i] {
			colors[i][j] = defaultColor
		}
	}
	return cells, colors
}

// randomColor generates a random RGB color with full opacity.
func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

// setAliveCells sets the specified cells to alive and assigns them random colors.
func setAliveCells(cells [][]bool, colors [][]color.RGBA, positions [][2]int, offsetX, offsetY int) {
	for _, pos := range positions {
		x := pos[0] + offsetX
		y := pos[1] + offsetY
		if y >= 0 && y < len(cells) && x >= 0 && x < len(cells[0]) {
			cells[y][x] = true
			colors[y][x] = randomColor()
		}
	}
}

// RandomConfig initializes the board with a random configuration.
func RandomConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float64() < 0.2 { // 20% chance to be alive
				cells[y][x] = true
				colors[y][x] = randomColor()
			}
		}
	}
	return cells, colors, "random"
}

// GliderConfig places a glider pattern in the center of the board.
func GliderConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	glider := [5][2]int{
		{0, 1},
		{1, 2},
		{2, 0},
		{2, 1},
		{2, 2},
	}
	setAliveCells(cells, colors, glider[:], midX, midY)
	return cells, colors, "glider"
}

// GunConfig places the Gosper Glider Gun pattern in the center of the board.
func GunConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	gunPattern := [36][2]int{
		{5, 1}, {5, 2}, {6, 1}, {6, 2},
		{5, 11}, {6, 11}, {7, 11}, {4, 12}, {8, 12}, {3, 13}, {9, 13},
		{3, 14}, {9, 14}, {6, 15}, {4, 16}, {8, 16}, {5, 17}, {6, 17},
		{7, 17}, {6, 18}, {3, 21}, {4, 21}, {5, 21}, {3, 22}, {4, 22},
		{5, 22}, {2, 23}, {6, 23}, {1, 25}, {2, 25}, {6, 25}, {7, 25},
		{3, 35}, {4, 35}, {3, 36}, {4, 36},
	}
	offsetX, offsetY := width/2-18, height/2-9
	setAliveCells(cells, colors, gunPattern[:], offsetX, offsetY)
	return cells, colors, "gosper gun"
}

// BlockConfig places a 2x2 block pattern in the center of the board.
func BlockConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	block := [4][2]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}
	setAliveCells(cells, colors, block[:], midX, midY)
	return cells, colors, "block"
}

// BlinkerConfig places a blinker (horizontal line) in the center of the board.
func BlinkerConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	blinker := [3][2]int{
		{-1, 0},
		{0, 0},
		{1, 0},
	}
	setAliveCells(cells, colors, blinker[:], midX, midY)
	return cells, colors, "blinker"
}

// ToadConfig places an oscillator in the center of the board
func ToadConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	toad := [6][2]int{
		{0, -1},
		{1, -1},
		{2, -1},
		{-1, 0},
		{0, 0},
		{1, 0},
	}
	setAliveCells(cells, colors, toad[:], midX, midY)
	return cells, colors, "toad"
}

func BeaconConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	beacon := [6][2]int{
		{0, 0}, {1, 0},
		{0, 1},
		{2, 2}, {3, 2},
		{3, 3},
	}
	setAliveCells(cells, colors, beacon[:], midX, midY)
	return cells, colors, "beacon"
}

func LWSSConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX := width/2 - 2
	midY := height/2 - 1
	lwss := [9][2]int{
		{1, 0}, {4, 0},
		{0, 1}, {4, 1},
		{4, 2},
		{0, 3}, {1, 3}, {2, 3}, {3, 3},
	}
	setAliveCells(cells, colors, lwss[:], midX, midY)
	return cells, colors, "lightweight space ship"
}

func PulsarConfig(height, width int) ([][]bool, [][]color.RGBA, string) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2
	pulsar := [48][2]int{
		{-4, -2}, {-4, -1}, {-4, 1}, {-4, 2},
		{-2, -4}, {-1, -4}, {1, -4}, {2, -4},
		{-3, -2}, {-3, 2}, {3, -2}, {3, 2},
		{-2, -3}, {-1, -3}, {1, -3}, {2, -3},
		{-2, 3}, {-1, 3}, {1, 3}, {2, 3},
		{-3, -1}, {-3, 1}, {3, -1}, {3, 1},
		{-1, -2}, {-1, 2}, {1, -2}, {1, 2},
		{-2, -1}, {-2, 1}, {2, -1}, {2, 1},
		{-1, -4}, {1, -4}, {-1, 4}, {1, 4},
		{-4, -3}, {-4, 3}, {4, -3}, {4, 3},
		{-3, -4}, {3, -4}, {-3, 4}, {3, 4},
		{-4, -2}, {4, -2}, {-4, 2}, {4, 2},
	}
	setAliveCells(cells, colors, pulsar[:], midX, midY)
	return cells, colors, "pulsar"
}
