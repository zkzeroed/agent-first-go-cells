package tokenrender

import (
	"encoding/json"
	"errors"
	"net/http"

	issueapi "example.com/agent-first-reference/internal/cells/token-issue/api"
)

type handler struct{ issuer issueapi.Issuer }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.issuer.Issue(r.Context(), r.PathValue("subject"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, issueapi.ErrEmptySubject) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}
