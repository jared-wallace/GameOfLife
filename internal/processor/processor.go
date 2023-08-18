package processor

import (
	"github.expedia.biz/jarwallace/gol/internal/display"
	"github.expedia.biz/jarwallace/gol/internal/models"
)

type processor struct {
	winXMax        int
	winYMax        int
	concurrencyMax int
}

func NewProcessor(winXMax, winYMax int) *processor {
	concurrencyMax := factors(winXMax / 10)[len(factors(winXMax/10))-1]
	return &processor{
		winXMax:        winXMax,
		winYMax:        winYMax,
		concurrencyMax: concurrencyMax,
	}
}

func factors(n int) []int {
	var result []int
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

func (pr *processor) AnalyzePop(pop map[models.Point]models.Cell) map[models.Point]models.Cell {
	nextGen := map[models.Point]models.Cell{}
	for x := 0; x < pr.winXMax/10; x++ {
		for y := 0; y < pr.winYMax/10; y++ {
			p := models.Point{X: x, Y: y}
			n := getNeighborCount(p, pop)
			if n == 3 || (n == 2 && pop[p].Alive) {
				if _, exists := pop[p]; exists {
					nextGen[p] = pop[p]
				} else {
					nextGen[p] = models.Cell{Point: p, Alive: true, Color: display.RandomColor()}
				}
			}
		}
	}
	return nextGen
}

func (pr *processor) AnalyzePopEfficiently(pop map[models.Point]models.Cell) map[models.Point]models.Cell {
	nextGen := map[models.Point]models.Cell{}

	// Create a set of points to analyze (alive cells and their neighbors)
	toAnalyze := map[models.Point]struct{}{}
	for p := range pop {
		toAnalyze[p] = struct{}{}
		for _, offset := range getNeighborOffsets() {
			toAnalyze[models.Point{X: p.X + offset.X, Y: p.Y + offset.Y}] = struct{}{}
		}
	}

	pointsList := make([]models.Point, 0, len(toAnalyze))
	for p := range toAnalyze {
		pointsList = append(pointsList, p)
	}

	numGoroutines := pr.concurrencyMax
	results := make(chan map[models.Point]models.Cell, numGoroutines)

	// Calculate number of points each goroutine will handle
	pointsPerGoroutine := len(pointsList) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		startIdx := i * pointsPerGoroutine
		endIdx := startIdx + pointsPerGoroutine
		if i == numGoroutines-1 { // last goroutine might handle more
			endIdx = len(pointsList)
		}

		go func(startIdx, endIdx int) {
			partialNextGen := map[models.Point]models.Cell{}
			for idx := startIdx; idx < endIdx; idx++ {
				p := pointsList[idx]
				n := getNeighborCount(p, pop)
				if n == 3 || (n == 2 && pop[p].Point == p) {
					if _, exists := pop[p]; exists {
						partialNextGen[p] = pop[p]
					} else {
						partialNextGen[p] = models.Cell{Point: p, Alive: true, Color: display.RandomColor()}
					}
				}
			}
			results <- partialNextGen
		}(startIdx, endIdx)
	}

	// Merge results from all goroutines
	for i := 0; i < numGoroutines; i++ {
		partialNextGen := <-results
		for k, v := range partialNextGen {
			nextGen[k] = v
		}
	}

	return nextGen
}

func getNeighborOffsets() []models.Point {
	return []models.Point{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
		{X: -1, Y: -1},
		{X: 1, Y: -1},
		{X: -1, Y: 1},
		{X: 1, Y: 1},
	}
}

func (pr *processor) AnalyzePopConcurrent(pop map[models.Point]models.Cell) map[models.Point]models.Cell {
	nextGen := map[models.Point]models.Cell{}
	numGoroutines := pr.concurrencyMax
	results := make(chan map[models.Point]models.Cell, numGoroutines)

	// calculate the step size for x based on the number of goroutines
	step := (pr.winXMax / 10) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		// Calculate the range for each goroutine
		startX := i * step
		endX := startX + step

		go func(startX, endX int) {
			partialNextGen := map[models.Point]models.Cell{}
			for x := startX; x < endX; x++ {
				for y := 0; y < pr.winYMax/10; y++ {
					p := models.Point{X: x, Y: y}
					n := getNeighborCount(p, pop)
					if n == 3 || (n == 2 && pop[p].Point == p) {
						if _, exists := pop[p]; exists {
							partialNextGen[p] = pop[p]
						} else {
							partialNextGen[p] = models.Cell{Point: p, Alive: true, Color: display.RandomColor()}
						}
					}
				}
			}
			results <- partialNextGen
		}(startX, endX)
	}

	// Merge results from all goroutines
	for i := 0; i < numGoroutines; i++ {
		partialNextGen := <-results
		for k, v := range partialNextGen {
			nextGen[k] = v
		}
	}

	return nextGen
}

func getNeighborCount(p models.Point, pop map[models.Point]models.Cell) int {
	count := 0
	offsets := []models.Point{{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: 1}}
	for _, offset := range offsets {
		neighbor := models.Point{X: p.X + offset.X, Y: p.Y + offset.Y}
		if pop[neighbor].Alive {
			count++
		}
	}
	return count
}
