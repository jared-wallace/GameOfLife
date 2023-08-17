package processor

import (
	"github.com/faiface/pixel"
	"github.expedia.biz/jarwallace/gol/internal/models"
	"math/rand"
)

const concurrencyMax = 204 // needs to be a factor of winXMax / 10

func analyzePop(pop map[models.Point]models.Cell, winXMax, winYMax int) map[models.Point]models.Cell {
	nextGen := map[models.Point]models.Cell{}
	for x := 0; x < winXMax/10; x++ {
		for y := 0; y < winYMax/10; y++ {
			p := models.Point{X: x, Y: y}
			n := getNeighborCount(p, pop)
			if n == 3 || (n == 2 && pop[p].Point == p) {
				if _, exists := pop[p]; exists {
					nextGen[p] = pop[p]
				} else {
					nextGen[p] = models.Cell{Point: p, Alive: true, Color: randomColor()}
				}
			}
		}
	}
	return nextGen
}

func AnalyzePopConcurrent(pop map[models.Point]models.Cell, winXMax, winYMax int) map[models.Point]models.Cell {
	nextGen := map[models.Point]models.Cell{}
	numGoroutines := concurrencyMax / (winXMax / 10)
	results := make(chan map[models.Point]models.Cell, numGoroutines)

	// calculate the step size for x based on the number of goroutines
	step := (winXMax / 10) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		// Calculate the range for each goroutine
		startX := i * step
		endX := startX + step

		go func(startX, endX int) {
			partialNextGen := map[models.Point]models.Cell{}
			for x := startX; x < endX; x++ {
				for y := 0; y < winYMax/10; y++ {
					p := models.Point{X: x, Y: y}
					n := getNeighborCount(p, pop)
					if n == 3 || (n == 2 && pop[p].Point == p) {
						if _, exists := pop[p]; exists {
							partialNextGen[p] = pop[p]
						} else {
							partialNextGen[p] = models.Cell{Point: p, Alive: true, Color: randomColor()}
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

func randomColor() pixel.RGBA {
	return pixel.RGBA{
		R: float64(rand.Float32()),
		G: float64(rand.Float32()),
		B: float64(rand.Float32()),
		A: 1,
	}
}
