package greetingcompose

import "context"

// Store supplies the greeting prefix without exposing its implementation.
type Store interface {
	Prefix(context.Context) (string, error)
}

type store struct{ prefix string }

func NewStore(prefix string) Store { return store{prefix: prefix} }

func (s store) Prefix(context.Context) (string, error) { return s.prefix, nil }
