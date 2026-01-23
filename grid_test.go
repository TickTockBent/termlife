package main

import (
	"testing"
)

func TestNewGrid(t *testing.T) {
	grid := NewGrid(10, 20, false)

	if grid.Width != 10 {
		t.Errorf("expected width 10, got %d", grid.Width)
	}
	if grid.Height != 20 {
		t.Errorf("expected height 20, got %d", grid.Height)
	}
	if grid.Wrap {
		t.Error("expected wrap to be false")
	}
	if len(grid.Cells) != 20 {
		t.Errorf("expected 20 rows, got %d", len(grid.Cells))
	}
	if len(grid.Cells[0]) != 10 {
		t.Errorf("expected 10 columns, got %d", len(grid.Cells[0]))
	}
}

func TestGridGetSet(t *testing.T) {
	grid := NewGrid(10, 10, false)

	// Initially all dead
	if grid.Get(5, 5) {
		t.Error("expected cell to be dead initially")
	}

	// Set a cell alive
	grid.Set(5, 5, true)
	if !grid.Get(5, 5) {
		t.Error("expected cell to be alive after Set")
	}

	// Set it back to dead
	grid.Set(5, 5, false)
	if grid.Get(5, 5) {
		t.Error("expected cell to be dead after Set(false)")
	}
}

func TestGridGetOutOfBounds(t *testing.T) {
	grid := NewGrid(10, 10, false)
	grid.Set(0, 0, true)

	// Out of bounds should return false when wrap is disabled
	if grid.Get(-1, 0) {
		t.Error("expected false for negative x")
	}
	if grid.Get(0, -1) {
		t.Error("expected false for negative y")
	}
	if grid.Get(10, 0) {
		t.Error("expected false for x >= width")
	}
	if grid.Get(0, 10) {
		t.Error("expected false for y >= height")
	}
}

func TestGridWrapping(t *testing.T) {
	grid := NewGrid(10, 10, true)
	grid.Set(0, 0, true)

	// With wrapping, negative coords should wrap around
	if !grid.Get(-10, 0) {
		t.Error("expected wrap to work for x=-10")
	}
	if !grid.Get(10, 0) {
		t.Error("expected wrap to work for x=10")
	}
	if !grid.Get(0, -10) {
		t.Error("expected wrap to work for y=-10")
	}
	if !grid.Get(0, 10) {
		t.Error("expected wrap to work for y=10")
	}
}

func TestGridToggle(t *testing.T) {
	grid := NewGrid(10, 10, false)

	grid.Toggle(5, 5)
	if !grid.Get(5, 5) {
		t.Error("expected cell to be alive after toggle")
	}

	grid.Toggle(5, 5)
	if grid.Get(5, 5) {
		t.Error("expected cell to be dead after second toggle")
	}
}

func TestCountNeighbors(t *testing.T) {
	grid := NewGrid(10, 10, false)

	// Set up a blinker pattern
	grid.Set(4, 5, true)
	grid.Set(5, 5, true)
	grid.Set(6, 5, true)

	// Center cell should have 2 neighbors
	if n := grid.CountNeighbors(5, 5); n != 2 {
		t.Errorf("expected 2 neighbors for center, got %d", n)
	}

	// Cell above center should have 3 neighbors
	if n := grid.CountNeighbors(5, 4); n != 3 {
		t.Errorf("expected 3 neighbors for cell above center, got %d", n)
	}

	// End cell should have 1 neighbor
	if n := grid.CountNeighbors(4, 5); n != 1 {
		t.Errorf("expected 1 neighbor for end cell, got %d", n)
	}
}

func TestCountNeighborsWithWrap(t *testing.T) {
	grid := NewGrid(5, 5, true)

	// Cell at corner with neighbor wrapped around
	grid.Set(0, 0, true)
	grid.Set(4, 4, true) // Should be neighbor of (0,0) with wrapping

	if n := grid.CountNeighbors(0, 0); n != 1 {
		t.Errorf("expected 1 neighbor with wrapping, got %d", n)
	}
}

func TestGridStep(t *testing.T) {
	grid := NewGrid(10, 10, false)
	rule := ConwayRule()

	// Set up a blinker (horizontal)
	grid.Set(4, 5, true)
	grid.Set(5, 5, true)
	grid.Set(6, 5, true)

	// After one step, should become vertical
	grid.Step(rule)

	// Check vertical configuration
	if !grid.Get(5, 4) {
		t.Error("expected (5,4) to be alive")
	}
	if !grid.Get(5, 5) {
		t.Error("expected (5,5) to be alive")
	}
	if !grid.Get(5, 6) {
		t.Error("expected (5,6) to be alive")
	}

	// Check horizontal cells are now dead
	if grid.Get(4, 5) {
		t.Error("expected (4,5) to be dead")
	}
	if grid.Get(6, 5) {
		t.Error("expected (6,5) to be dead")
	}
}

func TestGridClear(t *testing.T) {
	grid := NewGrid(10, 10, false)

	// Set some cells
	grid.Set(1, 1, true)
	grid.Set(5, 5, true)
	grid.Set(9, 9, true)

	grid.Clear()

	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			if grid.Get(x, y) {
				t.Errorf("expected all cells dead after clear, found alive at (%d,%d)", x, y)
			}
		}
	}
}

func TestGridRandomize(t *testing.T) {
	grid := NewGrid(100, 100, false)

	// Use a fixed "random" function for testing
	callCount := 0
	fixedRng := func() float64 {
		callCount++
		if callCount%2 == 0 {
			return 0.1 // Below density
		}
		return 0.9 // Above density
	}

	grid.Randomize(0.5, fixedRng)

	// Should have approximately half cells alive
	alive := grid.CountAlive()
	if alive == 0 || alive == 10000 {
		t.Error("randomize should produce a mix of alive and dead cells")
	}
}

func TestCountAlive(t *testing.T) {
	grid := NewGrid(10, 10, false)

	if count := grid.CountAlive(); count != 0 {
		t.Errorf("expected 0 alive cells initially, got %d", count)
	}

	grid.Set(0, 0, true)
	grid.Set(5, 5, true)
	grid.Set(9, 9, true)

	if count := grid.CountAlive(); count != 3 {
		t.Errorf("expected 3 alive cells, got %d", count)
	}
}
