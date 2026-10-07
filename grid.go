package main

// Grid holds the cellular automaton state
type Grid struct {
	Width  int
	Height int
	Cells  [][]bool
	Wrap   bool
}

// NewGrid creates a new grid with the given dimensions
func NewGrid(width, height int, wrap bool) *Grid {
	cells := make([][]bool, height)
	for y := range cells {
		cells[y] = make([]bool, width)
	}
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  cells,
		Wrap:   wrap,
	}
}

// Get returns the cell state at (x, y), handling wrapping if enabled
func (g *Grid) Get(x, y int) bool {
	if g.Wrap {
		x = ((x % g.Width) + g.Width) % g.Width
		y = ((y % g.Height) + g.Height) % g.Height
	} else if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return false
	}
	return g.Cells[y][x]
}

// Set sets the cell state at (x, y)
func (g *Grid) Set(x, y int, alive bool) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		g.Cells[y][x] = alive
	}
}

// Toggle flips the cell state at (x, y)
func (g *Grid) Toggle(x, y int) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		g.Cells[y][x] = !g.Cells[y][x]
	}
}

// CountNeighbors returns the number of alive neighbors for cell at (x, y)
func (g *Grid) CountNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if g.Get(x+dx, y+dy) {
				count++
			}
		}
	}
	return count
}

// Step advances the grid by one generation using the given rule
func (g *Grid) Step(rule *Rule) {
	// Create a new grid to hold the next state
	next := make([][]bool, g.Height)
	for y := range next {
		next[y] = make([]bool, g.Width)
	}

	// Apply rules to each cell
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			neighbors := g.CountNeighbors(x, y)
			alive := g.Cells[y][x]

			if alive {
				next[y][x] = rule.Survives(neighbors)
			} else {
				next[y][x] = rule.Births(neighbors)
			}
		}
	}

	g.Cells = next
}

// Clear sets all cells to dead
func (g *Grid) Clear() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.Cells[y][x] = false
		}
	}
}

// Randomize fills the grid randomly with the given density (0.0-1.0)
func (g *Grid) Randomize(density float64, rng func() float64) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.Cells[y][x] = rng() < density
		}
	}
}

// CountAlive returns the total number of alive cells
func (g *Grid) CountAlive() int {
	count := 0
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Cells[y][x] {
				count++
			}
		}
	}
	return count
}

// Resize changes the grid dimensions, keeping cells in the overlapping top-left region
func (g *Grid) Resize(width, height int) {
	resized := NewGrid(width, height, g.Wrap)
	for y := 0; y < min(g.Height, height); y++ {
		copy(resized.Cells[y], g.Cells[y][:min(g.Width, width)])
	}
	*g = *resized
}
