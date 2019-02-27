package resourceadapter

import (
	"context"

	. "github.com/blocksafe/go/protocols/aurora"
	"github.com/blocksafe/go/services/aurora/internal/db2/history"
)

func PopulateHistoryAccount(ctx context.Context, dest *HistoryAccount, row history.Account) {
	dest.ID = row.Address
	dest.AccountID = row.Address
}
