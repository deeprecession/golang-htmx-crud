package handlers

import (
	"errors"
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

func (h BaseHandler) ToggleDoneStatusTaskHandler(ctx echo.Context) error {
	log := ctx.Logger()

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.log.Error("Invalid id", "err", err)

		return ctx.String(BadRequestError, "Invalid id")
	}

	log.Info("PUT /task/:id", "id", taskID)

	task, err := h.taskStorage.GetTaskByID(taskID)
	if err != nil {
		h.log.Error("Task not found", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	newDoneStatus := !task.IsDone

	err = h.taskStorage.SetDoneStatus(taskID, newDoneStatus)
	if err != nil {
		h.log.Error("Task not found", "err", err)

		return ctx.String(NotFoundError, "Task is not found")
	}

	updatedTask, err := h.taskStorage.GetTaskByID(taskID)
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

	err = h.taskStorage.RemoveTask(taskID)
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

	task, err := h.taskStorage.NewTask(taskTitle, isDone)
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
