package session

import (
	"context"
	"errors"
)

var ErrSessionNotFound = errors.New("session not found")

const (
	SessionKey = "_session"
)

func GetFromCtx(ctx context.Context) (*Session, error) {
	ses, ok := ctx.Value(SessionKey).(*Session)
	if !ok {
		return nil, ErrSessionNotFound
	}
	return ses, nil
}
