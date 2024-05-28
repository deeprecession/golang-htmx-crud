package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserStorage interface {
	Register(login string, password string) error
	Login(login string, password string) error
}

type SessionStore interface {
	GetSession(*http.Request, string) (string, error)
	SetSession(*http.ResponseWriter, string, string) error
}

func AuthorizationCheckMiddleware(store SessionStore, log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			_, err := store.GetSession(ctx.Request(), "session")
			if err != nil {
				log.Info("Not authorized! Redirecting...", "err", err)

				return ctx.Redirect(http.StatusFound, "/login")
			}

			return next(ctx)
		}
	}
}

func LoginPageHandler(log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Info("GET /login")

		return ctx.Render(http.StatusOK, "login", nil)
	}
}

func LoginUserHandler(
	sessionStorage SessionStore,
	userStorage UserStorage,
	log *slog.Logger,
) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		log.Info("POST /login", "login", login, "password", password)

		err := userStorage.Login(login, password)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "failed to login: "+err.Error())
		}

		sessionStorage.SetSession(&ctx.Response().Writer, "session", login)

		return ctx.Redirect(http.StatusFound, "/")
	}
}

func RegisterPageHandler(log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Info("GET /register")

		return ctx.Render(http.StatusOK, "register", nil)
	}
}

type RegisterResponse struct {
	LoginValue    string
	PasswordValue string
	Error         string
}

func RegisterUserHandler(userStorage UserStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		log.Info("POST /register", "login", login, "password", password)

		err := userStorage.Register(login, password)
		if err != nil {
			registerResponse := RegisterResponse{
				login,
				password,
				err.Error(),
			}

			return ctx.Render(http.StatusOK, "register", registerResponse)
		}

		return ctx.Redirect(http.StatusFound, "/")
	}
}
