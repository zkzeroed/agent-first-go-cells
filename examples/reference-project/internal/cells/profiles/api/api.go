// Package api defines shared public contracts for the profiles domain.
package api

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("profiles: not found")

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Repository interface {
	Create(context.Context, Profile) error
	Find(context.Context, string) (Profile, error)
}
