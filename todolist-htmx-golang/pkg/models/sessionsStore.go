package models

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionNotFound = errors.New("session is not found")
	ErrBadRequest      = errors.New("bad request")
)

func NewSessionStore(rdb redis.UniversalClient) SessionStore {
	return SessionStore{
		rdb: rdb,
	}
}

type SessionStore struct {
	rdb redis.UniversalClient
}

func (store *SessionStore) GetSession(request *http.Request, key string) (string, error) {
	const errFuncMsg = "models.SessionStore.GetSession"

	cookie, err := request.Cookie(key)
	if err != nil {
		return "", fmt.Errorf("failed to get cookie: %w", ErrBadRequest)
	}

	ctx := context.Background()

	value, err := store.rdb.Get(ctx, cookie.Value).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrSessionNotFound
		}

		return "", fmt.Errorf("%s failed to get session: %w", errFuncMsg, err)
	}

	return value, nil
}

func (store *SessionStore) SetSession(response *http.ResponseWriter, key string, val string) error {
	const errFuncMsg = "models.SessionStore.SetSession"

	encodedValue := base64.StdEncoding.EncodeToString([]byte(val))

	const expireDuration = 24 * time.Hour

	cookie := &http.Cookie{
		Name:     key,
		Value:    encodedValue,
		Expires:  time.Now().Add(expireDuration),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	http.SetCookie(*response, cookie)

	ctx := context.Background()

	_, err := store.rdb.Set(ctx, encodedValue, val, expireDuration).Result()
	if err != nil {
		return fmt.Errorf("%s failed to set session: %w", errFuncMsg, err)
	}

	return nil
}
