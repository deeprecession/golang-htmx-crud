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
	server := echo.New()
	server.Renderer = newTemplate()

	slogHandlerOptions := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	stdoutTextHandler := slog.NewTextHandler(os.Stdout, &slogHandlerOptions)
	log := slog.New(stdoutTextHandler)

	server.Use(handlers.NewLoggerMiddleware(log))

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
	sessionStorage := models.NewSessionStore()

	baseGroup := server.Group("")

	baseGroup.Use(echoprometheus.NewMiddleware("myapp"))
	baseGroup.GET("/metrics", echoprometheus.NewHandler())

	authRequiredBaseGroup := baseGroup.Group("")
	authRequiredBaseGroup.Use(handlers.AuthorizationCheckMiddleware(&sessionStorage, log))

	authRequiredBaseGroup.GET("/", handlers.BaseHandler(&tasksStorage, log))
	authRequiredBaseGroup.PUT("/task/:id", handlers.ToggleDoneStatusTaskHandler(&tasksStorage, log))
	authRequiredBaseGroup.DELETE("/task/:id", handlers.RemoveTaskHandler(&tasksStorage, log))
	authRequiredBaseGroup.POST("/tasks", handlers.CreateTaskHandler(&tasksStorage, log))

	baseGroup.GET("/register", handlers.RegisterPageHandler(log))
	baseGroup.POST("/register", handlers.RegisterUserHandler(&userStorage, log))

	baseGroup.GET("/login", handlers.LoginPageHandler(log))
	baseGroup.POST("/login", handlers.LoginUserHandler(&sessionStorage, &userStorage, log))

	appPort := os.Getenv("APP_PORT")

	server.Logger.Fatal(server.Start(":" + appPort))
}
