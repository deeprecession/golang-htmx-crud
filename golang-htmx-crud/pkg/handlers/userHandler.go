package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAuth interface {
	Register(login string, password string) error
	Login(login string, password string) error
}

type SessionStore interface {
	GetSession(response *http.Request, key string) (string, error)
	SetSession(response *http.ResponseWriter, key string, val string) error
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

type LoginFormResponse struct {
	LoginValue    string
	PasswordValue string
	Error         string
}

func LoginUserHandler(
	sessionStorage SessionStore,
	userAuth UserAuth,
	log *slog.Logger,
) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		err := userAuth.Login(login, password)
		if err != nil {

			log.Debug(
				"POST /login failed to login",
				"login",
				login,
				"password",
				password,
				"err",
				err,
			)

			loginResponse := LoginFormResponse{
				LoginValue:    login,
				PasswordValue: password,
				Error:         err.Error(),
			}

			return ctx.Render(http.StatusOK, "login", loginResponse)
		}

		err = sessionStorage.SetSession(&ctx.Response().Writer, "session", login)
		if err != nil {
			log.Debug(
				"POST /login failed to login",
				"login",
				login,
				"password",
				password,
				"err",
				err,
			)

			loginResponse := LoginFormResponse{
				LoginValue:    login,
				PasswordValue: password,
				Error:         err.Error(),
			}

			return ctx.Render(http.StatusOK, "login", loginResponse)
		}

		log.Debug("POST /login logged successfully", "login", login, "password", password)

		return ctx.Redirect(http.StatusFound, "/")
	}
}

func RegisterPageHandler(log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Info("GET /register")

		return ctx.Render(http.StatusOK, "register", nil)
	}
}

type RegisterFormResponse struct {
	LoginValue    string
	PasswordValue string
	Error         string
}

func RegisterUserHandler(userAuth UserAuth, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		log.Info("POST /register", "login", login, "password", password)

		err := userAuth.Register(login, password)
		if err != nil {
			registerResponse := RegisterFormResponse{
				login,
				password,
				err.Error(),
			}

			return ctx.Render(http.StatusOK, "register", registerResponse)
		}

		return ctx.Redirect(http.StatusFound, "/")
	}
}
