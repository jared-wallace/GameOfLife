package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	_ "image/png"
	"time"
)

func main() {
	pixelgl.Run(run)
}

const REFRESH = 500

type point struct {
	x int
	y int
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
	init := getInitial(atlas)
	for !win.Closed() {
		win.Clear(colornames.Navy)
		for _, point := range init {
			point.Draw(win, pixel.IM.Scaled(point.Orig, 1))
		}
		// main loop
		currGen := init
		for popExists(currGen) {
			time.Sleep(REFRESH * time.Millisecond)
			nextGen := analyzePop(atlas)
			win.Clear(colornames.Navy)
			for _, point := range nextGen {
				point.Draw(win, pixel.IM.Scaled(point.Orig, 1))
			}
			win.Update()
		}

		win.Update()
	}
}

func drawPop(atlas *text.Atlas, win *pixelgl.Window, pop []point) {
	for _, point := range pop {
		x := point.x
		y := point.y
		t := text.New(pixel.V(float64(x*10), float64(y*10)), atlas)
		_, _ = t.WriteString("x")
		t.Draw(win, pixel.IM.Scaled(t.Orig, 1))
	}
}

func getInitial(atlas *text.Atlas) []*text.Text {
	res := make([]*text.Text, 0)
	for x := 10; x < 1020; x += 10 {
		for y := 10; y < 760; y += 10 {
			t := text.New(pixel.V(float64(x), float64(y)), atlas)
			t.WriteString("x")
			res = append(res, t)
		}
	}
	return res
}
