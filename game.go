package main

import "fmt"

// Game holds all runtime state
type Game struct {
	Grid       *Grid
	Rule       *Rule
	Theme      *Theme
	FPS        int
	Paused     bool
	Generation int
	CursorX    int
	CursorY    int
	Ages       [][]int // For rainbow mode age tracking
	AutoSize   bool    // Grid follows the terminal size on resize
}

// NewGame creates a game around an initialized grid
func NewGame(grid *Grid, rule *Rule, theme *Theme, fps int) *Game {
	game := &Game{
		Grid:    grid,
		Rule:    rule,
		Theme:   theme,
		FPS:     fps,
		CursorX: grid.Width / 2,
		CursorY: grid.Height / 2,
	}
	if theme.AgeColored {
		game.Ages = newAges(grid.Width, grid.Height)
	}
	return game
}

// SetupGrid creates a grid and fills it with the named pattern or random cells
func SetupGrid(width, height int, wrap bool, patternName string, density float64, rng func() float64) (*Grid, error) {
	grid := NewGrid(width, height, wrap)
	if patternName == "random" {
		grid.Randomize(density, rng)
		return grid, nil
	}
	pattern, ok := Patterns[patternName]
	if !ok {
		return nil, fmt.Errorf("unknown pattern: %s", patternName)
	}
	PlacePattern(grid, pattern, width/2, height/2)
	return grid, nil
}

// Resize changes the grid dimensions, keeping cells and ages in the overlapping region
func (game *Game) Resize(width, height int) {
	if width == game.Grid.Width && height == game.Grid.Height {
		return
	}
	oldWidth, oldHeight := game.Grid.Width, game.Grid.Height
	game.Grid.Resize(width, height)

	if game.Ages != nil {
		resizedAges := newAges(width, height)
		for y := 0; y < min(oldHeight, height); y++ {
			copy(resizedAges[y], game.Ages[y][:min(oldWidth, width)])
		}
		game.Ages = resizedAges
	}

	game.CursorX = min(game.CursorX, width-1)
	game.CursorY = min(game.CursorY, height-1)
}

// SetCell forces the cell at (x, y) alive or dead, resetting its age
func (game *Game) SetCell(x, y int, alive bool) {
	game.Grid.Set(x, y, alive)
	if game.Ages != nil && x >= 0 && x < game.Grid.Width && y >= 0 && y < game.Grid.Height {
		game.Ages[y][x] = 0
	}
}

// stepWithAges advances the simulation and updates cell ages
func stepWithAges(game *Game) {
	// Step swaps in a fresh cell slice, so the previous generation stays intact
	previousCells := game.Grid.Cells
	game.Grid.Step(game.Rule)
	game.Generation++

	if game.Ages == nil {
		return
	}

	for y := 0; y < game.Grid.Height; y++ {
		for x := 0; x < game.Grid.Width; x++ {
			if game.Grid.Cells[y][x] && previousCells[y][x] {
				game.Ages[y][x]++ // Survived
			} else {
				game.Ages[y][x] = 0 // Born or dead
			}
		}
	}
}

// resetAges resets all cell ages to 0
func resetAges(game *Game) {
	if game.Ages == nil {
		return
	}
	for y := range game.Ages {
		for x := range game.Ages[y] {
			game.Ages[y][x] = 0
		}
	}
}

func newAges(width, height int) [][]int {
	ages := make([][]int, height)
	for y := range ages {
		ages[y] = make([]int, width)
	}
	return ages
}
