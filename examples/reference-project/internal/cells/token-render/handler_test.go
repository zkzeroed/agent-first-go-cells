package tokenrender

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	issueapi "example.com/agent-first-reference/internal/cells/token-issue/api"
)

func TestHandler(t *testing.T) {
	handler := NewHandler(Deps{Issuer: fakeIssuer{}})
	req := httptest.NewRequest(http.MethodGet, "/token/ada", nil)
	req.SetPathValue("subject", "ada")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK || response.Body.Len() == 0 {
		t.Fatalf("response = %d %q", response.Code, response.Body.String())
	}
}

type fakeIssuer struct{}

func (fakeIssuer) Issue(_ context.Context, _ string) (issueapi.Token, error) {
	return issueapi.Token{Value: "token", Digest: "digest"}, nil
}
