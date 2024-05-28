package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
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

		return ctx.Render(http.StatusOK, "index", page)
	}
}
