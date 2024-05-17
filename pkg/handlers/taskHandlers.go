package handlers

import (
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"

	"gitlab.pg.innopolis.university/v.kishkovskiy/htmx-golang-crud/pkg/models"
)

type UserStorage interface {
	SetDoneStatus(id int, is_done bool) error
	RemoveTask(id int) error
	NewTask(taskTitle string, isDone bool) (models.Task, error)
	HasTask(taskTitle string) (bool, error)
	GetTask(taskId int) (models.Task, error)
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

func (h BaseHandler) ToggleDoneStatusTaskHandler(c echo.Context) error {
	log := c.Logger()

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid id", "err", err)
		return c.String(400, "Invalid id")
	}

	log.Info("PUT /task/:id", "id", taskId)

	task, err := h.storage.GetTask(taskId)
	if err != nil {
		h.log.Error("Task not found", "err", err)
		return c.String(404, "Task is not found")
	}

	newDoneStatus := !task.IsDone

	err = h.storage.SetDoneStatus(taskId, newDoneStatus)
	if err != nil {
		h.log.Error("Task not found", "err", err)
		return c.String(404, "Task is not found")
	}

	updatedTask, err := h.storage.GetTask(taskId)
	if err != nil {
		h.log.Error("failed to get a task", "err", err)
		return c.String(404, "Task is not found")
	}

	return c.Render(200, "task", updatedTask)
}

func (h BaseHandler) RemoveTaskHandler(c echo.Context) error {
	log := c.Logger()

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("Invalid id", "err", err)
		return c.String(400, "Invalid id")
	}

	log.Info("DELETE /task/:id", "id", taskId)

	err = h.storage.RemoveTask(taskId)
	if err != nil {
		h.log.Error("Task not found", "err", err)
		return c.String(404, "Task is not found")
	}

	return c.NoContent(200)
}

func (h BaseHandler) CreateTaskHandler(c echo.Context) error {
	log := c.Logger()

	taskTitle := c.FormValue("title")
	isDone := false

	task, err := h.storage.NewTask(taskTitle, isDone)
	if err != nil {
		h.log.Error("Failed to create a task", "err", err)
		return c.String(500, "Failed to create a task")
	}

	err = c.Render(200, "create-task-form", h.pageCreator.NewFormData())
	if err != nil {
		h.log.Error("Failed to create a form", "err", err)
		return c.String(500, "Failed to create a task")
	}

	log.Info("POST /tasks")

	return c.Render(200, "oob-task", task)
}
