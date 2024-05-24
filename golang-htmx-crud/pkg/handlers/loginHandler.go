package handlers

import "github.com/labstack/echo/v4"

func (h BaseHandler) LoginPageHandler(ctx echo.Context) error {
	h.log.Info("GET /login")

	return ctx.Render(OkResponse, "login", nil)
}

func (h BaseHandler) LoginUserHandler(ctx echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	h.log.Info("POST /login", "login", login, "password", password)

	err := h.userStorage.Login(login, password)
	if err != nil {
		return ctx.String(BadRequestError, "failed to login: "+err.Error())
	}

	return ctx.String(OkResponse, "logged!")
}
