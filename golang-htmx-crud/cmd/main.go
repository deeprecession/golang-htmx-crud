package main

import (
	"context"
	"html/template"
	"io"
	"log/slog"
	"os"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/deeprecession/golang-htmx-crud/pkg/db"
	"github.com/deeprecession/golang-htmx-crud/pkg/handlers"
	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

type Templates struct {
	templates *template.Template
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	log := slog.Default()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				log.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	db_connection, err := db.CreateDatabase(log)
	if err != nil {
		log.Error("Failed to create a db connecdtion: ", err)
		os.Exit(1)
	}
	defer func() {
		if err := db_connection.Close(); err != nil {
			log.Error("failed to close database connection:", err)
		}
	}()

	tasklist, err := models.GetTaskList(db_connection, log)
	if err != nil {
		log.Error("Failed to get tasklist: %q", err)
		os.Exit(1)
	}
	page := models.NewPage(tasklist)
	baseTaskHandler := handlers.NewBaseTaskHandler(&tasklist, &page, log)

	e.GET("/", func(c echo.Context) error {
		tasklist, err := models.GetTaskList(db_connection, log)
		if err != nil {
			log.Error("Failed to get tasklist: %q", err)
			os.Exit(1)
		}

		page := models.NewPage(tasklist)

		return c.Render(200, "index", page)
	})
	e.PUT("/task/:id", baseTaskHandler.ToggleDoneStatusTaskHandler)
	e.DELETE("/task/:id", baseTaskHandler.RemoveTaskHandler)
	e.POST("/tasks", baseTaskHandler.CreateTaskHandler)

	e.Use(echoprometheus.NewMiddleware("myapp"))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.Logger.Fatal(e.Start(":42069"))
}
