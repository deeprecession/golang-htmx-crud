package handlers

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

const (
	BadRequestError     = 400
	NotFoundError       = 404
	OkResponse          = 200
	InternalServerError = 500
)

type UserStorage interface {
	SetDoneStatus(id int, isDone bool) error
	RemoveTask(id int) error
	NewTask(taskTitle string, isDone bool) (models.Task, error)
	HasTask(taskTitle string) (bool, error)
	GetTaskByID(taskID int) (models.Task, error)
}

type PageCreator interface {
	NewFormData() models.FormData
}

type BaseHandler struct {
	storage     UserStorage
	pageCreator PageCreator
	log         *slog.Logger
}

func NewBaseTaskHandler(
	userStorage UserStorage,
	pageCreator PageCreator,
	logger *slog.Logger,
) BaseHandler {
	return BaseHandler{userStorage, pageCreator, logger}
}

func (h BaseHandler) ToggleDoneStatusTaskHandler(ctx echo.Context) error {
	log := ctx.Logger()

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.log.Error("Invalid id", "err", err)

		return ctx.String(BadRequestError, "Invalid id")
	}

	log.Info("PUT /task/:id", "id", taskID)

	task, err := h.storage.GetTaskByID(taskID)
	if err != nil {
		h.log.Error("Task not found", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	newDoneStatus := !task.IsDone

	err = h.storage.SetDoneStatus(taskID, newDoneStatus)
	if err != nil {
		h.log.Error("Task not found", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	updatedTask, err := h.storage.GetTaskByID(taskID)
	if err != nil {
		h.log.Error("failed to get a task", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	return ctx.Render(OkResponse, "task", updatedTask)
}

func (h BaseHandler) RemoveTaskHandler(ctx echo.Context) error {
	log := ctx.Logger()

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.log.Error("Invalid id", "err", err)

		return ctx.String(BadRequestError, "Invalid id")
	}

	log.Info("DELETE /task/:id", "id", taskID)

	err = h.storage.RemoveTask(taskID)
	if err != nil {
		h.log.Error("Task not found", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	return ctx.NoContent(OkResponse)
}

func (h BaseHandler) CreateTaskHandler(ctx echo.Context) error {
	log := ctx.Logger()

	taskTitle := ctx.FormValue("title")
	isDone := false

	task, err := h.storage.NewTask(taskTitle, isDone)
	if errors.Is(err, models.ErrTaskAlreadyExist) {
		newFormData := h.pageCreator.NewFormData()
		newFormData.Values["Title"] = taskTitle
		newFormData.Errors["Title"] = "Task already exist"

		return ctx.Render(OkResponse, "create-task-form", newFormData)
	}

	if err != nil {
		h.log.Error("Failed to create a task", "err", err)

		return ctx.String(InternalServerError, "Failed to create a task")
	}

	err = ctx.Render(OkResponse, "create-task-form", h.pageCreator.NewFormData())
	if err != nil {
		h.log.Error("Failed to create a form", "err", err)

		return ctx.String(InternalServerError, "Failed to create a task")
	}

	log.Info("POST /tasks")

	return ctx.Render(OkResponse, "oob-task", task)
}
