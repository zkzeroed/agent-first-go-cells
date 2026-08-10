package profilecreate

import (
	"context"
	"errors"
	"testing"

	profilesapi "example.com/agent-first-reference/internal/cells/profiles/api"
)

func TestCreate(t *testing.T) {
	creator := New(Deps{Repository: fakeRepository{}})
	profile, err := creator.Create(t.Context(), "ada", "Ada Lovelace")
	if err != nil || profile.Name != "Ada Lovelace" {
		t.Fatalf("Create() = %#v, %v", profile, err)
	}
	_, err = creator.Create(t.Context(), "", "Ada")
	if !errors.Is(err, ErrInvalidProfile) {
		t.Fatalf("Create(invalid) error = %v", err)
	}
}

type fakeRepository struct{}

func (fakeRepository) Create(context.Context, profilesapi.Profile) error { return nil }

func (fakeRepository) Find(context.Context, string) (profilesapi.Profile, error) {
	return profilesapi.Profile{}, profilesapi.ErrNotFound
}
