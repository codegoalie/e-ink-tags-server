package motivation

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RenderHandler(assets fs.ReadFileFS) echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")
		if text == "" {
			return c.String(http.StatusBadRequest, "text parameter is required")
		}

		imgBuf, err := RenderText(text, assets)
		if err != nil {
			return err
		}

		return c.Stream(http.StatusOK, "image/png", imgBuf)
	}
}
