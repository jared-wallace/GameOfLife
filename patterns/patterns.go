package patterns

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/jared-wallace/gol/pkg/PatternParser"
)

// Coordinate represents the position of a live cell
type Coordinate struct {
	X int
	Y int
}

// PatternGenerator manages available patterns
type PatternGenerator struct {
	patterns []string
	height   int
	width    int
}

// NewPatternGenerator initializes the PatternGenerator by loading all patterns from the patterns/ directory
func NewPatternGenerator(height, width int) *PatternGenerator {
	pg := &PatternGenerator{
		height: height,
		width:  width,
	}

	// Load all pattern names from the patterns directory
	patternNames, err := loadPatternNames("patterns")
	if err != nil {
		log.Fatalf("Failed to load patterns: %v", err)
	}

	pg.patterns = append([]string{"random"}, patternNames...)
	return pg
}

// loadPatternNames scans the directory and returns a list of pattern names (without extension)
func loadPatternNames(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read patterns directory: %v", err)
	}

	var patternNames []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".txt") || strings.HasSuffix(file.Name(), ".rle")) {
			name := strings.TrimSuffix(file.Name(), ".txt")
			name = strings.TrimSuffix(name, ".rle")
			patternNames = append(patternNames, name)
		}
	}
	sort.Strings(patternNames)

	return patternNames, nil
}

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

// GetConfig loads the specified pattern by index
func (pg *PatternGenerator) GetConfig(idx int) ([][]bool, [][]color.RGBA, string, error) {
	if idx < 0 || idx >= len(pg.patterns) {
		return nil, nil, "", fmt.Errorf("pattern index %d out of range", idx)
	}
	if idx == 0 {
		// Return random configuration
		cells, colors, name := RandomConfig(pg.height, pg.width)
		return cells, colors, name, nil
	}

	patternName := pg.patterns[idx]
	return LoadPatternConfig(pg.height, pg.width, patternName)
}

// GetPatternCount returns the number of available patterns
func (pg *PatternGenerator) GetPatternCount() int {
	return len(pg.patterns)
}

// SetHW sets the height and width of the board
func (pg *PatternGenerator) SetHW(height int, width int) {
	pg.height = height
	pg.width = width
}

func LoadPatternConfig(height, width int, patternName string) ([][]bool, [][]color.RGBA, string, error) {
	cells, colors := initializeBoard(height, width)
	midX, midY := width/2, height/2

	// Try to open .txt file first, if not found, try .rle
	var filePath string
	if _, err := os.Stat(fmt.Sprintf("patterns/%s.txt", patternName)); err == nil {
		filePath = fmt.Sprintf("patterns/%s.txt", patternName)
	} else if _, err := os.Stat(fmt.Sprintf("patterns/%s.rle", patternName)); err == nil {
		filePath = fmt.Sprintf("patterns/%s.rle", patternName)
	} else {
		return nil, nil, "", fmt.Errorf("pattern file for '%s' not found", patternName)
	}

	var coordinates []PatternParser.Coordinate
	var err error
	if strings.HasSuffix(filePath, ".txt") {
		// Read and parse the plaintext pattern file
		coordinates, err = PatternParser.ReadPatternFromFile(filePath)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to read pattern '%s': %v", patternName, err)
		}
	} else if strings.HasSuffix(filePath, ".rle") {
		// Read and parse the RLE pattern file
		coordinates, _, _, err = PatternParser.ReadRLEPatternFromFile(filePath)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to read RLE pattern '%s': %v", patternName, err)
		}
	} else {
		return nil, nil, "", fmt.Errorf("unknown file extension for pattern '%s'", patternName)
	}

	// Convert to [][2]int
	coordPairs := PatternParser.ParsePattern(coordinates)

	// Set the alive cells on the board
	setAliveCells(cells, colors, coordPairs, midX, midY, width, height)

	return cells, colors, patternName, nil
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
func setAliveCells(cells [][]bool, colors [][]color.RGBA, positions [][2]int, offsetX, offsetY int, width, height int) {
	for _, pos := range positions {
		x := (pos[0] + offsetX + width) % width
		y := (pos[1] + offsetY + height) % height

		cells[y][x] = true
		colors[y][x] = randomColor()
	}
}
