package handlers

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

type TaskStorage interface {
	SetDoneStatus(id int, isDone bool) error
	GetTasks() (models.Tasks, error)
	RemoveTask(id int) error
	NewTask(taskTitle string, isDone bool) (models.Task, error)
	HasTask(taskTitle string) (bool, error)
	GetTaskByID(taskID int) (models.Task, error)
}

func ToggleDoneStatusTaskHandler(taskStorage TaskStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		taskID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			log.Error("Invalid id", "err", err)

			return ctx.String(BadRequestError, "Invalid id")
		}

		log.Info("PUT /task/:id", "id", taskID)

		task, err := taskStorage.GetTaskByID(taskID)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(NotFoundError, "Task is not found")
		}

		newDoneStatus := !task.IsDone

		err = taskStorage.SetDoneStatus(taskID, newDoneStatus)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(NotFoundError, "Task is not found")
		}

		updatedTask, err := taskStorage.GetTaskByID(taskID)
		if err != nil {
			log.Error("failed to get a task", "err", err)

			return ctx.String(NotFoundError, "Task is not found")
		}

		return ctx.Render(OkResponse, "task", updatedTask)
	}
}

func RemoveTaskHandler(taskStorage TaskStorage, log *slog.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		taskID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			log.Error("Invalid id", "err", err)

			return ctx.String(BadRequestError, "Invalid id")
		}

		log.Info("DELETE /task/:id", "id", taskID)

		err = taskStorage.RemoveTask(taskID)
		if err != nil {
			log.Error("Task not found", "err", err)

			return ctx.String(NotFoundError, "Task is not found")
		}

		return ctx.NoContent(OkResponse)
	}
}

func CreateTaskHandler(
	taskStorage TaskStorage,
	log *slog.Logger,
) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		taskTitle := ctx.FormValue("title")
		isDone := false

		task, err := taskStorage.NewTask(taskTitle, isDone)
		if errors.Is(err, models.ErrTaskAlreadyExist) {
			newFormData := models.NewFormData()
			newFormData.Values["Title"] = taskTitle
			newFormData.Errors["Title"] = "Task already exist"

			return ctx.Render(OkResponse, "create-task-form", newFormData)
		}

		if err != nil {
			log.Error("Failed to create a task", "err", err)

			return ctx.String(InternalServerError, "Failed to create a task")
		}

		err = ctx.Render(OkResponse, "create-task-form", models.NewFormData())
		if err != nil {
			log.Error("Failed to create a form", "err", err)

			return ctx.String(InternalServerError, "Failed to create a task")
		}

		log.Info("POST /tasks")

		return ctx.Render(OkResponse, "oob-task", task)
	}
}
