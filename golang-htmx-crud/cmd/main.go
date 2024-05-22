package main

import (
	"html/template"
	"io"
	"log/slog"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

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

func (t *Templates) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	httpServer := echo.New()
	httpServer.Renderer = newTemplate()

	log := slog.Default()

	httpServer.Use(handlers.NewLoggerMiddleware(log))

	const sqlitePath = "./sqlite.db"

	dbConneciton, err := db.CreateSQLiteDatabase(sqlitePath)
	if err != nil {
		log.Error("Failed to create a db connecdtion", "err", err)

		return
	}

	defer func() {
		if err := dbConneciton.Close(); err != nil {
			log.Error("failed to close database connection", "err", err)
		}
	}()

	tasklist, err := models.GetTaskList(dbConneciton, log)
	if err != nil {
		log.Error("failed to get tasklist:", "err", err)

		return
	}

	page := models.NewPage(tasklist)
	baseTaskHandler := handlers.NewBaseTaskHandler(&tasklist, &page, log)

	httpServer.GET("/", func(ctx echo.Context) error {
		tasklist, err := models.GetTaskList(dbConneciton, log)
		if err != nil {
			log.Error("failed to get tasklist:", "err", err)
		}

		page := models.NewPage(tasklist)

		return ctx.Render(handlers.OkResponse, "index", page)
	})
	httpServer.PUT("/task/:id", baseTaskHandler.ToggleDoneStatusTaskHandler)
	httpServer.DELETE("/task/:id", baseTaskHandler.RemoveTaskHandler)
	httpServer.POST("/tasks", baseTaskHandler.CreateTaskHandler)

	httpServer.Use(echoprometheus.NewMiddleware("myapp"))
	httpServer.GET("/metrics", echoprometheus.NewHandler())

	httpServer.Logger.Fatal(httpServer.Start(":42069"))
}
