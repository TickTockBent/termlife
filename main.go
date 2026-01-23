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
}

func main() {
	// CLI flags
	patternFlag := flag.String("pattern", "random", "Initial pattern: random, glider, blinker, pulsar, gosper-gun, etc.")
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

	// Initialize random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randFunc := rng.Float64

	// Check if GIF mode is enabled
	gifMode := *gifFlag > 0

	if gifMode {
		// GIF mode - no terminal UI
		runGIFMode(
			*gifFlag,
			*gifOutFlag,
			*gifScaleFlag,
			*gifDelayFlag,
			*gifLoopFlag,
			*sizeFlag,
			*patternFlag,
			*wrapFlag,
			density,
			rule,
			theme,
			themeName,
			randFunc,
		)
		return
	}

	// Interactive mode - initialize terminal UI
	screen, err := InitScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	// Determine grid size
	var gridWidth, gridHeight int
	if *sizeFlag != "" {
		gridWidth, gridHeight, err = parseSize(*sizeFlag)
		if err != nil {
			screen.Fini()
			fmt.Fprintf(os.Stderr, "Error parsing size: %v\n", err)
			os.Exit(1)
		}
	} else {
		gridWidth, gridHeight = GetGridSize(screen)
	}

	// Create grid
	grid := NewGrid(gridWidth, gridHeight, *wrapFlag)

	// Initialize pattern
	patternName := strings.ToLower(*patternFlag)
	if patternName == "random" {
		grid.Randomize(density, randFunc)
	} else {
		pattern, ok := Patterns[patternName]
		if !ok {
			screen.Fini()
			fmt.Fprintf(os.Stderr, "Unknown pattern: %s\nAvailable: %s\n", *patternFlag, strings.Join(PatternNames(), ", "))
			os.Exit(1)
		}
		PlacePattern(grid, pattern, gridWidth/2, gridHeight/2)
	}

	// Create game state
	game := &Game{
		Grid:       grid,
		Rule:       rule,
		Theme:      theme,
		FPS:        fps,
		Paused:     false,
		Generation: 0,
		CursorX:    gridWidth / 2,
		CursorY:    gridHeight / 2,
	}

	// Initialize ages for rainbow mode
	if theme.Name == "Rainbow" {
		game.Ages = make([][]int, gridHeight)
		for y := range game.Ages {
			game.Ages[y] = make([]int, gridWidth)
		}
	}

	// Create renderer
	renderer := NewRenderer(screen, theme)

	// Run game loop
	runGameLoop(screen, renderer, game, density, randFunc)
}

func runGIFMode(
	frames int,
	output string,
	scale int,
	delay int,
	loop int,
	sizeFlag string,
	patternFlag string,
	wrap bool,
	density float64,
	rule *Rule,
	theme *Theme,
	themeName string,
	randFunc func() float64,
) {
	// Validate GIF config
	cfg := GIFConfig{
		Frames:    frames,
		Output:    output,
		Scale:     scale,
		Delay:     delay,
		Loop:      loop,
		ThemeName: themeName,
	}

	hasSize := sizeFlag != ""
	if err := ValidateGIFConfig(cfg, hasSize); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse size (required in GIF mode)
	gridWidth, gridHeight, err := parseSize(sizeFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing size: %v\n", err)
		os.Exit(1)
	}

	// Create grid
	grid := NewGrid(gridWidth, gridHeight, wrap)

	// Initialize pattern
	patternName := strings.ToLower(patternFlag)
	if patternName == "random" {
		grid.Randomize(density, randFunc)
	} else {
		pattern, ok := Patterns[patternName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown pattern: %s\nAvailable: %s\n", patternFlag, strings.Join(PatternNames(), ", "))
			os.Exit(1)
		}
		PlacePattern(grid, pattern, gridWidth/2, gridHeight/2)
	}

	// Create game state
	game := &Game{
		Grid:       grid,
		Rule:       rule,
		Theme:      theme,
		FPS:        10, // Not used in GIF mode
		Paused:     false,
		Generation: 0,
	}

	// Initialize ages for rainbow mode
	if themeName == "rainbow" {
		game.Ages = make([][]int, gridHeight)
		for y := range game.Ages {
			game.Ages[y] = make([]int, gridWidth)
		}
	}

	// Export GIF
	if err := ExportGIF(game, cfg, density, randFunc); err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting GIF: %v\n", err)
		os.Exit(1)
	}
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

			// Handle resize events
			if _, ok := ev.(*tcell.EventResize); ok {
				screen.Sync()
				continue
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
