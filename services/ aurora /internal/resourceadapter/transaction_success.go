package resourceadapter

import (
	"context"

	. "github.com/blocksafe/go/protocols/aurora"
	"github.com/blocksafe/go/services/aurora/internal/httpx"
	"github.com/blocksafe/go/services/aurora/internal/txsub"
	"github.com/blocksafe/go/support/render/hal"
)

// Populate fills out the details
func PopulateTransactionSuccess(ctx context.Context, dest *TransactionSuccess, result txsub.Result) {
	dest.Hash = result.Hash
	dest.Ledger = result.LedgerSequence
	dest.Env = result.EnvelopeXDR
	dest.Result = result.ResultXDR
	dest.Meta = result.ResultMetaXDR

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	dest.Links.Transaction = lb.Link("/transactions", result.Hash)
	return
}
