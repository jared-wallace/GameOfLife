package PatternParser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// QuadtreeNode represents a node in the quadtree
type QuadtreeNode struct {
	Size     int              // Size of the node (2^depth)
	IsLeaf   bool             // Indicates if the node is a leaf
	Alive    [][]bool         // For leaf nodes: 8x8 grid
	Children [4]*QuadtreeNode // For non-leaf nodes: NW, NE, SW, SE
}

// ReadMCMacrocellFromFile reads a Macrocell (.mc) file and returns live cell coordinates, rule, and generation
func ReadMCMacrocellFromFile(filePath string) ([]Coordinate, string, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to open MC file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse Header
	if !scanner.Scan() {
		return nil, "", 0, fmt.Errorf("empty MC file")
	}
	firstLine := scanner.Text()
	if !strings.HasPrefix(firstLine, "[M2]") {
		return nil, "", 0, fmt.Errorf("invalid MC file format: missing [M2] header")
	}

	var rule string
	var generation int

	// Initialize nodes slice with node 0 as nil
	nodes := []*QuadtreeNode{nil} // nodes[0] = nil

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			// Handle special comments
			if strings.HasPrefix(line, "#R") {
				// Rule definition
				parts := strings.Fields(line[2:])
				if len(parts) >= 1 {
					rule = strings.Join(parts, " ")
				}
			} else if strings.HasPrefix(line, "#G") {
				// Generation count
				genStr := strings.TrimSpace(line[2:])
				gen, err := strconv.Atoi(genStr)
				if err == nil {
					generation = gen
				}
			}
			continue
		}
		// Start of quadtree nodes

		// Process the first node line
		if line != "" {
			if strings.HasPrefix(line, ".") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "$") {
				// Leaf node
				alive, err := parseLeafNode(line)
				if err != nil {
					return nil, "", 0, fmt.Errorf("error parsing leaf node: %v", err)
				}
				node := &QuadtreeNode{
					IsLeaf: true,
					Alive:  alive,
				}
				nodes = append(nodes, node)
			} else {
				// Non-leaf node
				parts := strings.Fields(line)
				if len(parts) != 5 {
					return nil, "", 0, fmt.Errorf("invalid non-leaf node format: %s", line)
				}
				log2Size, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, "", 0, fmt.Errorf("invalid log2 size in node: %s", line)
				}
				nw, err1 := strconv.Atoi(parts[1])
				ne, err2 := strconv.Atoi(parts[2])
				sw, err3 := strconv.Atoi(parts[3])
				se, err4 := strconv.Atoi(parts[4])
				if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
					return nil, "", 0, fmt.Errorf("invalid child node numbers in node: %s", line)
				}
				// Validate child node numbers
				if nw >= len(nodes) || ne >= len(nodes) || sw >= len(nodes) || se >= len(nodes) {
					return nil, "", 0, fmt.Errorf("child node number out of range in node: %s", line)
				}
				node := &QuadtreeNode{
					IsLeaf: false,
					Children: [4]*QuadtreeNode{
						nodes[nw],
						nodes[ne],
						nodes[sw],
						nodes[se],
					},
					Size: 1 << log2Size, // Size = 2^log2Size
				}
				nodes = append(nodes, node)
			}
		}
	}

	// Continue parsing remaining nodes
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		if strings.HasPrefix(line, "#") {
			// Skip comments within quadtree (if any)
			continue
		}

		if strings.HasPrefix(line, ".") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "$") {
			// Leaf node
			alive, err := parseLeafNode(line)
			if err != nil {
				return nil, "", 0, fmt.Errorf("error parsing leaf node: %v", err)
			}
			node := &QuadtreeNode{
				IsLeaf: true,
				Alive:  alive,
			}
			nodes = append(nodes, node)
		} else {
			// Non-leaf node
			parts := strings.Fields(line)
			if len(parts) != 5 {
				return nil, "", 0, fmt.Errorf("invalid non-leaf node format: %s", line)
			}
			log2Size, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, "", 0, fmt.Errorf("invalid log2 size in node: %s", line)
			}
			nw, err1 := strconv.Atoi(parts[1])
			ne, err2 := strconv.Atoi(parts[2])
			sw, err3 := strconv.Atoi(parts[3])
			se, err4 := strconv.Atoi(parts[4])
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				return nil, "", 0, fmt.Errorf("invalid child node numbers in node: %s", line)
			}
			// Validate child node numbers
			if nw >= len(nodes) || ne >= len(nodes) || sw >= len(nodes) || se >= len(nodes) {
				return nil, "", 0, fmt.Errorf("child node number out of range in node: %s", line)
			}
			node := &QuadtreeNode{
				IsLeaf: false,
				Children: [4]*QuadtreeNode{
					nodes[nw],
					nodes[ne],
					nodes[sw],
					nodes[se],
				},
				Size: 1 << log2Size, // Size = 2^log2Size
			}
			nodes = append(nodes, node)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", 0, fmt.Errorf("error reading MC file: %v", err)
	}

	if len(nodes) < 2 {
		return nil, "", 0, fmt.Errorf("no nodes found in MC file")
	}

	// The root node is the last node
	root := nodes[len(nodes)-1]

	// Reconstruct the universe grid from the quadtree
	// Determine the size of the universe
	universeSize := root.Size
	universe := make([][]bool, universeSize)
	for i := 0; i < universeSize; i++ {
		universe[i] = make([]bool, universeSize)
	}

	// The coordinate system: (0,0) at top-left
	// The upper left cell of the southeast child of the root node is at (0,1)
	// Adjust the origin accordingly if needed

	reconstructGrid(root, 0, 0, universe)

	// Extract live cell coordinates
	var liveCells []Coordinate
	for y := 0; y < universeSize; y++ {
		for x := 0; x < universeSize; x++ {
			if universe[y][x] {
				liveCells = append(liveCells, Coordinate{X: x, Y: y})
			}
		}
	}

	return liveCells, rule, generation, nil
}

// reconstructGrid fills the universe grid based on the quadtree
func reconstructGrid(node *QuadtreeNode, x, y int, universe [][]bool) {
	if node == nil {
		// Node 0: No live cells in this quadrant
		return
	}

	if node.IsLeaf {
		// Each leaf node represents an 8x8 grid
		for dy := 0; dy < 8; dy++ {
			for dx := 0; dx < 8; dx++ {
				if y+dy >= len(universe) || x+dx >= len(universe[0]) {
					// Prevent out-of-bounds errors
					continue
				}
				universe[y+dy][x+dx] = node.Alive[dy][dx]
			}
		}
	} else {
		// Non-leaf node: divide the area into four quadrants
		childSize := node.Size / 2
		reconstructGrid(node.Children[0], x, y, universe)                     // NW
		reconstructGrid(node.Children[1], x+childSize, y, universe)           // NE
		reconstructGrid(node.Children[2], x, y+childSize, universe)           // SW
		reconstructGrid(node.Children[3], x+childSize, y+childSize, universe) // SE
	}
}

// parseLeafNode parses a leaf node line for two-state algorithms
func parseLeafNode(line string) ([][]bool, error) {
	// Each leaf node represents an 8x8 grid
	// The line contains characters '.', '*', '$'
	// '.' = dead cell
	// '*' = alive cell
	// '$' = end of line
	// Empty cells at the end of each row are suppressed

	grid := make([][]bool, 8)
	for i := 0; i < 8; i++ {
		grid[i] = make([]bool, 8)
	}

	row := 0
	col := 0

	for _, char := range line {
		if char == '$' {
			row++
			col = 0
			if row >= 8 {
				break
			}
			continue
		}
		if row >= 8 || col >= 8 {
			continue // Ignore extra cells
		}
		switch char {
		case '*':
			grid[row][col] = true
		case '.':
			grid[row][col] = false
		default:
			return nil, fmt.Errorf("invalid character '%c' in leaf node", char)
		}
		col++
	}

	return grid, nil
}
