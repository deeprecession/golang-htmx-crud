package models

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func NewSessionStore() SessionStore {
	return SessionStore{}
}

type SessionStore map[string]string

func (store *SessionStore) GetSession(request *http.Request, key string) (string, error) {
	cookie, err := request.Cookie(key)
	if err != nil {
		return "", fmt.Errorf("failed to get cookie")
	}

	value, ok := (*store)[string(cookie.Value)]
	if !ok {
		return "", fmt.Errorf("session is not found")
	}

	return value, nil
}

func (store *SessionStore) SetSession(response *http.ResponseWriter, key string, val string) error {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(val))

	cookie := &http.Cookie{
		Name:     key,
		Value:    encodedValue,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	http.SetCookie(*response, cookie)

	(*store)[encodedValue] = val

	return nil
}
