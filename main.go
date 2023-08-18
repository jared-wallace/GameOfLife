package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.expedia.biz/jarwallace/gol/internal/display"
	"github.expedia.biz/jarwallace/gol/internal/models"
	"github.expedia.biz/jarwallace/gol/internal/processor"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	_ "image/png"
	"io"
	"math/rand"
	"os"
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
	face, err := loadTTF("Consolas.ttf", 14)
	if err != nil {
		panic(err)
	}
	atlas := text.NewAtlas(face, text.ASCII)

	currGen := models.InjectRPentomino()
	pr := processor.NewProcessor(winXMax, winYMax)
	clicks := 0
	generation := 0
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			switch clicks % 4 {
			case 0:
				currGen = models.InjectGlider()
			case 1:
				currGen = models.InjectAcorn()
			case 2:
				currGen = getInitial(winXMax, winYMax)
			case 3:
				currGen = models.InjectRPentomino()
			}
			generation = 0
			clicks++
		}
		win.Clear(colornames.Black)
		drawPop(atlas, win, currGen)
		time.Sleep(REFRESH * time.Millisecond) // Delay to observe the generation

		if !popExists(currGen) {
			break // Exit the loop if there are no alive cells
		}
		currGen = pr.AnalyzePopEfficiently(currGen)
		txt := text.New(pixel.V(10, l-40), atlas) // This places the text near the top-left corner. Adjust as necessary.
		fmt.Fprintf(txt, "Generation: %d", generation)
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 1))
		generation++
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
		t.Draw(win, pixel.IM)
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

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ttfBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(ttfBytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size: size,
	}), nil
}
