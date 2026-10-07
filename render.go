package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Renderer handles all screen drawing operations
type Renderer struct {
	Screen tcell.Screen
	Theme  *Theme
}

// NewRenderer creates a new renderer with the given screen and theme
func NewRenderer(screen tcell.Screen, theme *Theme) *Renderer {
	return &Renderer{
		Screen: screen,
		Theme:  theme,
	}
}

// Clear clears the screen with the background color
func (r *Renderer) Clear() {
	r.Screen.Clear()
}

// DrawGrid renders the grid cells to the screen
// Each cell is rendered as two characters for a more square aspect ratio
func (r *Renderer) DrawGrid(game *Game) {
	for y := 0; y < game.Grid.Height; y++ {
		for x := 0; x < game.Grid.Width; x++ {
			screenX := x * 2 // Two characters per cell
			screenY := y

			color := r.cellColor(game, x, y)
			style := tcell.StyleDefault.Background(color).Foreground(color)

			// Draw two characters for each cell
			r.Screen.SetContent(screenX, screenY, ' ', nil, style)
			r.Screen.SetContent(screenX+1, screenY, ' ', nil, style)
		}
	}
}

// cellColor returns the display color for the cell at (x, y)
func (r *Renderer) cellColor(game *Game, x, y int) tcell.Color {
	if !game.Grid.Cells[y][x] {
		return r.Theme.Dead
	}
	if r.Theme.AgeColored && game.Ages != nil {
		return GetRainbowColor(game.Ages[y][x])
	}
	return r.Theme.Alive
}

// DrawCursor draws the cursor at the current position (when paused)
func (r *Renderer) DrawCursor(game *Game) {
	if !game.Paused {
		return
	}

	screenX := game.CursorX * 2
	screenY := game.CursorY

	// Draw cursor brackets over the cell's own color
	style := tcell.StyleDefault.Foreground(r.Theme.CursorFg).Background(r.cellColor(game, game.CursorX, game.CursorY))
	r.Screen.SetContent(screenX, screenY, '[', nil, style)
	r.Screen.SetContent(screenX+1, screenY, ']', nil, style)
}

// DrawStatusBar renders the status bar at the bottom of the screen
func (r *Renderer) DrawStatusBar(game *Game, screenWidth, screenHeight int) {
	statusY := screenHeight - 1

	// Build status string
	status := fmt.Sprintf(" Gen: %d | Cells: %d | FPS: %d | Rule: %s",
		game.Generation,
		game.Grid.CountAlive(),
		game.FPS,
		game.Rule.String(),
	)

	if game.Paused {
		status += " | [PAUSED]"
	}

	// Pad to fill width, with the help hint right-aligned when it fits
	helpHint := "? help "
	if len(status)+len(helpHint) < screenWidth {
		status += strings.Repeat(" ", screenWidth-len(status)-len(helpHint)) + helpHint
	}
	for len(status) < screenWidth {
		status += " "
	}

	style := tcell.StyleDefault.Foreground(r.Theme.StatusFg).Background(r.Theme.StatusBg)

	for x, ch := range status {
		if x >= screenWidth {
			break
		}
		r.Screen.SetContent(x, statusY, ch, nil, style)
	}
}

// DrawHelp renders the key reference as a box centered on the screen
func (r *Renderer) DrawHelp(game *Game, screenWidth, screenHeight int) {
	if !game.ShowHelp {
		return
	}

	keyColumnWidth := 0
	for _, binding := range KeyBindings {
		keyColumnWidth = max(keyColumnWidth, len(binding[0]))
	}
	lines := []string{"Controls", ""}
	for _, binding := range KeyBindings {
		lines = append(lines, fmt.Sprintf("%-*s  %s", keyColumnWidth, binding[0], binding[1]))
	}

	innerWidth := 0
	for _, line := range lines {
		innerWidth = max(innerWidth, len(line))
	}
	boxWidth := innerWidth + 4 // Border plus one space of padding per side
	boxHeight := len(lines) + 2
	left := max((screenWidth-boxWidth)/2, 0)
	top := max((screenHeight-1-boxHeight)/2, 0) // Center above the status bar

	style := tcell.StyleDefault.Foreground(r.Theme.StatusFg).Background(r.Theme.StatusBg)
	titleStyle := style.Bold(true)

	for row := 0; row < boxHeight; row++ {
		var text string
		var leftEdge, fill, rightEdge rune
		switch row {
		case 0:
			leftEdge, fill, rightEdge = '┌', '─', '┐'
		case boxHeight - 1:
			leftEdge, fill, rightEdge = '└', '─', '┘'
		default:
			leftEdge, fill, rightEdge = '│', ' ', '│'
			text = lines[row-1]
		}

		r.Screen.SetContent(left, top+row, leftEdge, nil, style)
		for col := 1; col < boxWidth-1; col++ {
			r.Screen.SetContent(left+col, top+row, fill, nil, style)
		}
		r.Screen.SetContent(left+boxWidth-1, top+row, rightEdge, nil, style)

		textStyle := style
		if row == 1 {
			textStyle = titleStyle
		}
		for col, ch := range text {
			r.Screen.SetContent(left+2+col, top+row, ch, nil, textStyle)
		}
	}
}

// Show presents the screen buffer
func (r *Renderer) Show() {
	r.Screen.Show()
}

// GetGridSize calculates the grid dimensions based on terminal size
// Returns width, height (accounting for 2-char cells and status bar)
func GetGridSize(screen tcell.Screen) (int, int) {
	screenWidth, screenHeight := screen.Size()

	// Each cell is 2 characters wide
	gridWidth := screenWidth / 2

	// Reserve 1 row for status bar
	gridHeight := screenHeight - 1

	// Minimum size constraints
	if gridWidth < 10 {
		gridWidth = 10
	}
	if gridHeight < 10 {
		gridHeight = 10
	}

	return gridWidth, gridHeight
}

// InitScreen initializes and returns a new tcell screen
func InitScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.EnableMouse()
	screen.Clear()

	return screen, nil
}
