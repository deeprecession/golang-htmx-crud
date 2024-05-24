package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h BaseHandler) RegisterPageHandler(ctx echo.Context) error {
	h.log.Info("GET /register")

	return ctx.Render(OkResponse, "register", nil)
}

func (h BaseHandler) RegisterUserHandler(ctx echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	h.log.Info("POST /register", "login", login, "password", password)

	err := h.userStorage.Register(login, password)
	if err != nil {
		return ctx.String(BadRequestError, "failed to register: "+err.Error())
	}

	return ctx.String(OkResponse, "registered!")
}
