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
		Level: slog.LevelDebug,
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

		reconnectSecondsTime := 5
		time.Sleep(time.Duration(reconnectSecondsTime) * time.Second)

		dbConnection, err = db.CreatePostgresDatabase(psqlInfo)
	}

	defer func() {
		if err := dbConnection.Close(); err != nil {
			log.Error("failed to close database connection", "err", err)
		}
	}()

	tasksStorage := models.NewTaskList(dbConnection, log)

	userStorage := models.GetUserStorage(log, dbConnection)

	httpServer.GET("/", handlers.BaseHandler(&tasksStorage, log))
	httpServer.PUT("/task/:id", handlers.ToggleDoneStatusTaskHandler(&tasksStorage, log))
	httpServer.DELETE("/task/:id", handlers.RemoveTaskHandler(&tasksStorage, log))
	httpServer.POST("/tasks", handlers.CreateTaskHandler(&tasksStorage, log))

	httpServer.Use(echoprometheus.NewMiddleware("myapp"))
	httpServer.GET("/metrics", echoprometheus.NewHandler())

	httpServer.GET("/register", handlers.RegisterPageHandler(log))
	httpServer.POST("/register", handlers.RegisterUserHandler(&userStorage, log))

	httpServer.GET("/login", handlers.LoginPageHandler(log))
	httpServer.POST("/login", handlers.LoginUserHandler(&userStorage, log))

	appPort := os.Getenv("APP_PORT")

	httpServer.Logger.Fatal(httpServer.Start(":" + appPort))
}
