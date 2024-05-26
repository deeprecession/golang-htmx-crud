package handlers

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

type UserStorage interface {
	Register(login string, password string) error
	Login(login string, password string) error
}

func LoginPageHandler(log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Info("GET /login")

		return ctx.Render(OkResponse, "login", nil)
	}
}

func LoginUserHandler(userStorage UserStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		log.Info("POST /login", "login", login, "password", password)

		err := userStorage.Login(login, password)
		if err != nil {
			return ctx.String(BadRequestError, "failed to login: "+err.Error())
		}

		return ctx.String(OkResponse, "logged!")
	}
}

func RegisterPageHandler(log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Info("GET /register")

		return ctx.Render(OkResponse, "register", nil)
	}
}

func RegisterUserHandler(userStorage UserStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login := ctx.FormValue("login")
		password := ctx.FormValue("password")

		log.Info("POST /register", "login", login, "password", password)

		err := userStorage.Register(login, password)
		if err != nil {
			return ctx.String(BadRequestError, "failed to register: "+err.Error())
		}

		return ctx.String(OkResponse, "registered!")
	}
}
