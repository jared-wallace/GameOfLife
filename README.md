# Conway's Game of Life in Go

Conway's Game of Life in Go is a high-performance, feature-rich implementation of John Conway's renowned cellular automaton. Built with efficiency and flexibility in mind, this application allows users to explore intricate patterns and observe their evolution over time.

## Features

    Efficient Simulation: Utilizes concurrent processing to handle large patterns smoothly.
    Flexible Grid Management: Supports dynamic resizing of the grid with adjustable cell sizes.
    Multiple Pattern Formats: Load patterns from .txt, .rle, and .mc files.
    Interactive Controls: Easily adjust simulation speed, cell size, and switch between patterns.
    User-Friendly Interface: Built with Ebiten, providing a responsive and intuitive GUI.

## Supported Pattern Formats

Plaintext (.txt): Simple format where each line represents a row of cells. Live cells are denoted by 'O', and dead cells by '.'.

Example:

```
# This is a comment
..O..
.O.O.
OOO..
```

Run-Length Encoded (.rle): Compact representation using counts and symbols to denote sequences of live or dead cells.

Example:

```
x = 3, y = 3, rule = B3/S23
bo$2bo$3o!
```

Macrocell (.mc): Advanced format used by Golly, supporting quadtree-based representations for efficient storage and manipulation of large universes.

Example:

```
[M2] (golly 2.0)
#R B3/S23
#G 0
$$..*$...*$.***$$$$
```

## Installation

### Prerequisites

Go Programming Language: Ensure you have Go installed. You can download it from golang.org.

### Steps

Clone the Repository:

`git clone https://github.com/jared-wallace/gol.git`

Build the Application:

`go build -o gameoflife`

Run the Application:

`./gameoflife`

## Usage

Upon launching the application, you'll be greeted with a window displaying the cellular grid.

### Controls

  - Spacebar: Cycle through available patterns.
  - '+' / '-': Increase or decrease the cell size for better visibility.
  - Up Arrow: Increase the simulation tick speed (TPS - Ticks Per Second).
  - Down Arrow: Decrease the simulation tick speed.
  - Escape: Exit the application.
