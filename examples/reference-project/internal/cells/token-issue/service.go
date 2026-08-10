package tokenissue

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"example.com/agent-first-reference/internal/cells/token-issue/api"
)

type service struct{ store Store }

func (s service) Issue(ctx context.Context, subject string) (api.Token, error) {
	if strings.TrimSpace(subject) == "" {
		return api.Token{}, api.ErrEmptySubject
	}
	bytes := make([]byte, tokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return api.Token{}, err
	}
	value := base64.RawURLEncoding.EncodeToString(bytes)
	digest := sha256.Sum256(bytes)
	encodedDigest := hex.EncodeToString(digest[:])
	if err := s.store.Save(ctx, subject, encodedDigest); err != nil {
		return api.Token{}, err
	}
	return api.Token{Value: value, Digest: encodedDigest}, nil
}
