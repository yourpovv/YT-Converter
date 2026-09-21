package utils

import (
	"fmt"
	"strings"
)

var (
	Minty   = []string{"00ffcc", "66ffcc", "99ffe6", "ccfff2"}
	Error   = []string{"cc0000", "e60000", "ff3333", "ff6666", "ff9999"}
	Success = []string{"006600", "009900", "33cc33", "66ff66", "99ff99"}
)

type Color struct {
	R, G, B int
}

// ParseHexColor parses a 3 or 6-digit hex color string
func ParseHexColor(hex string) (Color, error) {
	var c Color
	switch len(hex) {
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &c.R, &c.G, &c.B)
		if err != nil {
			return c, err
		}
	case 3:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &c.R, &c.G, &c.B)
		if err != nil {
			return c, err
		}
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		return c, fmt.Errorf("invalid hex color: %s", hex)
	}
	return c, nil
}

// Colorize wraps text in a 24-bit ANSI foreground escape
func Colorize(text string, color Color) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", color.R, color.G, color.B, text)
}

// Gradient streaks text across the palette left to right
func Gradient(text string, hexPalette []string) string {
	text = strings.TrimSpace(text)
	if len(text) == 0 || len(hexPalette) == 0 {
		return text
	}

	colors := parsePalette(hexPalette)
	if len(colors) == 0 {
		return text
	}

	var builder strings.Builder
	for index, ch := range text {
		builder.WriteString(Colorize(string(ch), colorAt(colors, index, len(text))))
	}

	return builder.String()
}

// parsePalette converts hex strings to colors
func parsePalette(hexPalette []string) []Color {
	colors := make([]Color, 0, len(hexPalette))
	for _, hex := range hexPalette {
		if c, err := ParseHexColor(hex); err == nil {
			colors = append(colors, c)
		}
	}
	return colors
}

// colorAt interpolates the palette color for one character position
func colorAt(colors []Color, index, charCount int) Color {
	if charCount <= 1 {
		return colors[0]
	}
	segmentCount := float64(len(colors) - 1)
	pos := float64(index) / float64(charCount-1) * segmentCount
	lower := int(pos)
	upper := lower + 1
	if upper >= len(colors) {
		upper = len(colors) - 1
	}
	return lerpColor(colors[lower], colors[upper], pos-float64(lower))
}

// lerpColor blends two colors by fraction in [0,1]
func lerpColor(first, second Color, fraction float64) Color {
	blend := func(a, b int) int {
		return int(float64(a)*(1-fraction) + float64(b)*fraction)
	}
	return Color{R: blend(first.R, second.R), G: blend(first.G, second.G), B: blend(first.B, second.B)}
}
