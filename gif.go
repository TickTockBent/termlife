package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"strings"
)

// GIFConfig holds settings for GIF export
type GIFConfig struct {
	Frames    int
	Output    string
	Scale     int
	Delay     int // milliseconds
	Loop      int
	ThemeName string // Theme name for color lookup
}

// DefaultGIFConfig returns default GIF export settings
func DefaultGIFConfig() GIFConfig {
	return GIFConfig{
		Frames:    0,
		Output:    "termlife.gif",
		Scale:     4,
		Delay:     100,
		Loop:      0,
		ThemeName: "white",
	}
}

// ValidateGIFConfig checks the GIF configuration for errors
func ValidateGIFConfig(cfg GIFConfig, hasSize bool) error {
	if cfg.Frames <= 0 {
		return fmt.Errorf("frame count must be positive, got %d", cfg.Frames)
	}
	if cfg.Scale < 1 {
		return fmt.Errorf("scale must be at least 1, got %d", cfg.Scale)
	}
	if !hasSize {
		return fmt.Errorf("--size is required when using --gif")
	}
	return nil
}

// GIF color definitions matching the spec
var GIFColors = map[string]struct {
	Alive color.RGBA
	Dead  color.RGBA
}{
	"white": {
		Alive: color.RGBA{255, 255, 255, 255},
		Dead:  color.RGBA{0, 0, 0, 255},
	},
	"green": {
		Alive: color.RGBA{0, 255, 0, 255},
		Dead:  color.RGBA{0, 0, 0, 255},
	},
	"matrix": {
		Alive: color.RGBA{0, 255, 0, 255},
		Dead:  color.RGBA{10, 48, 0, 255},
	},
	"amber": {
		Alive: color.RGBA{255, 191, 0, 255},
		Dead:  color.RGBA{61, 38, 0, 255},
	},
	"cyan": {
		Alive: color.RGBA{0, 255, 255, 255},
		Dead:  color.RGBA{0, 51, 51, 255},
	},
}

// ExportGIF renders the simulation to a GIF file
func ExportGIF(game *Game, cfg GIFConfig, density float64, rng func() float64) error {
	// Create the output first so a bad path fails before any rendering work
	f, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Print header
	fmt.Printf("termlife: rendering %d frames at %dx%d\n", cfg.Frames, game.Grid.Width, game.Grid.Height)

	// Build palette from theme name
	isRainbow := strings.ToLower(cfg.ThemeName) == "rainbow"
	palette := buildPalette(cfg.ThemeName, isRainbow)

	anim := &gif.GIF{
		LoopCount: cfg.Loop,
	}

	// Clamp delay minimum to 10ms (GIF spec minimum)
	delay := cfg.Delay
	if delay < 10 {
		delay = 10
	}

	// Estimate output size and warn if large
	estimatedBytes := cfg.Frames * (game.Grid.Width * cfg.Scale) * (game.Grid.Height * cfg.Scale) / 8
	if estimatedBytes > 10*1024*1024 {
		fmt.Printf("termlife: warning - estimated output size ~%d MB\n", estimatedBytes/1024/1024)
	}
	if game.Grid.Width > 500 || game.Grid.Height > 500 {
		fmt.Printf("termlife: warning - large grid may use significant memory\n")
	}

	for i := 0; i < cfg.Frames; i++ {
		// Render frame
		img := renderFrame(game, cfg.Scale, palette, isRainbow)
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, delay/10) // Convert ms to centiseconds
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)

		// Advance simulation
		stepWithAges(game)

		// Update progress
		printProgress(i+1, cfg.Frames)
	}

	fmt.Println() // Newline after progress bar

	// Write file
	if err := gif.EncodeAll(f, anim); err != nil {
		f.Close()
		os.Remove(cfg.Output)
		return fmt.Errorf("failed to encode GIF: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to write GIF: %w", err)
	}

	// Get file size
	info, err := os.Stat(cfg.Output)
	if err == nil {
		sizeKB := float64(info.Size()) / 1024
		if sizeKB > 1024 {
			fmt.Printf("termlife: wrote %s (%.1f MB)\n", cfg.Output, sizeKB/1024)
		} else {
			fmt.Printf("termlife: wrote %s (%.1f KB)\n", cfg.Output, sizeKB)
		}
	} else {
		fmt.Printf("termlife: wrote %s\n", cfg.Output)
	}

	return nil
}

func renderFrame(game *Game, scale int, palette color.Palette, rainbow bool) *image.Paletted {
	width := game.Grid.Width * scale
	height := game.Grid.Height * scale

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	// Fill with background (dead color at index 0)
	for i := range img.Pix {
		img.Pix[i] = 0
	}

	for y := 0; y < game.Grid.Height; y++ {
		for x := 0; x < game.Grid.Width; x++ {
			if !game.Grid.Cells[y][x] {
				continue // Already filled with dead color
			}

			var colorIndex uint8 = 1 // Default alive color

			// Rainbow mode uses age-based colors
			if rainbow && game.Ages != nil {
				age := game.Ages[y][x]
				// Clamp to palette size (index 0 is dead, 1-8 are rainbow colors)
				if age >= len(palette)-1 {
					colorIndex = uint8(len(palette) - 1)
				} else {
					colorIndex = uint8(age + 1)
				}
			}

			// Fill scaled cell
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetColorIndex(x*scale+dx, y*scale+dy, colorIndex)
				}
			}
		}
	}

	return img
}

func buildPalette(themeName string, rainbow bool) color.Palette {
	if rainbow {
		// Rainbow palette: dead + 8 age-based colors
		return color.Palette{
			color.RGBA{0, 0, 0, 255},     // 0: dead/black
			color.RGBA{255, 0, 0, 255},   // 1: red (age 0)
			color.RGBA{255, 165, 0, 255}, // 2: orange (age 1)
			color.RGBA{255, 255, 0, 255}, // 3: yellow (age 2)
			color.RGBA{0, 255, 0, 255},   // 4: green (age 3)
			color.RGBA{0, 255, 255, 255}, // 5: cyan (age 4)
			color.RGBA{0, 0, 255, 255},   // 6: blue (age 5)
			color.RGBA{128, 0, 255, 255}, // 7: violet (age 6)
			color.RGBA{255, 0, 255, 255}, // 8: magenta (age 7+)
		}
	}

	// Look up theme colors
	themeLower := strings.ToLower(themeName)
	if colors, ok := GIFColors[themeLower]; ok {
		return color.Palette{colors.Dead, colors.Alive}
	}

	// Default to white on black
	return color.Palette{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 255, 255, 255},
	}
}

func printProgress(current, total int) {
	barWidth := 40
	progress := float64(current) / float64(total)
	filled := int(progress * float64(barWidth))

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	fmt.Printf("\r[%s] %d/%d", bar, current, total)
}
