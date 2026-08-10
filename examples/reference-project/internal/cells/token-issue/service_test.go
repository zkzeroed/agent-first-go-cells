package tokenissue

import (
	"errors"
	"testing"

	"example.com/agent-first-reference/internal/cells/token-issue/api"
)

func TestIssue(t *testing.T) {
	issuer := New(Deps{Store: NewStore()})
	token, err := issuer.Issue(t.Context(), "ada")
	if err != nil || token.Value == "" || len(token.Digest) != 64 {
		t.Fatalf("Issue() = %#v, %v", token, err)
	}
	_, err = issuer.Issue(t.Context(), " ")
	if !errors.Is(err, api.ErrEmptySubject) {
		t.Fatalf("Issue(empty) error = %v", err)
	}
}
