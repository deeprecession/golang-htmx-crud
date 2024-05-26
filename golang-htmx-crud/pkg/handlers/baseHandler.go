package handlers

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

const (
	BadRequestError     = 400
	NotFoundError       = 404
	OkResponse          = 200
	InternalServerError = 500
)

type PageCreator interface {
	NewFormData() models.FormData
}

func BaseHandler(tasklist TaskStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tasklist, err := tasklist.GetTasks()
		if err != nil {
			log.Error("failed to get tasks:", "err", err)
		}

		page := models.NewPage(tasklist)

		return ctx.Render(OkResponse, "index", page)
	}
}
