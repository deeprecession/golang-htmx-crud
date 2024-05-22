package models

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Tasks []Task

type Task struct {
	Id     int
	Title  string
	IsDone bool
}

type TaskList struct {
	Tasks Tasks
	log   *slog.Logger
	db    *sql.DB
}

func GetTaskList(db *sql.DB, logger *slog.Logger) (TaskList, error) {
	const funcErrMsg = "models.GetTaskList"

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return TaskList{}, fmt.Errorf("%s: failed to query tasks table: %q", funcErrMsg, err)
	}

	tasks := Tasks{}
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.Id, &task.Title, &task.IsDone)
		if err != nil {
			return TaskList{}, fmt.Errorf("%s: failed to scan rows: %q", funcErrMsg, err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return TaskList{}, fmt.Errorf("%s: rows error: %q", funcErrMsg, err)
	}

	tasklist := TaskList{
		tasks,
		logger,
		db,
	}
	return tasklist, nil
}

func (tl *TaskList) NewTask(title string, isDone bool) (Task, error) {
	const funcErrMsg = "models.NewTask"

	stmt, err := tl.db.Prepare("INSERT INTO tasks(title, is_done) VALUES ($1, $2)")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	_, err = stmt.Exec(title, isDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a statement: %q", funcErrMsg, err)
	}

	err = stmt.Close()
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to close a stmt: %q", funcErrMsg, err)
	}

	task, err := tl.GetTaskByTitle(title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to get a task by title %q", funcErrMsg, err)
	}

	tl.log.Info("Succsesfully inserted:", "id", task.Id, "title", task.Title, "isDone", task.IsDone)

	return task, nil
}

func (tl *TaskList) HasTask(title string) (bool, error) {
	const funcErrMsg = "models.HasTask"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE title = $1")
	if err != nil {
		return false, fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	result, err := stmt.Query(title)
	if err != nil {
		return false, fmt.Errorf("%s: failed to execute a query: %q", funcErrMsg, err)
	}

	hasTask := result.Next()

	err = result.Close()
	if err != nil {
		return false, fmt.Errorf("%s: failed to close a result: %q", funcErrMsg, err)
	}

	return hasTask, nil
}

func (tl *TaskList) RemoveTask(id int) error {
	const funcErrMsg = "models.RemoveTask"

	stmt, err := tl.db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute a query: %q", funcErrMsg, err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close a stmt: %q", funcErrMsg, err)
	}

	return nil
}

func (tl *TaskList) GetTask(id int) (Task, error) {
	const funcErrMsg = "models.GetTask"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE id = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	result, err := stmt.Query(id)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a query: %q", funcErrMsg, err)
	}

	hasNext := result.Next()
	if !hasNext {
		return Task{}, fmt.Errorf("%s: task with id=%d not found", funcErrMsg, id)
	}

	task := Task{}
	err = result.Scan(&task.Id, &task.Title, &task.IsDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to scan a query respone: %q", funcErrMsg, err)
	}

	err = result.Close()
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to close a result: %q", funcErrMsg, err)
	}

	return task, nil
}

func (tl *TaskList) GetTaskByTitle(title string) (Task, error) {
	const funcErrMsg = "models.GetTaskByTitle"

	stmt, err := tl.db.Prepare("SELECT * FROM tasks WHERE title = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	result, err := stmt.Query(title)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a query: %q", funcErrMsg, err)
	}

	hasNext := result.Next()
	if !hasNext {
		return Task{}, fmt.Errorf("%s: task with id=%s not found", funcErrMsg, title)
	}

	task := Task{}
	err = result.Scan(&task.Id, &task.Title, &task.IsDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to scan a query respone: %q", funcErrMsg, err)
	}

	err = result.Close()
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to close a result: %q", funcErrMsg, err)
	}

	return task, nil
}

func (tl *TaskList) SetDoneStatus(id int, is_done bool) error {
	const funcErrMsg = "models.ToggleDoneStatus"

	stmt, err := tl.db.Prepare("UPDATE tasks SET is_done = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare a statement: %q", funcErrMsg, err)
	}

	result, err := stmt.Exec(is_done, id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute a query: %q", funcErrMsg, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %q", funcErrMsg, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows affected", funcErrMsg)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close a stmt: %q", funcErrMsg, err)
	}

	return nil
}
