package main

// Pattern represents a preset pattern as a list of relative coordinates
type Pattern struct {
	Name   string
	Coords [][2]int
}

// Patterns contains all available preset patterns
var Patterns = map[string]*Pattern{
	"glider": {
		Name: "Glider",
		Coords: [][2]int{
			{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2},
		},
	},
	"blinker": {
		Name: "Blinker",
		Coords: [][2]int{
			{0, 0}, {1, 0}, {2, 0},
		},
	},
	"toad": {
		Name: "Toad",
		Coords: [][2]int{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {1, 1}, {2, 1},
		},
	},
	"beacon": {
		Name: "Beacon",
		Coords: [][2]int{
			{0, 0}, {1, 0},
			{0, 1}, {1, 1},
			{2, 2}, {3, 2},
			{2, 3}, {3, 3},
		},
	},
	"pulsar": {
		Name: "Pulsar",
		Coords: [][2]int{
			// Top section
			{2, 0}, {3, 0}, {4, 0}, {8, 0}, {9, 0}, {10, 0},
			{0, 2}, {5, 2}, {7, 2}, {12, 2},
			{0, 3}, {5, 3}, {7, 3}, {12, 3},
			{0, 4}, {5, 4}, {7, 4}, {12, 4},
			{2, 5}, {3, 5}, {4, 5}, {8, 5}, {9, 5}, {10, 5},
			// Bottom section (mirror)
			{2, 7}, {3, 7}, {4, 7}, {8, 7}, {9, 7}, {10, 7},
			{0, 8}, {5, 8}, {7, 8}, {12, 8},
			{0, 9}, {5, 9}, {7, 9}, {12, 9},
			{0, 10}, {5, 10}, {7, 10}, {12, 10},
			{2, 12}, {3, 12}, {4, 12}, {8, 12}, {9, 12}, {10, 12},
		},
	},
	"gosper-gun": {
		Name: "Gosper Glider Gun",
		Coords: [][2]int{
			// Left square
			{0, 4}, {0, 5}, {1, 4}, {1, 5},
			// Left part
			{10, 4}, {10, 5}, {10, 6},
			{11, 3}, {11, 7},
			{12, 2}, {12, 8},
			{13, 2}, {13, 8},
			{14, 5},
			{15, 3}, {15, 7},
			{16, 4}, {16, 5}, {16, 6},
			{17, 5},
			// Right part
			{20, 2}, {20, 3}, {20, 4},
			{21, 2}, {21, 3}, {21, 4},
			{22, 1}, {22, 5},
			{24, 0}, {24, 1}, {24, 5}, {24, 6},
			// Right square
			{34, 2}, {34, 3}, {35, 2}, {35, 3},
		},
	},
	"diehard": {
		Name: "Diehard",
		Coords: [][2]int{
			{6, 0},
			{0, 1}, {1, 1},
			{1, 2}, {5, 2}, {6, 2}, {7, 2},
		},
	},
	"acorn": {
		Name: "Acorn",
		Coords: [][2]int{
			{1, 0},
			{3, 1},
			{0, 2}, {1, 2}, {4, 2}, {5, 2}, {6, 2},
		},
	},
	"r-pentomino": {
		Name: "R-pentomino",
		Coords: [][2]int{
			{1, 0}, {2, 0},
			{0, 1}, {1, 1},
			{1, 2},
		},
	},
	"lwss": {
		Name: "Lightweight Spaceship",
		Coords: [][2]int{
			{1, 0}, {4, 0},
			{0, 1},
			{0, 2}, {4, 2},
			{0, 3}, {1, 3}, {2, 3}, {3, 3},
		},
	},
	"block": {
		Name: "Block",
		Coords: [][2]int{
			{0, 0}, {1, 0},
			{0, 1}, {1, 1},
		},
	},
	"beehive": {
		Name: "Beehive",
		Coords: [][2]int{
			{1, 0}, {2, 0},
			{0, 1}, {3, 1},
			{1, 2}, {2, 2},
		},
	},
	"loaf": {
		Name: "Loaf",
		Coords: [][2]int{
			{1, 0}, {2, 0},
			{0, 1}, {3, 1},
			{1, 2}, {3, 2},
			{2, 3},
		},
	},
	"pentadecathlon": {
		Name: "Pentadecathlon",
		Coords: [][2]int{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{0, 2}, {4, 2},
			{1, 3}, {2, 3}, {3, 3},
			{1, 6}, {2, 6}, {3, 6},
			{0, 7}, {4, 7},
			{0, 8}, {4, 8},
			{1, 9}, {2, 9}, {3, 9},
		},
	},
}

// PatternNames returns a sorted list of available pattern names
func PatternNames() []string {
	return []string{
		"random", "glider", "blinker", "toad", "beacon", "pulsar",
		"gosper-gun", "diehard", "acorn", "r-pentomino", "lwss",
		"block", "beehive", "loaf", "pentadecathlon",
	}
}

// PlacePattern places a pattern on the grid, centered at the given position
func PlacePattern(grid *Grid, pattern *Pattern, centerX, centerY int) {
	if pattern == nil {
		return
	}

	// Find pattern bounds to center it properly
	minX, minY, maxX, maxY := patternBounds(pattern)
	patternWidth := maxX - minX + 1
	patternHeight := maxY - minY + 1

	// Calculate offset to center the pattern
	offsetX := centerX - patternWidth/2 - minX
	offsetY := centerY - patternHeight/2 - minY

	for _, coord := range pattern.Coords {
		x := coord[0] + offsetX
		y := coord[1] + offsetY
		grid.Set(x, y, true)
	}
}

// patternBounds returns the bounding box of a pattern
func patternBounds(pattern *Pattern) (minX, minY, maxX, maxY int) {
	if len(pattern.Coords) == 0 {
		return 0, 0, 0, 0
	}

	minX = pattern.Coords[0][0]
	maxX = pattern.Coords[0][0]
	minY = pattern.Coords[0][1]
	maxY = pattern.Coords[0][1]

	for _, coord := range pattern.Coords {
		if coord[0] < minX {
			minX = coord[0]
		}
		if coord[0] > maxX {
			maxX = coord[0]
		}
		if coord[1] < minY {
			minY = coord[1]
		}
		if coord[1] > maxY {
			maxY = coord[1]
		}
	}
	return
}
