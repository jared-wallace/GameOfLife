package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.expedia.biz/jarwallace/gol/internal/models"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	_ "image/png"
	"math/rand"
	"time"
)

func main() {
	pixelgl.Run(run)
}

const REFRESH = 16
const winXMax = 2048
const winYMax = 1536
const concurrencyMax = 204

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Conway's Game of Life",
		Bounds: pixel.R(0, 0, winXMax, winYMax),
		VSync:  true,
		//Monitor: pixelgl.PrimaryMonitor(),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	atlas := text.NewAtlas(
		basicfont.Face7x13,
		[]rune{'x'},
	)

	currGen := getInitial()
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			currGen = getInitial()
		}
		win.Clear(colornames.Black)
		drawPop(atlas, win, currGen)
		time.Sleep(REFRESH * time.Millisecond) // Delay to observe the generation

		if !popExists(currGen) {
			break // Exit the loop if there are no alive cells
		}
		currGen = analyzePopConcurrent(currGen)
		win.Update()
	}
}

func popExists(pop map[models.Point]models.Cell) bool {
	for _, cell := range pop {
		if cell.Alive {
			return true
		}
	}
	return false
}

func analyzePop(pop map[models.Point]models.Cell) map[models.Point]models.Cell {
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

func analyzePopConcurrent(pop map[models.Point]models.Cell) map[models.Point]models.Cell {
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

func drawPop(atlas *text.Atlas, win *pixelgl.Window, pop map[models.Point]models.Cell) {
	t := text.New(pixel.V(0, 0), atlas)
	for point, cell := range pop {
		if !cell.Alive {
			continue
		}
		t.Clear()
		t.Color = cell.Color
		t.Dot = pixel.V(float64(point.X*10), float64(point.Y*10))
		_, _ = t.WriteString("x")
		t.Draw(win, pixel.IM.Scaled(t.Orig, 1))
	}
}

func getPoint(x int, y int) models.Point {
	return models.Point{
		X: x,
		Y: y,
	}
}

func getInitial() map[models.Point]models.Cell {
	res := map[models.Point]models.Cell{}
	min := 0
	maxX := winXMax / 10
	maxY := winYMax / 10
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20000; i++ {
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		pt := getPoint(xVal, yVal)
		res[pt] = models.Cell{Point: pt, Alive: true, Color: randomColor()}
	}
	return res
}

func randomColor() pixel.RGBA {
	return pixel.RGBA{
		R: float64(rand.Float32()),
		G: float64(rand.Float32()),
		B: float64(rand.Float32()),
		A: 1,
	}
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
