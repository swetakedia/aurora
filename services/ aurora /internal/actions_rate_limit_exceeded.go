package aurora

import (
	"net/http"

	hProblem "github.com/blocksafe/go/services/aurora/internal/render/problem"
	"github.com/blocksafe/go/support/render/problem"
)

// RateLimitExceededAction renders a 429 response
type RateLimitExceededAction struct {
	Action
	App *App
}

func (action RateLimitExceededAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(w, r)
	ap.App = action.App
	problem.Render(action.R.Context(), action.W, hProblem.RateLimitExceeded)
}
