package aurora

import (
	"net/http"

	"github.com/blocksafe/go/services/aurora/internal/db2/core"
	"github.com/blocksafe/go/services/aurora/internal/db2/history"
	"github.com/blocksafe/go/services/aurora/internal/txsub"
	results "github.com/blocksafe/go/services/aurora/internal/txsub/results/db"
	"github.com/blocksafe/go/services/aurora/internal/txsub/sequence"
)

func initSubmissionSystem(app *App) {
	cq := &core.Q{Session: app.CoreSession(nil)}

	app.submitter = &txsub.System{
		Pending:         txsub.NewDefaultSubmissionList(),
		Submitter:       txsub.NewDefaultSubmitter(http.DefaultClient, app.config.BlocksafeCoreURL),
		SubmissionQueue: sequence.NewManager(),
		Results: &results.DB{
			Core:    cq,
			History: &history.Q{Session: app.AuroraSession(nil)},
		},
		Sequences:         cq.SequenceProvider(),
		NetworkPassphrase: app.config.NetworkPassphrase,
	}
}

func init() {
	appInit.Add("txsub", initSubmissionSystem, "app-context", "log", "aurora-db", "core-db")
}
