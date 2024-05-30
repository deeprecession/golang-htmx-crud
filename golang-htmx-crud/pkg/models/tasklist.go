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

type User struct {
	id       int
	login    string
	password string
	db       *sql.DB
	log      *slog.Logger
}

func (user *User) GetTasks() (Tasks, error) {
	const funcErrMsg = "models.User.GetTasks"

	const query = `
		SELECT task.id, task.title, task.is_done FROM task
			JOIN user_task ON task.id = user_task.task_id
			JOIN "user" ON "user".id = user_task.user_id
			WHERE "user".login = $1;
		`

	stmt, err := user.db.Prepare(query)
	if err != nil {
		return Tasks{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(user.login)
	if err != nil {
		return Tasks{}, fmt.Errorf("%s: failed to query tasks table: %w", funcErrMsg, err)
	}

	defer rows.Close()

	tasks := Tasks{}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Title, &task.IsDone)
		if err != nil {
			return Tasks{}, fmt.Errorf("%s: failed to scan rows: %w", funcErrMsg, err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return Tasks{}, fmt.Errorf("%s: rows error: %w", funcErrMsg, err)
	}

	return tasks, nil
}

func (user *User) NewTask(title string, isDone bool) (Task, error) {
	const funcErrMsg = "models.User.NewTask"

	stmt, err := user.db.Prepare("INSERT INTO task(title, is_done) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(title, isDone)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a statement: %w", funcErrMsg, err)
	}

	if !rows.Next() {
		return Task{}, fmt.Errorf(
			"%s: failed to get an id of a last inserted task: %w",
			funcErrMsg,
			err,
		)
	}

	if err = rows.Err(); err != nil {
		return Task{}, fmt.Errorf("%s: rows.Err(): %w", funcErrMsg, err)
	}

	var taskID int

	err = rows.Scan(&taskID)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to scan a task id : %w", funcErrMsg, err)
	}

	defer rows.Close()

	stmt, err = user.db.Prepare("INSERT INTO user_task(user_id, task_id) VALUES ($1, $2)")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.id, taskID)
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to execute a statement: %w", funcErrMsg, err)
	}

	user.log.Info(
		"Succsesfully inserted:",
		"id",
		taskID,
		"title",
		title,
		"isDone",
		isDone,
	)

	task := Task{taskID, title, isDone}

	return task, nil
}

func (user *User) RemoveTask(taskID int) error {
	const funcErrMsg = "models.User.RemoveTask"

	stmt, err := user.db.Prepare("DELETE FROM task WHERE id = $1")
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

func (user *User) GetTaskByID(taskID int) (Task, error) {
	const funcErrMsg = "models.User.GetTask"

	stmt, err := user.db.Prepare("SELECT * FROM task WHERE id = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to prepare a statement: %w", funcErrMsg, err)
	}

	task, err := user.GetTask(stmt, strconv.Itoa(taskID))
	if err != nil {
		return Task{}, fmt.Errorf("%s: failed to get a task: %w", funcErrMsg, err)
	}

	return task, nil
}

func (user *User) GetTask(stmt *sql.Stmt, title string) (Task, error) {
	const funcErrMsg = "models.User.GetTask"

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

func (user *User) SetDoneStatus(taskID int, isDone bool) error {
	const funcErrMsg = "models.User.ToggleDoneStatus"

	stmt, err := user.db.Prepare("UPDATE task SET is_done = $1 WHERE id = $2")
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
