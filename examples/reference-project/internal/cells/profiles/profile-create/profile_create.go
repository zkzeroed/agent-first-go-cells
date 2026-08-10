package profilecreate

import (
	"context"
	"errors"
	"strings"

	profilesapi "example.com/agent-first-reference/internal/cells/profiles/api"
	"example.com/agent-first-reference/internal/cells/profiles/profile-create/api"
)

var ErrInvalidProfile = errors.New("profile-create: id and name are required")

type Deps struct{ Repository profilesapi.Repository }

func New(deps Deps) api.Creator { return service{repository: deps.Repository} }

type service struct{ repository profilesapi.Repository }

func (s service) Create(ctx context.Context, id, name string) (profilesapi.Profile, error) {
	profile := profilesapi.Profile{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name)}
	if profile.ID == "" || profile.Name == "" {
		return profilesapi.Profile{}, ErrInvalidProfile
	}
	if err := s.repository.Create(ctx, profile); err != nil {
		return profilesapi.Profile{}, err
	}
	return profile, nil
}
