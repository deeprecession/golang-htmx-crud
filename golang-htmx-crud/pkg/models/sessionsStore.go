package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session is not found")
	ErrBadRequest      = errors.New("bad request")
)

func NewSessionStore() SessionStore {
	return SessionStore{}
}

type SessionStore map[string]string

func (store *SessionStore) GetSession(request *http.Request, key string) (string, error) {
	cookie, err := request.Cookie(key)
	if err != nil {
		return "", fmt.Errorf("failed to get cookie: %w", ErrBadRequest)
	}

	value, ok := (*store)[cookie.Value]
	if !ok {
		return "", ErrSessionNotFound
	}

	return value, nil
}

func (store *SessionStore) SetSession(response *http.ResponseWriter, key string, val string) error {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(val))

	const dayDuration = 24 * time.Hour

	cookie := &http.Cookie{
		Name:     key,
		Value:    encodedValue,
		Expires:  time.Now().Add(dayDuration),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	http.SetCookie(*response, cookie)

	(*store)[encodedValue] = val

	return nil
}
