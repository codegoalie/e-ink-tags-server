package main

import (
	"bytes"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		dc := gg.NewContext(296, 128)

		// white background
		dc.DrawRectangle(0, 0, 296, 128)
		dc.SetRGB(1, 1, 1)
		dc.Fill()

		// red circle
		dc.DrawCircle(50, 50, 40)
		dc.SetRGB(1, 0, 0)
		dc.Fill()

		imgBuf := bytes.Buffer{}
		err := dc.EncodePNG(&imgBuf)
		if err != nil {
			return err
		}

		return c.Stream(http.StatusOK, "image/png", &imgBuf)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
