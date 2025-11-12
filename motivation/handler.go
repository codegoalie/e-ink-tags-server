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
)

func Handler(assets embed.FS, database *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		motivation, err := database.GetRandom()
		if err != nil {
			if strings.Contains(err.Error(), "no motivations found") {
				return c.String(http.StatusNotFound, "No motivations found")
			}
			return c.String(http.StatusInternalServerError, "Error retrieving motivation")
		}

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

		// Try progressively smaller font sizes until text fits
		fontSize := 64.0
		minFontSize := 12.0
		var lines []string
		var totalHeight float64

		for fontSize >= minFontSize {
			face, err := opentype.NewFace(
				metropolisFont,
				&opentype.FaceOptions{Size: fontSize, DPI: 72},
			)
			if err != nil {
				log.Fatal(err)
			}
			dc.SetFontFace(face)

			// Wrap text to fit within maxWidth
			lines = dc.WordWrap(motivation, maxWidth)

			// Calculate total height needed for all lines
			_, lineHeight := dc.MeasureString("M") // Measure a sample character for line height
			totalHeight = float64(len(lines)) * lineHeight

			// Check if it fits within maxHeight
			if totalHeight <= maxHeight {
				break
			}

			// Try smaller font
			fontSize -= 2.0
		}

		// Calculate starting Y position to center text vertically
		_, lineHeight := dc.MeasureString("M")
		startY := (height - totalHeight) / 2

		// Draw each line centered
		for i, line := range lines {
			y := startY + float64(i)*lineHeight + lineHeight/2
			dc.DrawStringAnchored(line, width/2, y, 0.5, 0.5)
		}

		imgBuf := bytes.Buffer{}
		err = dc.EncodePNG(&imgBuf)
		if err != nil {
			return err
		}

		return c.Stream(http.StatusOK, "image/png", &imgBuf)
	}
}
