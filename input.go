package main

import (
	"github.com/gdamore/tcell/v2"
)

// Action represents a user input action
type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionPause
	ActionStep
	ActionRandomize
	ActionClear
	ActionSpeedUp
	ActionSlowDown
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionToggleCell
)

// HandleInput processes keyboard and mouse events, returns the action to take
func HandleInput(ev tcell.Event, game *Game) Action {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return handleKeyEvent(ev, game)
	case *tcell.EventMouse:
		return handleMouseEvent(ev, game)
	case *tcell.EventResize:
		// Resize is handled in the main loop
		return ActionNone
	}
	return ActionNone
}

func handleKeyEvent(ev *tcell.EventKey, game *Game) Action {
	// Check for Ctrl+C first
	if ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	switch ev.Key() {
	case tcell.KeyEscape:
		return ActionQuit
	case tcell.KeyRune:
		return handleRuneKey(ev.Rune(), game)
	case tcell.KeyUp:
		return ActionMoveUp
	case tcell.KeyDown:
		return ActionMoveDown
	case tcell.KeyLeft:
		return ActionMoveLeft
	case tcell.KeyRight:
		return ActionMoveRight
	case tcell.KeyEnter:
		return ActionToggleCell
	}
	return ActionNone
}

func handleRuneKey(r rune, game *Game) Action {
	switch r {
	case 'q', 'Q':
		return ActionQuit
	case ' ':
		return ActionPause
	case 'n', 'N':
		return ActionStep
	case 'r', 'R':
		return ActionRandomize
	case 'c', 'C':
		return ActionClear
	case '+', '=':
		return ActionSpeedUp
	case '-':
		return ActionSlowDown
	}
	return ActionNone
}

func handleMouseEvent(ev *tcell.EventMouse, game *Game) Action {
	// Mouse support for future enhancement
	// Could toggle cells on click when paused
	_ = ev
	_ = game
	return ActionNone
}

// ApplyAction performs the action on the game state
// Returns true if the game should quit
func ApplyAction(action Action, game *Game, density float64, rng func() float64) bool {
	switch action {
	case ActionQuit:
		return true

	case ActionPause:
		game.Paused = !game.Paused

	case ActionStep:
		if game.Paused {
			stepWithAges(game)
		}

	case ActionRandomize:
		game.Grid.Randomize(density, rng)
		game.Generation = 0
		resetAges(game)

	case ActionClear:
		game.Grid.Clear()
		game.Generation = 0
		resetAges(game)

	case ActionSpeedUp:
		if game.FPS < 60 {
			game.FPS += 5
			if game.FPS > 60 {
				game.FPS = 60
			}
		}

	case ActionSlowDown:
		if game.FPS > 1 {
			game.FPS -= 5
			if game.FPS < 1 {
				game.FPS = 1
			}
		}

	case ActionMoveUp:
		if game.Paused && game.CursorY > 0 {
			game.CursorY--
		}

	case ActionMoveDown:
		if game.Paused && game.CursorY < game.Grid.Height-1 {
			game.CursorY++
		}

	case ActionMoveLeft:
		if game.Paused && game.CursorX > 0 {
			game.CursorX--
		}

	case ActionMoveRight:
		if game.Paused && game.CursorX < game.Grid.Width-1 {
			game.CursorX++
		}

	case ActionToggleCell:
		if game.Paused {
			game.SetCell(game.CursorX, game.CursorY, !game.Grid.Get(game.CursorX, game.CursorY))
		}
	}
	return false
}
