// Package api defines the public contract for the profile-create action.
package api

import (
	"context"

	profilesapi "example.com/agent-first-reference/internal/cells/profiles/api"
)

type Creator interface {
	Create(context.Context, string, string) (profilesapi.Profile, error)
}
