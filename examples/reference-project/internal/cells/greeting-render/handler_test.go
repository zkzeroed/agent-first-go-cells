package greetingrender

import (
	"net/http"
	"net/http/httptest"
	"testing"

	greetingcompose "example.com/agent-first-reference/internal/cells/greeting-compose"
)

func TestHandler(t *testing.T) {
	composer := greetingcompose.New(greetingcompose.Deps{Store: greetingcompose.NewStore("Hello")})
	handler := NewHandler(Deps{Composer: composer})
	req := httptest.NewRequest(http.MethodGet, "/greeting/Ada", nil)
	req.SetPathValue("name", "Ada")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK || response.Body.String() != "Hello, Ada!" {
		t.Fatalf("response = %d %q", response.Code, response.Body.String())
	}
}
