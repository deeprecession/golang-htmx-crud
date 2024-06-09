package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

type UserStorage interface {
	GetUserWithLogin(login string) (models.User, error)
}

func ToggleDoneStatusTaskHandler(
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
			log.Error("failed to get a user with login", "err", err)

			return ctx.String(http.StatusInternalServerError, "Invalid id")
		}

		taskID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			log.Error("Invalid id", "err", err)

			return ctx.String(http.StatusBadRequest, "Invalid id")
		}

		log.Info("PUT /task/:id", "id", taskID)

		task, err := user.GetTaskByID(taskID)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(http.StatusNotFound, "Task is not found")
		}

		newDoneStatus := !task.IsDone

		err = user.SetDoneStatus(taskID, newDoneStatus)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(http.StatusNotFound, "Task is not found")
		}

		updatedTask, err := user.GetTaskByID(taskID)
		if err != nil {
			log.Error("failed to get a task", "err", err)

			return ctx.String(http.StatusNotFound, "Task is not found")
		}

		return ctx.Render(http.StatusOK, "task", updatedTask)
	}
}

func RemoveTaskHandler(
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
			log.Error("failed to get a user with login", "err", err)

			return ctx.String(http.StatusInternalServerError, "Invalid id")
		}

		taskID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			log.Error("Invalid id", "err", err)

			return ctx.String(http.StatusBadRequest, "Invalid id")
		}

		log.Info("DELETE /task/:id", "id", taskID)

		err = user.RemoveTask(taskID)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(http.StatusNotFound, "Task is not found")
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func CreateTaskHandler(
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
			log.Error("failed to get a user with login", "err", err)

			return ctx.String(http.StatusInternalServerError, "Invalid id")
		}

		taskTitle := ctx.FormValue("title")
		isDone := false

		task, err := user.NewTask(taskTitle, isDone)
		if errors.Is(err, models.ErrTaskAlreadyExist) {
			newFormData := models.NewFormData()
			newFormData.Values["Title"] = taskTitle
			newFormData.Errors["Title"] = "Task already exist"

			return ctx.Render(http.StatusOK, "create-task-form", newFormData)
		}

		if err != nil {
			log.Error("Failed to create a task", "err", err)

			return ctx.String(http.StatusInternalServerError, "Failed to create a task")
		}

		err = ctx.Render(http.StatusOK, "create-task-form", models.NewFormData())
		if err != nil {
			log.Error("Failed to create a form", "err", err)

			return ctx.String(http.StatusInternalServerError, "Failed to create a task")
		}

		log.Info("POST /tasks")

		return ctx.Render(http.StatusOK, "oob-task", task)
	}
}
