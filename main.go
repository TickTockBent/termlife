package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Version is set at build time
var Version = "dev"

func main() {
	// CLI flags
	patternFlag := flag.String("pattern", "random", "Initial pattern: "+strings.Join(PatternNames(), ", "))
	ruleFlag := flag.String("rule", "B3/S23", "Birth/survival rule string (e.g., B3/S23, B36/S23)")
	colorFlag := flag.String("color", "white", "Color theme: white, green, matrix, amber, cyan, rainbow")
	fpsFlag := flag.Int("fps", 10, "Frames per second (1-60)")
	sizeFlag := flag.String("size", "", "Grid dimensions as WxH (default: auto-fit terminal)")
	wrapFlag := flag.Bool("wrap", false, "Enable toroidal wrapping")
	densityFlag := flag.Float64("density", 0.25, "Cell density for random initialization (0.0-1.0)")
	versionFlag := flag.Bool("version", false, "Show version")

	// GIF export flags
	gifFlag := flag.Int("gif", 0, "Number of frames to render for GIF export (enables GIF mode)")
	gifOutFlag := flag.String("gif-out", "termlife.gif", "Output file path for GIF")
	gifScaleFlag := flag.Int("gif-scale", 4, "Pixels per cell in GIF output")
	gifDelayFlag := flag.Int("gif-delay", 100, "Milliseconds between GIF frames")
	gifLoopFlag := flag.Int("gif-loop", 0, "Number of times to loop GIF (0 = infinite)")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("termlife %s\n", Version)
		os.Exit(0)
	}

	// Validate FPS
	fps := *fpsFlag
	if fps < 1 {
		fps = 1
	} else if fps > 60 {
		fps = 60
	}

	// Validate density
	density := *densityFlag
	if density < 0 {
		density = 0
	} else if density > 1 {
		density = 1
	}

	// Parse rule
	rule, err := ParseRule(*ruleFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing rule: %v\n", err)
		os.Exit(1)
	}

	// Get theme
	themeName := strings.ToLower(*colorFlag)
	theme, ok := Themes[themeName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown theme: %s\nAvailable: %s\n", *colorFlag, strings.Join(ThemeNames(), ", "))
		os.Exit(1)
	}

	// Validate pattern and size up front so errors print before the terminal is taken over
	patternName := strings.ToLower(*patternFlag)
	if _, ok := Patterns[patternName]; !ok && patternName != "random" {
		fmt.Fprintf(os.Stderr, "Unknown pattern: %s\nAvailable: %s\n", *patternFlag, strings.Join(PatternNames(), ", "))
		os.Exit(1)
	}

	var gridWidth, gridHeight int
	if *sizeFlag != "" {
		gridWidth, gridHeight, err = parseSize(*sizeFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing size: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randFunc := rng.Float64

	// GIF mode - no terminal UI
	if *gifFlag > 0 {
		cfg := GIFConfig{
			Frames:    *gifFlag,
			Output:    *gifOutFlag,
			Scale:     *gifScaleFlag,
			Delay:     *gifDelayFlag,
			Loop:      *gifLoopFlag,
			ThemeName: themeName,
		}
		if err := ValidateGIFConfig(cfg, *sizeFlag != ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		grid, err := SetupGrid(gridWidth, gridHeight, *wrapFlag, patternName, density, randFunc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		game := NewGame(grid, rule, theme, fps)
		if err := ExportGIF(game, cfg, density, randFunc); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting GIF: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Interactive mode - initialize terminal UI
	screen, err := InitScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	autoSize := *sizeFlag == ""
	if autoSize {
		gridWidth, gridHeight = GetGridSize(screen)
	}

	grid, err := SetupGrid(gridWidth, gridHeight, *wrapFlag, patternName, density, randFunc)
	if err != nil {
		screen.Fini()
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	game := NewGame(grid, rule, theme, fps)
	game.AutoSize = autoSize

	renderer := NewRenderer(screen, theme)
	runGameLoop(screen, renderer, game, density, randFunc)
}

func runGameLoop(screen tcell.Screen, renderer *Renderer, game *Game, density float64, randFunc func() float64) {
	// Channel for events
	eventChan := make(chan tcell.Event)
	quitChan := make(chan struct{})

	// Poll events in separate goroutine
	go func() {
		for {
			ev := screen.PollEvent()
			select {
			case eventChan <- ev:
			case <-quitChan:
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second / time.Duration(game.FPS))
	defer ticker.Stop()

	lastFPS := game.FPS

	for {
		select {
		case ev := <-eventChan:
			if ev == nil {
				continue
			}

			// Handle resize events; auto-sized grids follow the terminal
			if _, ok := ev.(*tcell.EventResize); ok {
				if game.AutoSize {
					game.Resize(GetGridSize(screen))
				}
				screen.Sync()
				break
			}

			action := HandleInput(ev, game)
			if ApplyAction(action, game, density, randFunc) {
				close(quitChan)
				return
			}

			// Update ticker if FPS changed
			if game.FPS != lastFPS {
				ticker.Reset(time.Second / time.Duration(game.FPS))
				lastFPS = game.FPS
			}

		case <-ticker.C:
			// Step simulation if not paused
			if !game.Paused {
				stepWithAges(game)
			}
		}

		// Render
		renderer.Clear()
		renderer.DrawGrid(game)
		renderer.DrawCursor(game)
		screenWidth, screenHeight := screen.Size()
		renderer.DrawStatusBar(game, screenWidth, screenHeight)
		renderer.DrawHelp(game, screenWidth, screenHeight)
		renderer.Show()
	}
}

func parseSize(s string) (int, int, error) {
	parts := strings.Split(strings.ToLower(s), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected WxH format, got %q", s)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %v", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %v", err)
	}

	if width < 10 || height < 10 {
		return 0, 0, fmt.Errorf("minimum size is 10x10")
	}

	return width, height, nil
}
