package PatternParser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Coordinate represents the position of a live cell
type Coordinate struct {
	X int
	Y int
}

// ReadPatternFromFile reads a plaintext Game of Life pattern from a file
// and returns a slice of Coordinates where each Coordinate represents
// a live cell ('O') in the pattern.
func ReadPatternFromFile(filePath string) ([]Coordinate, error) {
	var coordinates []Coordinate

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Use a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, char := range line {
			if char == 'O' {
				coordinates = append(coordinates, Coordinate{X: x, Y: y})
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return coordinates, nil
}

// ParsePattern converts a slice of Coordinate to a slice of [2]int pairs.
func ParsePattern(coordinates []Coordinate) [][2]int {
	var coordPairs [][2]int
	for _, coord := range coordinates {
		coordPairs = append(coordPairs, [2]int{coord.X, coord.Y})
	}
	return coordPairs
}

// ReadRLEPatternFromFile reads an RLE format Game of Life pattern from a file
// and returns a slice of Coordinates where each Coordinate represents
// a live cell ('o') in the pattern, along with the xMax and yMax.
func ReadRLEPatternFromFile(filePath string) ([]Coordinate, int, int, error) {
	var coordinates []Coordinate

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var xSize, ySize int
	headerParsed := false

	// Read the header
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			// Skip comments and empty lines
			continue
		}
		if strings.HasPrefix(line, "x") {
			// Parse the header line
			headerRegex := regexp.MustCompile(`x\s*=\s*(\d+)\s*,\s*y\s*=\s*(\d+)`)
			matches := headerRegex.FindStringSubmatch(line)
			if len(matches) >= 3 {
				xSize, err = strconv.Atoi(matches[1])
				if err != nil {
					return nil, 0, 0, fmt.Errorf("invalid x size in header: %v", err)
				}
				ySize, err = strconv.Atoi(matches[2])
				if err != nil {
					return nil, 0, 0, fmt.Errorf("invalid y size in header: %v", err)
				}
				headerParsed = true
			} else {
				return nil, 0, 0, fmt.Errorf("invalid header line: %s", line)
			}
			break // Exit after parsing header
		}
	}

	if !headerParsed {
		return nil, 0, 0, fmt.Errorf("RLE header not found in file")
	}

	// Read the pattern data
	var patternLines []string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			// Skip comments and empty lines
			continue
		}
		patternLines = append(patternLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error reading file: %v", err)
	}

	patternData := strings.Join(patternLines, "")
	// Remove any whitespace
	patternData = strings.ReplaceAll(patternData, " ", "")
	// Now parse the pattern data

	x, y := 0, 0
	count := 0
	number := ""

	for i := 0; i < len(patternData); i++ {
		c := patternData[i]
		switch c {
		case 'b', 'o', '$', '!':
			if number != "" {
				count, err = strconv.Atoi(number)
				if err != nil {
					return nil, 0, 0, fmt.Errorf("invalid number in pattern data: %v", err)
				}
				number = ""
			} else {
				count = 1
			}

			switch c {
			case 'b':
				x += count
			case 'o':
				for j := 0; j < count; j++ {
					coordinates = append(coordinates, Coordinate{X: x, Y: y})
					x++
				}
			case '$':
				y += count
				x = 0
			case '!':
				// End of pattern
				return coordinates, xSize, ySize, nil
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			number += string(c)
		default:
			return nil, 0, 0, fmt.Errorf("unexpected character '%c' in pattern data", c)
		}
	}

	return coordinates, xSize, ySize, nil
}
