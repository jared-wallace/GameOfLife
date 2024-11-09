package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jared-wallace/gol/engine"
)

// main initializes and runs the game.
func main() {
	// Initial grid size
	initialGridWidth, initialGridHeight := 320, 256 // Adjusted for better performance

	// Create a new game instance
	game := engine.NewGame(initialGridWidth, initialGridHeight)

	// Configure Ebiten window
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Set initial window size based on grid size and cell size
	ebiten.SetWindowSize(initialGridWidth*game.GetCellSize(), initialGridHeight*game.GetCellSize())
	ebiten.SetWindowTitle("Conway's Game of Life")
	ebiten.SetTPS(30)
	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
