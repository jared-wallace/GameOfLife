package display

import (
	"github.com/faiface/pixel"
	"math/rand"
)

func RandomColor() pixel.RGBA {
	return pixel.RGBA{
		R: float64(rand.Float32()),
		G: float64(rand.Float32()),
		B: float64(rand.Float32()),
		A: 1,
	}
}
