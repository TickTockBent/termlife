package main

import (
	"testing"
)

func TestPatternsDefined(t *testing.T) {
	// Verify all expected patterns exist
	expectedPatterns := []string{
		"glider", "blinker", "toad", "beacon", "pulsar",
		"gosper-gun", "diehard", "acorn", "r-pentomino",
		"lwss", "block", "beehive", "loaf", "pentadecathlon",
	}

	for _, name := range expectedPatterns {
		if _, ok := Patterns[name]; !ok {
			t.Errorf("expected pattern %q to be defined", name)
		}
	}
}

func TestPatternNames(t *testing.T) {
	names := PatternNames()

	// Should include "random" as first option
	if names[0] != "random" {
		t.Errorf("expected first pattern name to be 'random', got %q", names[0])
	}

	// Should have more than just "random"
	if len(names) < 5 {
		t.Errorf("expected at least 5 pattern names, got %d", len(names))
	}
}

func TestPlacePattern(t *testing.T) {
	grid := NewGrid(20, 20, false)
	pattern := Patterns["glider"]

	PlacePattern(grid, pattern, 10, 10)

	// Count alive cells - glider has 5 cells
	alive := grid.CountAlive()
	if alive != 5 {
		t.Errorf("expected 5 alive cells for glider, got %d", alive)
	}
}

func TestPlacePatternCentered(t *testing.T) {
	grid := NewGrid(20, 20, false)
	pattern := Patterns["block"] // 2x2 pattern

	PlacePattern(grid, pattern, 10, 10)

	// Block should be roughly centered at (10, 10)
	// Check that cells exist near the center
	foundNearCenter := false
	for y := 8; y <= 12; y++ {
		for x := 8; x <= 12; x++ {
			if grid.Get(x, y) {
				foundNearCenter = true
				break
			}
		}
	}

	if !foundNearCenter {
		t.Error("expected pattern to be placed near center")
	}
}

func TestPlacePatternNil(t *testing.T) {
	grid := NewGrid(20, 20, false)

	// Should not panic with nil pattern
	PlacePattern(grid, nil, 10, 10)

	if grid.CountAlive() != 0 {
		t.Error("expected no cells alive after placing nil pattern")
	}
}

func TestPatternBounds(t *testing.T) {
	testCases := []struct {
		name     string
		expected struct{ minX, minY, maxX, maxY int }
	}{
		{
			"glider",
			struct{ minX, minY, maxX, maxY int }{0, 0, 2, 2},
		},
		{
			"blinker",
			struct{ minX, minY, maxX, maxY int }{0, 0, 2, 0},
		},
		{
			"block",
			struct{ minX, minY, maxX, maxY int }{0, 0, 1, 1},
		},
	}

	for _, tc := range testCases {
		pattern := Patterns[tc.name]
		minX, minY, maxX, maxY := patternBounds(pattern)

		if minX != tc.expected.minX || minY != tc.expected.minY ||
			maxX != tc.expected.maxX || maxY != tc.expected.maxY {
			t.Errorf("%s bounds: expected (%d,%d)-(%d,%d), got (%d,%d)-(%d,%d)",
				tc.name,
				tc.expected.minX, tc.expected.minY, tc.expected.maxX, tc.expected.maxY,
				minX, minY, maxX, maxY)
		}
	}
}

func TestGliderPattern(t *testing.T) {
	grid := NewGrid(10, 10, false)
	rule := ConwayRule()

	PlacePattern(grid, Patterns["glider"], 3, 3)

	// Count initial cells
	initialCount := grid.CountAlive()
	if initialCount != 5 {
		t.Errorf("expected 5 initial cells, got %d", initialCount)
	}

	// After 4 generations, glider should have moved but still have 5 cells
	for i := 0; i < 4; i++ {
		grid.Step(rule)
	}

	finalCount := grid.CountAlive()
	if finalCount != 5 {
		t.Errorf("expected 5 cells after 4 generations, got %d", finalCount)
	}
}

func TestBlinkerPattern(t *testing.T) {
	grid := NewGrid(10, 10, false)
	rule := ConwayRule()

	PlacePattern(grid, Patterns["blinker"], 5, 5)

	// Blinker has period 2
	// After 2 steps, should be back to original configuration

	// Capture original state (find which cells are alive)
	original := make(map[[2]int]bool)
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			if grid.Get(x, y) {
				original[[2]int{x, y}] = true
			}
		}
	}

	grid.Step(rule)
	grid.Step(rule)

	// Compare to original
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			expected := original[[2]int{x, y}]
			actual := grid.Get(x, y)
			if expected != actual {
				t.Errorf("blinker period-2 failed at (%d,%d): expected %v, got %v", x, y, expected, actual)
			}
		}
	}
}
