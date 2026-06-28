package motivation

import (
	"os"
	"testing"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/opentype"
)

func loadFont(t *testing.T) *opentype.Font {
	t.Helper()

	ttf, err := os.ReadFile("../assets/Metropolis-Regular.ttf")
	if err != nil {
		t.Fatalf("reading font: %v", err)
	}
	font, err := opentype.Parse(ttf)
	if err != nil {
		t.Fatalf("parsing font: %v", err)
	}
	return font
}

// widestLine measures the widest of the given lines at the chosen font size.
func widestLine(t *testing.T, dc *gg.Context, font *opentype.Font, fontSize float64, lines []string) float64 {
	t.Helper()

	face, err := opentype.NewFace(font, &opentype.FaceOptions{Size: fontSize, DPI: 72})
	if err != nil {
		t.Fatalf("creating face: %v", err)
	}
	dc.SetFontFace(face)

	var widest float64
	for _, line := range lines {
		w, _ := dc.MeasureString(line)
		if w > widest {
			widest = w
		}
	}
	return widest
}

func TestFitTextLongWordFitsWidth(t *testing.T) {
	font := loadFont(t)
	dc := gg.NewContext(296, 128)

	const (
		maxWidth    = 276.0 // 296 - 2*10 padding
		maxHeight   = 108.0 // 128 - 2*10 padding
		maxFontSize = 64.0
		minFontSize = 12.0
	)

	// "superpower" is a long, unbreakable word that previously clipped.
	lines, fontSize, _ := fitText(dc, font, "Speed is a superpower", maxWidth, maxHeight, maxFontSize, minFontSize)

	widest := widestLine(t, dc, font, fontSize, lines)
	if widest > maxWidth {
		t.Errorf("widest line %.1f exceeds maxWidth %.1f at font size %.1f", widest, maxWidth, fontSize)
	}
}
