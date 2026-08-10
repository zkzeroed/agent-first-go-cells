package tokenissue

import (
	"context"
	"sync"

	"example.com/agent-first-reference/internal/cells/token-issue/api"
)

// Store records digests without retaining raw tokens.
type Store interface {
	Save(context.Context, string, string) error
}

type store struct {
	mu      sync.Mutex
	records map[string]string
}

func NewStore() Store { return &store{records: make(map[string]string)} }

func (s *store) Save(_ context.Context, subject, digest string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[subject] = digest
	return nil
}

var _ api.Issuer = service{}
