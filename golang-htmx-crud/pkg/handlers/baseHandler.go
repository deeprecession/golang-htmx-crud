package handlers

import (
	"log/slog"

	"github.com/deeprecession/golang-htmx-crud/pkg/models"
)

type TaskStorage interface {
	SetDoneStatus(id int, isDone bool) error
	RemoveTask(id int) error
	NewTask(taskTitle string, isDone bool) (models.Task, error)
	HasTask(taskTitle string) (bool, error)
	GetTaskByID(taskID int) (models.Task, error)
}

type PageCreator interface {
	NewFormData() models.FormData
}

type UserStorage interface {
	Register(login string, password string) error
	Login(login string, password string) error
}

type BaseHandler struct {
	taskStorage TaskStorage
	userStorage UserStorage
	pageCreator PageCreator
	log         *slog.Logger
}

func NewBaseTaskHandler(
	taskStorage TaskStorage,
	pageCreator PageCreator,
	userStorage UserStorage,
	logger *slog.Logger,
) BaseHandler {
	return BaseHandler{taskStorage, userStorage, pageCreator, logger}
}
