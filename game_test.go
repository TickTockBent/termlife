package main

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func newTestGame(width, height int, themeName string) *Game {
	return NewGame(NewGrid(width, height, false), ConwayRule(), Themes[themeName], 10)
}

func TestSetupGrid(t *testing.T) {
	grid, err := SetupGrid(20, 20, false, "glider", 0.25, func() float64 { return 1 })
	if err != nil {
		t.Fatalf("SetupGrid(glider) error = %v", err)
	}
	if grid.CountAlive() != len(Patterns["glider"].Coords) {
		t.Errorf("glider grid has %d cells, want %d", grid.CountAlive(), len(Patterns["glider"].Coords))
	}

	grid, err = SetupGrid(20, 20, true, "random", 0.5, func() float64 { return 0 })
	if err != nil {
		t.Fatalf("SetupGrid(random) error = %v", err)
	}
	if grid.CountAlive() != 400 || !grid.Wrap {
		t.Errorf("random grid: alive = %d, wrap = %v; want 400, true", grid.CountAlive(), grid.Wrap)
	}

	if _, err := SetupGrid(20, 20, false, "nope", 0.25, nil); err == nil {
		t.Error("SetupGrid(unknown pattern) should return an error")
	}
}

func TestNewGameAges(t *testing.T) {
	if game := newTestGame(10, 10, "white"); game.Ages != nil {
		t.Error("non-rainbow theme should not track ages")
	}
	game := newTestGame(12, 10, "rainbow")
	if len(game.Ages) != 10 || len(game.Ages[0]) != 12 {
		t.Errorf("rainbow ages sized %dx%d, want 12x10", len(game.Ages[0]), len(game.Ages))
	}
	if game.CursorX != 6 || game.CursorY != 5 {
		t.Errorf("cursor at (%d,%d), want centered (6,5)", game.CursorX, game.CursorY)
	}
}

func TestStepWithAges(t *testing.T) {
	game := newTestGame(10, 10, "rainbow")
	// Block (still life) plus a lone cell that dies
	for _, coord := range [][2]int{{1, 1}, {2, 1}, {1, 2}, {2, 2}, {7, 7}} {
		game.SetCell(coord[0], coord[1], true)
	}

	stepWithAges(game)
	stepWithAges(game)

	if game.Generation != 2 {
		t.Errorf("Generation = %d, want 2", game.Generation)
	}
	if game.Ages[1][1] != 2 {
		t.Errorf("surviving block cell age = %d, want 2", game.Ages[1][1])
	}
	if game.Grid.Get(7, 7) || game.Ages[7][7] != 0 {
		t.Errorf("lone cell should be dead with age 0, got alive=%v age=%d", game.Grid.Get(7, 7), game.Ages[7][7])
	}
}

func TestGameResize(t *testing.T) {
	game := newTestGame(20, 20, "rainbow")
	game.SetCell(2, 3, true)
	game.SetCell(15, 15, true)
	game.Ages[3][2] = 5
	game.CursorX, game.CursorY = 18, 18

	game.Resize(10, 12)

	if game.Grid.Width != 10 || game.Grid.Height != 12 {
		t.Fatalf("grid is %dx%d, want 10x12", game.Grid.Width, game.Grid.Height)
	}
	if !game.Grid.Get(2, 3) || game.Ages[3][2] != 5 {
		t.Error("cell and age inside the kept region should survive resize")
	}
	if game.Grid.CountAlive() != 1 {
		t.Errorf("alive = %d, want 1 (cell outside new bounds dropped)", game.Grid.CountAlive())
	}
	if len(game.Ages) != 12 || len(game.Ages[0]) != 10 {
		t.Errorf("ages sized %dx%d, want 10x12", len(game.Ages[0]), len(game.Ages))
	}
	if game.CursorX != 9 || game.CursorY != 11 {
		t.Errorf("cursor at (%d,%d), want clamped to (9,11)", game.CursorX, game.CursorY)
	}

	game.Resize(30, 30)
	if !game.Grid.Get(2, 3) || game.Grid.Width != 30 {
		t.Error("growing the grid should keep existing cells")
	}
}

func TestMousePaintAndErase(t *testing.T) {
	game := newTestGame(10, 10, "white")

	// Screen column 7 is the second half of cell 3
	action := HandleInput(tcell.NewEventMouse(7, 4, tcell.ButtonPrimary, tcell.ModNone), game)
	if action != ActionPaintCell || game.CursorX != 3 || game.CursorY != 4 {
		t.Fatalf("left click: action=%v cursor=(%d,%d), want paint at (3,4)", action, game.CursorX, game.CursorY)
	}
	ApplyAction(action, game, 0, nil)
	if !game.Grid.Get(3, 4) {
		t.Error("left click should make the cell alive")
	}

	action = HandleInput(tcell.NewEventMouse(7, 4, tcell.ButtonSecondary, tcell.ModNone), game)
	ApplyAction(action, game, 0, nil)
	if action != ActionEraseCell || game.Grid.Get(3, 4) {
		t.Error("right click should make the cell dead")
	}

	if action := HandleInput(tcell.NewEventMouse(7, 4, tcell.ButtonNone, tcell.ModNone), game); action != ActionNone {
		t.Errorf("motion without buttons should be ignored, got %v", action)
	}
	if action := HandleInput(tcell.NewEventMouse(5, 10, tcell.ButtonPrimary, tcell.ModNone), game); action != ActionNone {
		t.Errorf("click below the grid (status bar) should be ignored, got %v", action)
	}
}

func TestToggleCellRequiresPause(t *testing.T) {
	game := newTestGame(10, 10, "rainbow")
	ApplyAction(ActionToggleCell, game, 0, nil)
	if game.Grid.Get(5, 5) {
		t.Error("toggle while running should do nothing")
	}
	game.Paused = true
	ApplyAction(ActionToggleCell, game, 0, nil)
	if !game.Grid.Get(5, 5) {
		t.Error("toggle while paused should flip the cell")
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input      string
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{"80x40", 80, 40, false},
		{"80X40", 80, 40, false},
		{"10x10", 10, 10, false},
		{"9x40", 0, 0, true},
		{"80", 0, 0, true},
		{"axb", 0, 0, true},
		{"80x40x2", 0, 0, true},
	}
	for _, tt := range tests {
		width, height, err := parseSize(tt.input)
		if (err != nil) != tt.wantErr || width != tt.wantWidth || height != tt.wantHeight {
			t.Errorf("parseSize(%q) = %d, %d, %v; want %d, %d, err=%v",
				tt.input, width, height, err, tt.wantWidth, tt.wantHeight, tt.wantErr)
		}
	}
}
