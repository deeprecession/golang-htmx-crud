package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
	ErrBadPassword      = errors.New("bad password")
)

type User struct {
	login    string
	password string
}

type UserStorage struct {
	log      *slog.Logger
	database *sql.DB
}

func GetUserStorage(log *slog.Logger, database *sql.DB) UserStorage {
	return UserStorage{log, database}
}

func (storage *UserStorage) Register(login, password string) error {
	const funcErrMsg = "storage.UserStorage.Register"

	isLoginTaken, err := storage.loginIsTaken(login)
	if err != nil {
		return fmt.Errorf(
			"%s failed to check is login is taken: %w",
			funcErrMsg,
			ErrUserAlreadyExist,
		)
	}

	if isLoginTaken {
		return fmt.Errorf("%w", ErrUserAlreadyExist)
	}

	user := User{login, password}

	err = storage.addUser(user)
	if err != nil {
		return fmt.Errorf("%s failed to add a user: %w", funcErrMsg, err)
	}

	storage.log.Info("inseted a new user", "user", user)

	return nil
}

func (storage *UserStorage) Login(login, password string) error {
	const funcErrMsg = "storage.UserStorage.Login"

	storedUser, err := storage.GetUserWithLogin(login)
	if err != nil {
		return fmt.Errorf("%s failed to login: %w", funcErrMsg, err)
	}

	if storedUser.password != password {
		return fmt.Errorf("%s failed to login: %w", funcErrMsg, ErrBadPassword)
	}

	return nil
}

func (storage *UserStorage) GetUserWithLogin(login string) (User, error) {
	const funcErrMsg = "storage.UserStorage.GetUserWithLogin"

	const query = `SELECT login, password FROM users WHERE login = $1`

	stmt, err := storage.database.Prepare(query)
	if err != nil {
		return User{}, fmt.Errorf("%s failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(login)
	if err != nil {
		return User{}, fmt.Errorf("%s failed to query a statement: %w", funcErrMsg, err)
	}

	defer rows.Close()

	isUserNotExist := !rows.Next()

	if err = rows.Err(); err != nil {
		return User{}, fmt.Errorf("%s failed to check rows.Err(): %w", funcErrMsg, err)
	}

	if isUserNotExist {
		return User{}, fmt.Errorf("%w", ErrUserNotFound)
	}

	user := User{}

	err = rows.Scan(&user.login, &user.password)
	if err != nil {
		return User{}, fmt.Errorf("%s failed to scan a user: %w", funcErrMsg, err)
	}

	return user, nil
}

func (storage *UserStorage) addUser(user User) error {
	const funcErrMsg = "storage.UserStorage.addUser"

	const query = `INSERT INTO users(login, password) VALUES ($1, $2);`

	stmt, err := storage.database.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.login, user.password)
	if err != nil {
		return fmt.Errorf("%s failed to execute a statement: %w", funcErrMsg, err)
	}

	return nil
}

func (storage UserStorage) loginIsTaken(login string) (bool, error) {
	const funcErrMsg = "storage.UserStorage.loginIsTaken"

	const query = `SELECT * FROM users WHERE login = $1`

	stmt, err := storage.database.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s failed to prepare a statement: %w", funcErrMsg, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(login)
	if err != nil {
		return false, fmt.Errorf("%s failed to query a statement: %w", funcErrMsg, err)
	}

	defer rows.Close()

	isUserAlreadyExist := rows.Next()

	if err = rows.Err(); err != nil {
		return false, fmt.Errorf("%s failed to check rows.Err(): %w", funcErrMsg, err)
	}

	return isUserAlreadyExist, nil
}
