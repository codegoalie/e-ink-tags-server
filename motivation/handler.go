package motivation

import (
	"bytes"
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/codegoale/e-ink-tag-server/db"
	"github.com/fogleman/gg"
	"github.com/labstack/echo/v4"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

// fitText chooses the largest font size (between minFontSize and maxFontSize)
// at which the text, wrapped to maxWidth, fits within both maxWidth and
// maxHeight. It returns the wrapped lines, the chosen font size, and the line
// height. The font size is driven by the longest line so that long, unbreakable
// words are not clipped horizontally.
func fitText(
	dc *gg.Context,
	font *sfnt.Font,
	text string,
	maxWidth, maxHeight, maxFontSize, minFontSize float64,
) (lines []string, fontSize, lineHeight float64) {
	fontSize = maxFontSize

	for fontSize >= minFontSize {
		face, err := opentype.NewFace(
			font,
			&opentype.FaceOptions{Size: fontSize, DPI: 72},
		)
		if err != nil {
			log.Fatal(err)
		}
		dc.SetFontFace(face)

		// Wrap text to fit within maxWidth
		lines = dc.WordWrap(text, maxWidth)

		// Measure a sample character for line height
		_, lineHeight = dc.MeasureString("M")
		totalHeight := float64(len(lines)) * lineHeight

		// Find the width of the longest rendered line so we shrink the font
		// when a single word is too wide to wrap.
		var widestLine float64
		for _, line := range lines {
			lineWidth, _ := dc.MeasureString(line)
			if lineWidth > widestLine {
				widestLine = lineWidth
			}
		}

		// Check if it fits within both maxWidth and maxHeight
		if widestLine <= maxWidth && totalHeight <= maxHeight {
			break
		}

		// Try smaller font
		fontSize -= 2.0
	}

	return lines, fontSize, lineHeight
}

func RenderText(text string, assets embed.FS) (*bytes.Buffer, error) {
	dc := gg.NewContext(296, 128)

	// white background
	height := 128.0
	width := 296.0
	dc.DrawRectangle(0, 0, width, height)
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// motivation text
	dc.SetRGB(0, 0, 0)

	metropolisTTF, err := assets.ReadFile("assets/Metropolis-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	metropolisFont, err := opentype.Parse(metropolisTTF)
	if err != nil {
		log.Fatal(err)
	}

	// Add padding to avoid text touching edges
	padding := 10.0
	maxWidth := width - (2 * padding)
	maxHeight := height - (2 * padding)

	// Find the largest font size that fits, driven by the longest line so long
	// words are not clipped.
	lines, _, lineHeight := fitText(dc, metropolisFont, text, maxWidth, maxHeight, 64.0, 12.0)
	totalHeight := float64(len(lines)) * lineHeight

	// Calculate starting Y position to center text vertically
	startY := (height - totalHeight - 10) / 2

	// Draw each line centered
	for i, line := range lines {
		y := startY + float64(i)*lineHeight + lineHeight/2
		dc.DrawStringAnchored(line, width/2, y, 0.5, 0.5)
	}

	imgBuf := bytes.Buffer{}
	err = dc.EncodePNG(&imgBuf)
	if err != nil {
		return nil, err
	}

	return &imgBuf, nil
}

func Handler(assets embed.FS, database *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		motivation, err := database.GetRandom()
		if err != nil {
			if strings.Contains(err.Error(), "no motivations found") {
				return c.String(http.StatusNotFound, "No motivations found")
			}
			return c.String(http.StatusInternalServerError, "Error retrieving motivation")
		}

		imgBuf, err := RenderText(motivation, assets)
		if err != nil {
			return err
		}

		return c.Stream(http.StatusOK, "image/png", imgBuf)
	}
}
