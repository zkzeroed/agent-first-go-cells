package profiles

import (
	"context"
	"sync"

	"example.com/agent-first-reference/internal/cells/profiles/api"
)

type store struct {
	mu       sync.RWMutex
	profiles map[string]api.Profile
}

func NewStore() api.Repository { return &store{profiles: make(map[string]api.Profile)} }

func (s *store) Create(_ context.Context, profile api.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[profile.ID] = profile
	return nil
}

func (s *store) Find(_ context.Context, id string) (api.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	profile, ok := s.profiles[id]
	if !ok {
		return api.Profile{}, api.ErrNotFound
	}
	return profile, nil
}
