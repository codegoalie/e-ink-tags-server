package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/codegoale/e-ink-tag-server/db"
	"github.com/codegoale/e-ink-tag-server/motivation"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Initialize database
	database, err := db.New(db.GetDBPath())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	e.GET("/disney-countdown/:days", disneyCountdownHandler)
	e.GET("/motivation", motivation.Handler(assets, database))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
