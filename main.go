package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/image/font/opentype"
)

//go:embed assets
var assets embed.FS

func main() {
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.GET("/disney-countdown/:days", func(c echo.Context) error {
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
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
