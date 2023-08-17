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

const REFRESH = 100

type point struct {
	x int
	y int
}

func getNeighborCount(p point, pop map[point]bool) int {
	count := 0
	offsets := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for _, offset := range offsets {
		neighbor := point{p.x + offset.x, p.y + offset.y}
		if pop[neighbor] {
			count++
		}
	}
	return count
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Conway's Game of Life",
		Bounds: pixel.R(0, 0, 1024, 768),
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

	win.Clear(colornames.Navy)
	currGen := getInitial()
	for !win.Closed() {
		win.Clear(colornames.Navy)
		drawPop(atlas, win, currGen)
		win.Update()
		time.Sleep(REFRESH * time.Millisecond) // Delay to observe the generation

		if !popExists(currGen) {
			break // Exit the loop if there are no alive cells
		}
		currGen = analyzePop(currGen)
		win.Update()
	}
}

func popExists(pop map[point]bool) bool {
	for _, alive := range pop {
		if alive {
			return true
		}
	}
	return false
}

func analyzePop(pop map[point]bool) map[point]bool {
	nextGen := map[point]bool{}
	for x := 0; x < 102; x++ {
		for y := 0; y < 76; y++ {
			p := point{x: x, y: y}
			n := getNeighborCount(p, pop)
			if n == 3 || (n == 2 && pop[p]) {
				nextGen[p] = true
			}
		}
	}
	return nextGen
}

func drawPop(atlas *text.Atlas, win *pixelgl.Window, pop map[point]bool) {
	t := text.New(pixel.V(0, 0), atlas)
	for point, alive := range pop {
		if !alive {
			continue
		}
		t.Clear()
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

func getOscillator() map[point]bool {
	res := map[point]bool{}
	res[getPoint(10, 10)] = true
	res[getPoint(10, 11)] = true
	res[getPoint(10, 12)] = true
	return res
}

func getInitial() map[point]bool {
	//getOscillator()
	res := map[point]bool{}
	min := 0
	maxX := 102
	maxY := 76
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		pt := getPoint(xVal, yVal)
		if res[pt] {
			continue
		}
		res[pt] = true
	}
	return res
}
