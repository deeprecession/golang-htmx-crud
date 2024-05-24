package main

import (
	"html/template"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
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

func (t *Templates) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	httpServer := echo.New()
	httpServer.Renderer = newTemplate()

	slogHandlerOptions := slog.HandlerOptions{
		Level: slog.Level(slog.LevelDebug),
	}
	stdoutTextHandler := slog.NewTextHandler(os.Stdout, &slogHandlerOptions)
	log := slog.New(stdoutTextHandler)

	httpServer.Use(handlers.NewLoggerMiddleware(log))

	psqlInfo, err := db.GetPsqlInfoFromEnv()
	if err != nil {
		log.Error("failed to get posqlInfo:", "err", err)

		return
	}

	dbConnection, err := db.CreatePostgresDatabase(psqlInfo)

	for err != nil {
		log.Error("failed to connect to a database:", "psqlInfo", psqlInfo, "err", err)

		time.Sleep(time.Second * 5)

		dbConnection, err = db.CreatePostgresDatabase(psqlInfo)
	}

	defer func() {
		if err := dbConnection.Close(); err != nil {
			log.Error("failed to close database connection", "err", err)
		}
	}()

	tasklist, err := models.GetTaskList(dbConnection, log)
	if err != nil {
		log.Error("failed to get tasklist:", "err", err)

		return
	}

	userStorage := models.GetUserStorage(log, dbConnection)

	page := models.NewPage(tasklist)
	baseTaskHandler := handlers.NewBaseTaskHandler(&tasklist, &page, &userStorage, log)

	httpServer.GET("/", func(ctx echo.Context) error {
		tasklist, err := models.GetTaskList(dbConnection, log)
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

	httpServer.GET("/register", baseTaskHandler.RegisterPageHandler)
	httpServer.POST("/register", baseTaskHandler.RegisterUserHandler)

	httpServer.GET("/login", baseTaskHandler.LoginPageHandler)
	httpServer.POST("/login", baseTaskHandler.LoginUserHandler)

	appPort := os.Getenv("APP_PORT")

	httpServer.Logger.Fatal(httpServer.Start(":" + appPort))
}
