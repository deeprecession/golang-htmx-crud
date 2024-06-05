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

func BaseHandler(
	sessionStore SessionStore,
	userStorage UserStorage,
	log *slog.Logger,
) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		login, err := sessionStore.GetSession(ctx.Request(), "session")
		if err != nil {
			log.Info("Not authorized! Redirecting...", "err", err)

			return ctx.Redirect(http.StatusFound, "/login")
		}

		user, err := userStorage.GetUserWithLogin(login)
		if err != nil {
			log.Error("failed to get user by login", "err", err)

			return ctx.String(http.StatusInternalServerError, "failed to get user by login")
		}

		tasklist, err := user.GetTasks()
		if err != nil {
			log.Error("failed to get tasks:", "err", err)
		}

		page := models.NewPage(tasklist, user)

		return ctx.Render(http.StatusOK, "tasklist-page", page)
	}
}
