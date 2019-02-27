package aurora

import (
	"net/http"

	"github.com/blocksafe/go/protocols/aurora"
	"github.com/blocksafe/go/services/aurora/internal/actions"
	"github.com/blocksafe/go/services/aurora/internal/db2"
	"github.com/blocksafe/go/services/aurora/internal/db2/history"
	hProblem "github.com/blocksafe/go/services/aurora/internal/render/problem"
	"github.com/blocksafe/go/services/aurora/internal/render/sse"
	"github.com/blocksafe/go/services/aurora/internal/resourceadapter"
	"github.com/blocksafe/go/services/aurora/internal/txsub"
	"github.com/blocksafe/go/support/errors"
	"github.com/blocksafe/go/support/render/hal"
	"github.com/blocksafe/go/support/render/problem"
)

// This file contains the actions:
//
// TransactionIndexAction: pages of transactions
// TransactionShowAction: single transaction by sequence, by hash or id

// Interface verifications
var _ actions.JSONer = (*TransactionIndexAction)(nil)
var _ actions.EventStreamer = (*TransactionIndexAction)(nil)

// TransactionIndexAction renders a page of ledger resources, identified by
// a normal page query.
type TransactionIndexAction struct {
	Action
	LedgerFilter  int32
	AccountFilter string
	PagingParams  db2.PageQuery
	Records       []history.Transaction
	Page          hal.Page
	IncludeFailed bool
}

// JSON is a method for actions.JSON
func (action *TransactionIndexAction) JSON() error {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() { hal.Render(action.W, action.Page) },
	)
	return action.Err
}

// SSE is a method for actions.SSE
func (action *TransactionIndexAction) SSE(stream *sse.Stream) error {
	action.Setup(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.ValidateCursorWithinHistory,
	)
	action.Do(
		action.loadRecords,
		func() {
			stream.SetLimit(int(action.PagingParams.Limit))
			records := action.Records[stream.SentCount():]

			for _, record := range records {
				var res aurora.Transaction
				resourceadapter.PopulateTransaction(action.R.Context(), &res, record)
				stream.Send(sse.Event{ID: res.PagingToken(), Data: res})
			}
		},
	)

	return action.Err
}

func (action *TransactionIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.AccountFilter = action.GetAddress("account_id")
	action.LedgerFilter = action.GetInt32("ledger_id")
	action.PagingParams = action.GetPageQuery()
	action.IncludeFailed = action.GetBool("include_failed")

	if action.IncludeFailed == true && !action.App.config.IngestFailedTransactions {
		err := errors.New("`include_failed` parameter is unavailable when Aurora is not ingesting failed " +
			"transactions. Set `INGEST_FAILED_TRANSACTIONS=true` to start ingesting them.")
		action.Err = problem.MakeInvalidFieldProblem("include_failed", err)
		return
	}
}

func (action *TransactionIndexAction) loadRecords() {
	q := action.HistoryQ()
	txs := q.Transactions()

	switch {
	case action.AccountFilter != "":
		txs.ForAccount(action.AccountFilter)
	case action.LedgerFilter > 0:
		txs.ForLedger(action.LedgerFilter)
	}

	if !action.IncludeFailed {
		txs.SuccessfulOnly()
	}

	action.Err = txs.Page(action.PagingParams).Select(&action.Records)
}

func (action *TransactionIndexAction) loadPage() {
	for _, record := range action.Records {
		var res aurora.Transaction
		resourceadapter.PopulateTransaction(action.R.Context(), &res, record)
		action.Page.Add(res)
	}

	action.Page.FullURL = action.FullURL()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

// Interface verification
var _ actions.JSONer = (*TransactionShowAction)(nil)

// TransactionShowAction renders a ledger found by its sequence number.
type TransactionShowAction struct {
	Action
	Hash     string
	Record   history.Transaction
	Resource aurora.Transaction
}

func (action *TransactionShowAction) loadParams() {
	action.Hash = action.GetString("tx_id")
}

func (action *TransactionShowAction) loadRecord() {
	action.Err = action.HistoryQ().TransactionByHash(&action.Record, action.Hash)
}

func (action *TransactionShowAction) loadResource() {
	resourceadapter.PopulateTransaction(action.R.Context(), &action.Resource, action.Record)
}

// JSON is a method for actions.JSON
func (action *TransactionShowAction) JSON() error {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() { hal.Render(action.W, action.Resource) },
	)
	return action.Err
}

// Interface verification
var _ actions.JSONer = (*TransactionCreateAction)(nil)

// TransactionCreateAction submits a transaction to the blocksafe-core network
// on behalf of the requesting client.
type TransactionCreateAction struct {
	Action
	TX       string
	Result   txsub.Result
	Resource aurora.TransactionSuccess
}

// JSON format action handler
func (action *TransactionCreateAction) JSON() error {
	action.Do(
		action.loadTX,
		action.loadResult,
		action.loadResource,
		func() { hal.Render(action.W, action.Resource) },
	)
	return action.Err
}

func (action *TransactionCreateAction) loadTX() {
	action.ValidateBodyType()
	action.TX = action.GetString("tx")
}

func (action *TransactionCreateAction) loadResult() {
	submission := action.App.submitter.Submit(action.R.Context(), action.TX)

	select {
	case result := <-submission:
		action.Result = result
	case <-action.R.Context().Done():
		action.Err = &hProblem.Timeout
	}
}

func (action *TransactionCreateAction) loadResource() {
	if action.Result.Err == nil {
		resourceadapter.PopulateTransactionSuccess(action.R.Context(), &action.Resource, action.Result)
		return
	}

	if action.Result.Err == txsub.ErrTimeout {
		action.Err = &hProblem.Timeout
		return
	}

	if action.Result.Err == txsub.ErrCanceled {
		action.Err = &hProblem.Timeout
		return
	}

	switch err := action.Result.Err.(type) {
	case *txsub.FailedTransactionError:
		rcr := aurora.TransactionResultCodes{}
		resourceadapter.PopulateTransactionResultCodes(action.R.Context(), &rcr, err)

		action.Err = &problem.P{
			Type:   "transaction_failed",
			Title:  "Transaction Failed",
			Status: http.StatusBadRequest,
			Detail: "The transaction failed when submitted to the blocksafe network. " +
				"The `extras.result_codes` field on this response contains further " +
				"details.  Descriptions of each code can be found at: " +
				"https://www.blocksafe.org/developers/learn/concepts/list-of-operations.html",
			Extras: map[string]interface{}{
				"envelope_xdr": action.Result.EnvelopeXDR,
				"result_xdr":   err.ResultXDR,
				"result_codes": rcr,
			},
		}
	case *txsub.MalformedTransactionError:
		action.Err = &problem.P{
			Type:   "transaction_malformed",
			Title:  "Transaction Malformed",
			Status: http.StatusBadRequest,
			Detail: "Aurora could not decode the transaction envelope in this " +
				"request. A transaction should be an XDR TransactionEnvelope struct " +
				"encoded using base64.  The envelope read from this request is " +
				"echoed in the `extras.envelope_xdr` field of this response for your " +
				"convenience.",
			Extras: map[string]interface{}{
				"envelope_xdr": err.EnvelopeXDR,
			},
		}
	default:
		action.Err = err
	}
}
