// Package api defines the public contract for the token-issue cell.
package api

import (
	"context"
	"errors"
)

var ErrEmptySubject = errors.New("token-issue: empty subject")

// Issuer creates cryptographically random tokens and their SHA-256 digests.
type Issuer interface {
	Issue(context.Context, string) (Token, error)
}

// Token is the public representation of an issued credential.
type Token struct {
	Value  string `json:"value"`
	Digest string `json:"digest"`
}
