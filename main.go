package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	_ "image/png"
	"math/rand"
	"time"
)

func main() {
	pixelgl.Run(run)
}

const REFRESH = 10
const winXMax = 2048
const winYMax = 1536
const concurrencyMax = 500

type point struct {
	x int
	y int
}

type cell struct {
	point
	alive bool
	color pixel.RGBA
}

func getNeighborCount(p point, pop map[point]cell) int {
	count := 0
	offsets := []point{{x: -1, y: 0}, {x: 1, y: 0}, {x: 0, y: -1}, {x: 0, y: 1}, {x: -1, y: -1}, {x: 1, y: -1}, {x: -1, y: 1}, {x: 1, y: 1}}
	for _, offset := range offsets {
		neighbor := point{x: p.x + offset.x, y: p.y + offset.y}
		if pop[neighbor].alive {
			count++
		}
	}
	return count
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Conway's Game of Life",
		Bounds: pixel.R(0, 0, winXMax, winYMax),
		VSync:  true,
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
		//currGen = analyzePop(currGen)
		currGen = analyzePopConcurrent(currGen)
		win.Update()
	}
}

func popExists(pop map[point]cell) bool {
	for _, cell := range pop {
		if cell.alive {
			return true
		}
	}
	return false
}

func analyzePop(pop map[point]cell) map[point]cell {
	nextGen := map[point]cell{}
	for x := 0; x < winXMax/10; x++ {
		for y := 0; y < winYMax/10; y++ {
			p := point{x: x, y: y}
			n := getNeighborCount(p, pop)
			if n == 3 || (n == 2 && pop[p].point == p) {
				if _, exists := pop[p]; exists {
					nextGen[p] = pop[p]
				} else {
					nextGen[p] = cell{point: p, alive: true, color: randomColor()}
				}
			}
		}
	}
	return nextGen
}

func analyzePopConcurrent(pop map[point]cell) map[point]cell {
	nextGen := map[point]cell{}
	numGoroutines := concurrencyMax / (winXMax / 10)
	results := make(chan map[point]cell, numGoroutines)

	// calculate the step size for x based on the number of goroutines
	step := (winXMax / 10) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		// Calculate the range for each goroutine
		startX := i * step
		endX := startX + step

		go func(startX, endX int) {
			partialNextGen := map[point]cell{}
			for x := startX; x < endX; x++ {
				for y := 0; y < winYMax/10; y++ {
					p := point{x: x, y: y}
					n := getNeighborCount(p, pop)
					if n == 3 || (n == 2 && pop[p].point == p) {
						if _, exists := pop[p]; exists {
							partialNextGen[p] = pop[p]
						} else {
							partialNextGen[p] = cell{point: p, alive: true, color: randomColor()}
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

func drawPop(atlas *text.Atlas, win *pixelgl.Window, pop map[point]cell) {
	t := text.New(pixel.V(0, 0), atlas)
	for point, cell := range pop {
		if !cell.alive {
			continue
		}
		t.Clear()
		t.Color = cell.color
		t.Dot = pixel.V(float64(point.x*10), float64(point.y*10))
		_, _ = t.WriteString("x")
		t.Draw(win, pixel.IM.Scaled(t.Orig, 1))
	}
}

func getPoint(x int, y int) point {
	return point{
		x: x,
		y: y,
	}
}

func getInitial() map[point]cell {
	res := map[point]cell{}
	min := 0
	maxX := winXMax / 10
	maxY := winYMax / 10
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20000; i++ {
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		pt := getPoint(xVal, yVal)
		res[pt] = cell{point: pt, alive: true, color: randomColor()}
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
