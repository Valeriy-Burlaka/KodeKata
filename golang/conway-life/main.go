package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

const GRID_MIN_WIDTH = 5
const GRID_MIN_HEIGHT = 5
const GRID_MAX_WIDTH = 1<<16 - 1
const GRID_MAX_HEIGHT = 1<<16 - 1


type Pattern string

func (p Pattern) Parse() [][]bool {
	parsed := slices.DeleteFunc(
		strings.Split(string(p), "\n"),
		func (s string) bool {
			return s == ""
		},
	)
	for i, _ := range parsed {
		if len(parsed[i]) != len(parsed[0]) {
			log.Fatalf("Invalid pattern:\n%s\nAll rows must have the same length = %d but row %d has length = %d",
				p, len(parsed[0]), i, len(parsed[i]))
		}
	}
	
	result := make([][]bool, len(parsed))

	for i, row := range parsed {
		result[i] = make([]bool, len(row))
		for j, sym := range row {
			result[i][j] = string(sym) == "X"
		}
	}

	return result
}

var Kickback Pattern = `
-X-
X--
XXX
`

type Cell struct {
	// y uint16
	// x uint16
	isAlive bool
	neighbors []*Cell
}

// func NewCell(y, x uint16) Cell {
// 	c := Cell{y: y, x: x, neighbors: make([]*Cell, 0, 8)}

// 	return c
// }


// func NewCell() Cell {
// 	c := Cell{neighbors: make([]*Cell, 0, 8)}

// 	return c
// }

func (c *Cell) AddNeighbor (n *Cell) {
	c.neighbors = append(c.neighbors, n)
}

func (c *Cell) AddNeighbors (nn... *Cell) {
	for _, n := range nn {
		c.neighbors = append(c.neighbors, n)
	}
}

type Grid struct {
	width uint16
	height uint16
	cells []Cell
	cellIndex map[uint16][]*Cell
}

func NewGrid(width, height uint16) (*Grid, error) {
	if width > GRID_MAX_WIDTH {
		return nil, fmt.Errorf("failed to create new grid, width %d is too big (max width=%d)", width, GRID_MAX_WIDTH)
	} else if width < GRID_MIN_WIDTH {
		return nil, fmt.Errorf("failed to create new grid, width %d is too small (min width=%d)", width, GRID_MIN_WIDTH)
	}
	if height > GRID_MAX_HEIGHT {
		return nil, fmt.Errorf("failed to create new grid, height %d is too big (max height=%d)", height, GRID_MIN_HEIGHT)
	} else if height < GRID_MIN_HEIGHT {
		return nil, fmt.Errorf("failed to create new grid, height %d is too small (min height=%d)", width, GRID_MIN_HEIGHT)
	}

	g := Grid{
		width: width,
		height: height,
		cells: make([]Cell, width*height),
		cellIndex: make(map[uint16][]*Cell, height),
	}

	// Create grid cells.
	var y uint16
	var x uint16
	for y = 0; y < height; y++ {
		row := make([]*Cell, width)
		for x = 0; x < width; x++ {
			c := Cell{}
			g.cells[y*width + x] = c
			row[x] = &c
		}
		g.cellIndex[y] = row
	}

	// Calculate and populate each cell's neighbors.
	var i uint16
	max := uint16(len(g.cells))
	for i = 0; i < max; i++ {
		y := i / height
		x := i % width

		switch y {
		// top row
		case 0:
			switch x {
			case 0:
				// top left corner
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y][x+1],
					g.cellIndex[y+1][x],
					g.cellIndex[y+1][x+1],
				}
			case width-1:
				// top right corner
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y][x-1],
					g.cellIndex[y+1][x],
					g.cellIndex[y+1][x-1],
				}
			default:
				// non-corner cell in the top row
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y][x-1],
					g.cellIndex[y][x+1],
					g.cellIndex[y+1][x-1],
					g.cellIndex[y+1][x],
					g.cellIndex[y+1][x+1],
				}
			}
		case height-1:
			// bottom row
			switch x {
			case 0:
				// bottom left corner
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x+1],
					g.cellIndex[y][x+1],
				}
			case width-1:
				// bottom right corner
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x-1],
					g.cellIndex[y][x-1],
				}
			default:
				// non-corner cell in the bottom row
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y][x-1],
					g.cellIndex[y][x+1],
					g.cellIndex[y-1][x-1],
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x+1],
				}
			}
		default:
			// middle rows
			switch x {
			case 0:
				// left-edge column
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x+1],
					g.cellIndex[y][x+1],
					g.cellIndex[y+1][x+1],
					g.cellIndex[y+1][x],
				}
			case width-1:
				// right-edge column
				g.cells[i].neighbors = []*Cell{
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x-1],
					g.cellIndex[y][x-1],
					g.cellIndex[y+1][x-1],
					g.cellIndex[y+1][x],
				}
			default:
				// middle cells (each has 8 neighbors)
				g.cells[i].neighbors = []*Cell{
					// from top-left neighbor, clockwise
					g.cellIndex[y-1][x-1],
					g.cellIndex[y-1][x],
					g.cellIndex[y-1][x+1],
					g.cellIndex[y][x+1],
					g.cellIndex[y+1][x+1],
					g.cellIndex[y+1][x],
					g.cellIndex[y+1][x-1],
					g.cellIndex[y][x-1],
				}
			}
		}
	}

	// Populate living cells from the seed pattern

	return &g, nil
}

// Grid.Render()
// Grid.Evolve()

func main() {

	// Init a Grid
	// g, err := NewGrid(5, 5)
	// if err != nil {
	// 	log.Fatalf("failed to run: %v", err)
	// }

	// fmt.Println(g)

	p := Kickback.Parse()
	fmt.Println(p)
}
