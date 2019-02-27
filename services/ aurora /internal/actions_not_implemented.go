package aurora

import (
	"github.com/blocksafe/go/services/aurora/internal/actions"
	hProblem "github.com/blocksafe/go/services/aurora/internal/render/problem"
	"github.com/blocksafe/go/support/render/problem"
)

// Interface verification
var _ actions.JSONer = (*NotImplementedAction)(nil)

// NotImplementedAction renders a NotImplemented prblem
type NotImplementedAction struct {
	Action
}

// JSON is a method for actions.JSON
func (action *NotImplementedAction) JSON() error {
	problem.Render(action.R.Context(), action.W, hProblem.NotImplemented)
	return action.Err
}
