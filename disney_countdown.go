package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/fogleman/gg"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/image/font/opentype"
)

func disneyCountdownHandler(c echo.Context) error {
	days := c.Param("days")

	dc := gg.NewContext(296, 128)

	// white background
	height := 128.0
	width := 296.0
	dc.DrawRectangle(0, 0, width, height)
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// castle icon
	castleFile, err := assets.Open("assets/cinderella-castle-icon-thumb.jpg")
	if err != nil {
		err := fmt.Errorf("opening castle icon: %w", err)
		return err
	}
	im, _, err := image.Decode(castleFile)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawImage(im, 0, 0)

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
	dc.DrawStringAnchored(days, 180, height/2, 0.5, 0.5)

	unitFace, err := opentype.NewFace(
		metropolisFont,
		&opentype.FaceOptions{Size: 18, DPI: 72},
	)
	if err != nil {
		log.Fatal(err)
	}
	dc.SetFontFace(unitFace)
	dc.DrawStringAnchored("days", 250, height/2+16, 0.5, 0.5)

	imgBuf := bytes.Buffer{}
	err = dc.EncodePNG(&imgBuf)
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, "image/png", &imgBuf)
}
