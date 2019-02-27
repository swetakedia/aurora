package ingest

import (
	"database/sql"
	"fmt"

	"github.com/blocksafe/go/services/aurora/internal/db2/core"
	"github.com/blocksafe/go/support/db"
	"github.com/blocksafe/go/support/errors"
)

// Load runs queries against `core` to fill in the records of the bundle.
func (lb *LedgerBundle) Load(db *db.Session) error {
	q := &core.Q{Session: db}
	// Load Header
	err := q.LedgerHeaderBySequence(&lb.Header, lb.Sequence)
	if err != nil {
		// Remove when Aurora is able to handle gaps in blocksafe-core DB.
		// More info:
		// * https://github.com/blocksafe/go/issues/335
		// * https://www.blocksafe.org/developers/software/known-issues.html#gaps-detected
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("Gap detected in blocksafe-core database (ledger=%d). More information: https://www.blocksafe.org/developers/software/known-issues.html#gaps-detected", lb.Sequence))
		}
		return errors.Wrap(err, "failed to load header")
	}

	// Load transactions
	err = q.TransactionsByLedger(&lb.Transactions, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to load transactions")
	}

	err = q.TransactionFeesByLedger(&lb.TransactionFees, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to load transaction fees")
	}

	return nil
}
