package greetingrender

import (
	"errors"
	"net/http"

	composeapi "example.com/agent-first-reference/internal/cells/greeting-compose/api"
)

type handler struct{ composer composeapi.Composer }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	greeting, err := h.composer.Compose(r.Context(), r.PathValue("name"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, composeapi.ErrEmptyName) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}
	_, _ = w.Write([]byte(greeting))
}
