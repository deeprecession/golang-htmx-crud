package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
)

type Tasks []Task

type Task struct {
	ID     int
	Title  string
	IsDone bool
}

type TaskList struct {
	Tasks Tasks
	log   *slog.Logger
	db    *sql.DB
}

func GetTaskList(database *sql.DB, logger *slog.Logger) (TaskList, error) {
	const funcErrMsg = "models.GetTaskList"

	rows, err := database.Query("SELECT * FROM tasks")
	if err != nil {
		return TaskList{}, fmt.Errorf("%s: failed to query tasks table: %w", funcErrMsg, err)
	}

	defer rows.Close()

	tasks := Tasks{}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Title, &task.IsDone)
		if err != nil {
			return TaskList{}, fmt.Errorf("%s: failed to scan rows: %w", funcErrMsg, err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return TaskList{}, fmt.Errorf("%s: rows error: %w", funcErrMsg, err)
	}

	tasklist := TaskList{
		tasks,
		logger,
		database,
	}

	return tasklist, nil
}

func (tl *TaskList) NewTask(title string, isDone bool) (Task, error) {
	const funcErrMsg = "models.NewTask"

	isTaskExist, err := tl.HasTask(title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to check is task exist: %w", funcErrMsg, err)
	}

	if isTaskExist {
		return Task{}, ErrTaskAlreadyExist
	}

	stmt, err := tl.db.Prepare("INSERT INTO tasks(title, is_done) VALUES ($1, $2)")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(title, isDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a statement: %w", funcErrMsg, err)
	}

	task, err := tl.GetTaskByTitle(title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to get a task by title %w", funcErrMsg, err)
	}

	tl.log.Info("Succsesfully inserted:", "id", task.ID, "title", task.Title, "isDone", task.IsDone)

	return task, nil
}

func (tl *TaskList) HasTask(title string) (bool, error) {
	const funcErrMsg = "models.HasTask"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE title = $1")
	if err != nil {
		return false, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	result, err := stmt.Query(title)
	if err != nil {
		return false, fmt.Errorf("%s: failed to execute a query: %w", funcErrMsg, err)
	}

	defer result.Close()

	hasTask := result.Next()

	if err = result.Err(); err != nil {
		return false, fmt.Errorf("%s: failed to get next value: %w", funcErrMsg, err)
	}

	return hasTask, nil
}

func (tl *TaskList) RemoveTask(taskID int) error {
	const funcErrMsg = "models.RemoveTask"

	stmt, err := tl.db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskID)
	if err != nil {
		return fmt.Errorf("%s: failed to execute a query: %w", funcErrMsg, err)
	}

	return nil
}

func (tl *TaskList) GetTaskByID(taskID int) (Task, error) {
	const funcErrMsg = "models.GetTask"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE id = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	task, err := tl.GetTask(stmt, strconv.Itoa(taskID))
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to get a task: %w", funcErrMsg, err)
	}

	return task, nil
}

func (tl *TaskList) GetTaskByTitle(title string) (Task, error) {
	const funcErrMsg = "models.GetTaskByTitle"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE title = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	task, err := tl.GetTask(stmt, title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to get a task: %w", funcErrMsg, err)
	}

	return task, nil
}

func (tl *TaskList) GetTask(stmt *sql.Stmt, title string) (Task, error) {
	const funcErrMsg = "models.GetTask"

	result, err := stmt.Query(title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a query: %w", funcErrMsg, err)
	}

	defer result.Close()

	hasNext := result.Next()
	if !hasNext {
		return Task{}, fmt.Errorf("%s: task is not found", funcErrMsg)
	}

	if err = result.Err(); err != nil {
		return Task{}, fmt.Errorf("%s: task is not found: %w", funcErrMsg, err)
	}

	task := Task{}

	err = result.Scan(&task.ID, &task.Title, &task.IsDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to scan a query response: %w", funcErrMsg, err)
	}

	return task, nil
}

func (tl *TaskList) SetDoneStatus(taskID int, isDone bool) error {
	const funcErrMsg = "models.ToggleDoneStatus"

	stmt, err := tl.db.Prepare("UPDATE tasks SET is_done = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(isDone, taskID)
	if err != nil {
		return fmt.Errorf("%s: failed to execute a query: %w", funcErrMsg, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", funcErrMsg, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows affected", funcErrMsg)
	}

	return nil
}
