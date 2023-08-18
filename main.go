package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.expedia.biz/jarwallace/gol/internal/display"
	"github.expedia.biz/jarwallace/gol/internal/models"
	"github.expedia.biz/jarwallace/gol/internal/processor"
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

//const winXMax = 1536
//const winYMax = 1152

func run() {
	monitor := pixelgl.PrimaryMonitor()
	w, l := monitor.Size()
	winXMax := int(w)
	winYMax := int(l)
	cfg := pixelgl.WindowConfig{
		Title:     "Conway's Game of Life",
		Bounds:    pixel.R(0, 0, w, l),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	atlas := text.NewAtlas(
		basicfont.Face7x13,
		[]rune{'x'},
	)

	//currGen := getInitial(winXMax, winYMax)
	currGen := models.InjectGlider()
	pr := processor.NewProcessor(winXMax, winYMax)
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			currGen = getInitial(winXMax, winYMax)
		}
		win.Clear(colornames.Black)
		drawPop(atlas, win, currGen)
		time.Sleep(REFRESH * time.Millisecond) // Delay to observe the generation

		if !popExists(currGen) {
			break // Exit the loop if there are no alive cells
		}
		//currGen = pr.AnalyzePopEfficiently(currGen)
		currGen = pr.AnalyzePop(currGen)
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

func getInitial(winXMax, winYMax int) map[models.Point]models.Cell {
	res := map[models.Point]models.Cell{}
	min := 0
	maxX := winXMax / 10
	maxY := winYMax / 10
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20000; i++ {
		xVal := rand.Intn(maxX-min+1) + min
		yVal := rand.Intn(maxY-min+1) + min
		pt := getPoint(xVal, yVal)
		res[pt] = models.Cell{Point: pt, Alive: true, Color: display.RandomColor()}
	}
	return res
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
