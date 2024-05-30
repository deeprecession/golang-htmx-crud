package main

import (
	"database/sql"
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

func main() {
	log := getLogger()

	dbCon := getDB(log)
	defer func() {
		if err := dbCon.Close(); err != nil {
			log.Error("failed to close database connection", "err", err)
		}
	}()

	templates := newTemplate()

	startServer(log, dbCon, templates)
}

type Templates struct {
	templates *template.Template
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
}

func getDB(log *slog.Logger) *sql.DB {
	psqlInfo, err := db.GetPsqlInfoFromEnv()
	if err != nil {
		log.Error("failed to get posqlInfo:", "err", err)

		os.Exit(1)
	}

	dbConnection, err := db.CreatePostgresDatabase(psqlInfo)

	for err != nil {
		log.Error("failed to connect to a database:", "psqlInfo", psqlInfo, "err", err)

		reconnectSecondsTime := 5
		time.Sleep(time.Duration(reconnectSecondsTime) * time.Second)

		dbConnection, err = db.CreatePostgresDatabase(psqlInfo)
	}

	return dbConnection
}

func getLogger() *slog.Logger {
	slogHandlerOptions := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	stdoutTextHandler := slog.NewTextHandler(os.Stdout, &slogHandlerOptions)

	return slog.New(stdoutTextHandler)
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func startServer(log *slog.Logger, dbCon *sql.DB, templates *Templates) {
	server := echo.New()
	server.Renderer = templates

	server.Use(handlers.NewLoggerMiddleware(log))

	userStorage := models.GetUserStorage(log, dbCon)
	sessionStorage := models.NewSessionStore()

	baseGroup := server.Group("")

	baseGroup.Use(echoprometheus.NewMiddleware("myapp"))
	baseGroup.GET("/metrics", echoprometheus.NewHandler())

	authRequiredBaseGroup := baseGroup.Group("")
	authRequiredBaseGroup.Use(handlers.AuthorizationCheckMiddleware(&sessionStorage, log))

	authRequiredBaseGroup.GET("/", handlers.BaseHandler(&sessionStorage, &userStorage, log))
	authRequiredBaseGroup.PUT(
		"/task/:id",
		handlers.ToggleDoneStatusTaskHandler(&sessionStorage, &userStorage, log),
	)
	authRequiredBaseGroup.DELETE(
		"/task/:id",
		handlers.RemoveTaskHandler(&sessionStorage, &userStorage, log),
	)
	authRequiredBaseGroup.POST(
		"/tasks",
		handlers.CreateTaskHandler(&sessionStorage, &userStorage, log),
	)

	baseGroup.GET("/register", handlers.RegisterPageHandler(log))
	baseGroup.POST("/register", handlers.RegisterUserHandler(&userStorage, log))

	baseGroup.GET("/login", handlers.LoginPageHandler(log))
	baseGroup.POST("/login", handlers.LoginUserHandler(&sessionStorage, &userStorage, log))

	appPort := os.Getenv("APP_PORT")

	server.Logger.Fatal(server.Start(":" + appPort))
}
