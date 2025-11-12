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

		// countdown text
		dc.SetRGB(0, 0, 0)

		metropolisTTF, err := assets.ReadFile("assets/Metropolis-Regular.ttf")
		if err != nil {
			log.Fatal(err)
		}
		metropolisFont, err := opentype.Parse(metropolisTTF)
		if err != nil {
			log.Fatal(err)
		}
		daysFace, err := opentype.NewFace(
			metropolisFont,
			&opentype.FaceOptions{Size: 64, DPI: 72},
		)
		if err != nil {
			log.Fatal(err)
		}
		dc.SetFontFace(daysFace)
		dc.DrawStringAnchored(motivation, 180, height/2, 0.5, 0.5)

		imgBuf := bytes.Buffer{}
		err = dc.EncodePNG(&imgBuf)
		if err != nil {
			return err
		}

		return c.Stream(http.StatusOK, "image/png", &imgBuf)
	}
}
