// provides ansi color gradients for terminal&ssh applications
package simple_colors

import (
	"fmt"
	"math"
	"strings"
)

// rgb color struct
type Color struct {
	R, G, B uint8
}

type hsl struct {
	h, s, l float64 // h: 0-360, s: 0-1, l: 0-1
}

// function which create rgb color object using red, green and blue colors
func Rgb(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

// just color string in one rgb color
func Paint(text string, color Color) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", color.R, color.G, color.B, text)
}

// create gradient, flexmode = true enable smooth, flexmode = false disable it. text - text, which will be colored. colors - some Rgb() colors
func Gradient(flexMode bool, text string, colors ...Color) string {
	if len(text) == 0 || len(colors) == 0 {
		return text
	}

	if len(colors) == 1 {
		return colors[0].wrap(text)
	}

	hslColors := make([]hsl, len(colors))
	for i, c := range colors {
		hslColors[i] = rgbToHsl(c)
	}

	segments := len(colors) - 1
	segmentLen := len(text) / segments
	extraChars := len(text) % segments

	if flexMode {
		segmentLen = adjustMultiSteps(hslColors, len(text))
	}

	var builder strings.Builder
	pos := 0

	for seg := 0; seg < segments; seg++ {
		start := hslColors[seg]
		end := hslColors[seg+1]
		charsInSegment := segmentLen
		if seg == segments-1 {
			charsInSegment += extraChars
		}

		for i := 0; i < charsInSegment; i++ {
			if pos >= len(text) {
				break
			}

			t := float64(i) / float64(charsInSegment-1)
			color := interpolate(start, end, t)
			builder.WriteString(color.wrapChar(text[pos]))
			pos++
		}
	}

	return builder.String()
}

func (c Color) wrap(text string) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.R, c.G, c.B, text)
}

func (c Color) wrapChar(char byte) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c", c.R, c.G, c.B, char)
}

func rgbToHsl(c Color) hsl {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	h := 0.0
	s := 0.0
	l := (max + min) / 2.0

	if delta != 0 {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2.0 - max - min)
		}

		switch max {
		case r:
			h = (g - b) / delta
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/delta + 2
		case b:
			h = (r-g)/delta + 4
		}
		h *= 60
	}
	return hsl{h: h, s: s, l: l}
}

func hslToRgb(h, s, l float64) Color {
	if s == 0 {
		val := uint8(math.Round(l * 255))
		return Rgb(val, val, val)
	}

	var r, g, b float64
	h /= 360.0

	q := l + s - l*s
	if l < 0.5 {
		q = l * (1 + s)
	}
	p := 2*l - q

	r = hueToRgb(p, q, h+1.0/3)
	g = hueToRgb(p, q, h)
	b = hueToRgb(p, q, h-1.0/3)

	return Rgb(
		uint8(math.Round(r*255)),
		uint8(math.Round(g*255)),
		uint8(math.Round(b*255)),
	)
}

func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 0.5:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

func interpolate(start, end hsl, t float64) Color {
	h := lerp(start.h, end.h, t, true) // Учитываем циклическую природу Hue
	s := lerp(start.s, end.s, t, false)
	l := lerp(start.l, end.l, t, false)
	return hslToRgb(h, s, l)
}

func lerp(a, b, t float64, isHue bool) float64 {
	if !isHue {
		return a + (b-a)*t
	}

	if math.Abs(b-a) > 180 {
		if a < b {
			a += 360
		} else {
			b += 360
		}
	}
	result := a + (b-a)*t
	return math.Mod(result, 360)
}

func adjustMultiSteps(hslColors []hsl, totalLen int) int {
	if len(hslColors) < 2 {
		return totalLen
	}

	totalHueDiff := 0.0
	for i := 0; i < len(hslColors)-1; i++ {
		hueDiff := math.Min(
			math.Abs(hslColors[i].h-hslColors[i+1].h),
			360-math.Abs(hslColors[i].h-hslColors[i+1].h),
		)
		totalHueDiff += hueDiff
	}

	if totalHueDiff < 60 {
		return totalLen / (len(hslColors) - 1)
	}

	segmentLens := make([]int, len(hslColors)-1)
	remaining := totalLen

	for i := 0; i < len(hslColors)-2; i++ {
		hueDiff := math.Min(
			math.Abs(hslColors[i].h-hslColors[i+1].h),
			360-math.Abs(hslColors[i].h-hslColors[i+1].h),
		)
		segmentLens[i] = int(float64(totalLen) * hueDiff / totalHueDiff)
		remaining -= segmentLens[i]
	}
	segmentLens[len(hslColors)-2] = remaining

	minLen := totalLen / (3 * (len(hslColors) - 1))
	for i, l := range segmentLens {
		if l < minLen {
			segmentLens[i] = minLen
		}
	}

	sum := 0
	for _, l := range segmentLens {
		sum += l
	}
	return sum / len(segmentLens)
}
