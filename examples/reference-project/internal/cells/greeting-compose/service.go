package greetingcompose

import (
	"context"
	"strings"
)

type service struct{ store Store }

func (s service) Compose(ctx context.Context, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrEmptyName
	}
	prefix, err := s.store.Prefix(ctx)
	if err != nil {
		return "", err
	}
	return prefix + ", " + name + "!", nil
}
