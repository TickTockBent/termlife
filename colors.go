package main

import "github.com/gdamore/tcell/v2"

// Theme defines the color scheme for rendering
type Theme struct {
	Name       string
	Alive      tcell.Color
	Dead       tcell.Color
	Background tcell.Color
	StatusFg   tcell.Color
	StatusBg   tcell.Color
	CursorFg   tcell.Color
	AgeColored bool // Alive cells are colored by age (see GetRainbowColor)
}

// Themes contains all available color themes
var Themes = map[string]*Theme{
	"white": {
		Name:       "White",
		Alive:      tcell.ColorWhite,
		Dead:       tcell.ColorBlack,
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorWhite,
		StatusBg:   tcell.ColorDarkGray,
		CursorFg:   tcell.ColorYellow,
	},
	"green": {
		Name:       "Green",
		Alive:      tcell.ColorGreen,
		Dead:       tcell.ColorBlack,
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorGreen,
		StatusBg:   tcell.ColorDarkGray,
		CursorFg:   tcell.ColorYellow,
	},
	"matrix": {
		Name:       "Matrix",
		Alive:      tcell.ColorLime,
		Dead:       tcell.ColorDarkGreen,
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorLime,
		StatusBg:   tcell.ColorDarkGreen,
		CursorFg:   tcell.ColorWhite,
	},
	"amber": {
		Name:       "Amber",
		Alive:      tcell.ColorOrange,
		Dead:       tcell.NewRGBColor(40, 20, 0),
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorOrange,
		StatusBg:   tcell.NewRGBColor(40, 20, 0),
		CursorFg:   tcell.ColorYellow,
	},
	"cyan": {
		Name:       "Cyan",
		Alive:      tcell.ColorAqua,
		Dead:       tcell.ColorDarkBlue,
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorAqua,
		StatusBg:   tcell.ColorDarkBlue,
		CursorFg:   tcell.ColorYellow,
	},
	"rainbow": {
		Name:       "Rainbow",
		Alive:      tcell.ColorWhite, // Not used directly - see age-based coloring
		Dead:       tcell.ColorBlack,
		Background: tcell.ColorBlack,
		StatusFg:   tcell.ColorWhite,
		StatusBg:   tcell.ColorDarkGray,
		CursorFg:   tcell.ColorYellow,
		AgeColored: true,
	},
}

// ThemeNames returns a list of available theme names
func ThemeNames() []string {
	return []string{"white", "green", "matrix", "amber", "cyan", "rainbow"}
}

// RainbowColors provides age-based colors for the rainbow theme
var RainbowColors = []tcell.Color{
	tcell.ColorRed,
	tcell.ColorOrange,
	tcell.ColorYellow,
	tcell.ColorGreen,
	tcell.ColorAqua,
	tcell.ColorBlue,
	tcell.ColorPurple,
	tcell.ColorFuchsia,
}

// GetRainbowColor returns a color based on cell age (capped at max colors)
func GetRainbowColor(age int) tcell.Color {
	if age <= 0 {
		return tcell.ColorRed
	}
	if age >= len(RainbowColors) {
		return RainbowColors[len(RainbowColors)-1]
	}
	return RainbowColors[age]
}
