package main

import (
	"fmt"
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

const REFRESH = 500

type point struct {
	x     int
	y     int
	alive bool
}

func (p *point) getNeighborCount(pop []point) int {
	count := 0
	for _, point := range pop {
		if point.x == p.x && point.y == p.y {
			continue
		}
		if (point.x == p.x-1 || point.x == p.x+1 || point.x == p.x) &&
			(point.y == p.y-1 || point.y == p.y+1 || point.y == p.y) {
			if point.alive {
				count++
			}
		}
	}
	fmt.Println("Neighbor count: ", count, " and pop: ", pop)
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
	init := getInitial()
	for !win.Closed() {
		win.Clear(colornames.Navy)
		drawPop(atlas, win, init)
		// main loop
		currGen := init
		for popExists(currGen) {
			time.Sleep(REFRESH * time.Millisecond)
			nextGen := analyzePop(currGen)
			win.Clear(colornames.Navy)
			drawPop(atlas, win, nextGen)
			currGen = nextGen
			win.Update()
		}

		win.Update()
	}
}

func popExists(pop []point) bool {
	fmt.Println("Checking if pop exists")
	return len(pop) > 0
}

func analyzePop(pop []point) []point {
	fmt.Println("Analyzing")
	currCount := 0
	for _, p := range pop {
		if p.alive {
			currCount++
		}
	}
	fmt.Println("Current count: ", currCount)
	for _, p := range pop {
		n := p.getNeighborCount(pop)
		switch n {
		case 0, 1:
			p.alive = false
		case 2:
			fallthrough
		case 3:
			p.alive = true
		default:
			p.alive = false
		}
	}
	currCount = 0
	for _, p := range pop {
		if p.alive {
			currCount++
		}
	}
	fmt.Println("New count: ", currCount)
	return pop
}

func drawPop(atlas *text.Atlas, win *pixelgl.Window, pop []point) {
	fmt.Println("Drawing")
	for _, point := range pop {
		if !point.alive {
			continue
		}
		x := point.x
		y := point.y
		t := text.New(pixel.V(float64(x*10), float64(y*10)), atlas)
		_, _ = t.WriteString("x")
		t.Draw(win, pixel.IM.Scaled(t.Orig, 1))
	}
}

func getInitial() []point {
	res := make([]point, 0)
	min := 0
	maxX := 102
	maxY := 76
	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		t := point{
			x:     xVal,
			y:     yVal,
			alive: true,
		}
		res = append(res, t)
	}
	rand.Seed(time.Now().UnixNano())
	xVal := rand.Intn(maxX-min+1) + min
	yVal := rand.Intn(maxY-min+1) + min
	t := point{
		x:     xVal,
		y:     yVal,
		alive: true,
	}
	res = append(res, t)
	return res
}
