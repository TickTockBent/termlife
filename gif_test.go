package main

import (
	"image/gif"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateGIFConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     GIFConfig
		hasSize bool
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     GIFConfig{Frames: 10, Output: "test.gif", Scale: 4, Delay: 100},
			hasSize: true,
			wantErr: false,
		},
		{
			name:    "zero frames",
			cfg:     GIFConfig{Frames: 0, Output: "test.gif", Scale: 4, Delay: 100},
			hasSize: true,
			wantErr: true,
		},
		{
			name:    "negative frames",
			cfg:     GIFConfig{Frames: -1, Output: "test.gif", Scale: 4, Delay: 100},
			hasSize: true,
			wantErr: true,
		},
		{
			name:    "zero scale",
			cfg:     GIFConfig{Frames: 10, Output: "test.gif", Scale: 0, Delay: 100},
			hasSize: true,
			wantErr: true,
		},
		{
			name:    "missing size",
			cfg:     GIFConfig{Frames: 10, Output: "test.gif", Scale: 4, Delay: 100},
			hasSize: false,
			wantErr: true,
		},
		{
			name:    "low delay still valid",
			cfg:     GIFConfig{Frames: 10, Output: "test.gif", Scale: 4, Delay: 5},
			hasSize: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGIFConfig(tt.cfg, tt.hasSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGIFConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildPalette(t *testing.T) {
	tests := []struct {
		name       string
		themeName  string
		rainbow    bool
		wantColors int
	}{
		{"white theme", "white", false, 2},
		{"green theme", "green", false, 2},
		{"matrix theme", "matrix", false, 2},
		{"amber theme", "amber", false, 2},
		{"cyan theme", "cyan", false, 2},
		{"rainbow theme", "rainbow", true, 9},
		{"unknown theme defaults", "unknown", false, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			palette := buildPalette(tt.themeName, tt.rainbow)
			if len(palette) != tt.wantColors {
				t.Errorf("buildPalette() returned %d colors, want %d", len(palette), tt.wantColors)
			}
		})
	}
}

func TestGIFExport(t *testing.T) {
	// Create a temporary directory for test output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.gif")

	// Create a simple grid with a blinker pattern
	grid := NewGrid(20, 20, false)
	PlacePattern(grid, Patterns["blinker"], 10, 10)

	game := &Game{
		Grid:       grid,
		Rule:       ConwayRule(),
		Theme:      Themes["white"],
		FPS:        10,
		Generation: 0,
	}

	cfg := GIFConfig{
		Frames:    5,
		Output:    outputPath,
		Scale:     2,
		Delay:     100,
		Loop:      0,
		ThemeName: "white",
	}

	// Use a deterministic "random" function
	randFunc := func() float64 { return 0.5 }

	err := ExportGIF(game, cfg, 0.25, randFunc)
	if err != nil {
		t.Fatalf("ExportGIF() error = %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("GIF file was not created")
	}

	// Open and verify the GIF
	f, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open GIF: %v", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		t.Fatalf("Failed to decode GIF: %v", err)
	}

	// Verify frame count
	if len(g.Image) != 5 {
		t.Errorf("GIF has %d frames, want 5", len(g.Image))
	}

	// Verify dimensions (20 cells * 2 scale = 40 pixels)
	if g.Image[0].Bounds().Dx() != 40 || g.Image[0].Bounds().Dy() != 40 {
		t.Errorf("GIF dimensions = %dx%d, want 40x40",
			g.Image[0].Bounds().Dx(), g.Image[0].Bounds().Dy())
	}

	// Verify delay values (100ms / 10 = 10 centiseconds)
	for i, delay := range g.Delay {
		if delay != 10 {
			t.Errorf("Frame %d delay = %d, want 10", i, delay)
		}
	}
}

func TestGIFExportRainbow(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "rainbow.gif")

	grid := NewGrid(15, 15, false)
	PlacePattern(grid, Patterns["glider"], 7, 7)

	// Initialize ages for rainbow mode
	ages := make([][]int, 15)
	for y := range ages {
		ages[y] = make([]int, 15)
	}

	game := &Game{
		Grid:       grid,
		Rule:       ConwayRule(),
		Theme:      Themes["rainbow"],
		FPS:        10,
		Generation: 0,
		Ages:       ages,
	}

	cfg := GIFConfig{
		Frames:    10,
		Output:    outputPath,
		Scale:     4,
		Delay:     100,
		Loop:      0,
		ThemeName: "rainbow",
	}

	randFunc := func() float64 { return 0.5 }

	err := ExportGIF(game, cfg, 0.25, randFunc)
	if err != nil {
		t.Fatalf("ExportGIF() error = %v", err)
	}

	// Verify the file was created and is a valid GIF
	f, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open GIF: %v", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		t.Fatalf("Failed to decode GIF: %v", err)
	}

	// Rainbow palette should have at least 9 colors (GIF may add more)
	if len(g.Image[0].Palette) < 9 {
		t.Errorf("Rainbow GIF palette has %d colors, want at least 9", len(g.Image[0].Palette))
	}
}

func TestGIFScaleMinimum(t *testing.T) {
	cfg := GIFConfig{
		Frames: 10,
		Output: "test.gif",
		Scale:  0,
		Delay:  100,
	}

	err := ValidateGIFConfig(cfg, true)
	if err == nil {
		t.Error("ValidateGIFConfig() should error on scale < 1")
	}
}

func TestGIFRequiresSize(t *testing.T) {
	cfg := GIFConfig{
		Frames: 10,
		Output: "test.gif",
		Scale:  4,
		Delay:  100,
	}

	err := ValidateGIFConfig(cfg, false)
	if err == nil {
		t.Error("ValidateGIFConfig() should error when size not provided")
	}
}

func TestGIFColorMapping(t *testing.T) {
	// Test that each theme produces correct colors
	themes := []string{"white", "green", "matrix", "amber", "cyan"}

	for _, themeName := range themes {
		colors, ok := GIFColors[themeName]
		if !ok {
			t.Errorf("GIFColors missing theme %q", themeName)
			continue
		}

		// Verify colors have full alpha
		if colors.Alive.A != 255 || colors.Dead.A != 255 {
			t.Errorf("Theme %q colors should have full alpha", themeName)
		}
	}
}
