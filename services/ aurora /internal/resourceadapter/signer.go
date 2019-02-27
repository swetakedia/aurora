package resourceadapter

import (
	"context"

	. "github.com/blocksafe/go/protocols/aurora"
	"github.com/blocksafe/go/services/aurora/internal/db2/core"
)

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func PopulateSigner(ctx context.Context, dest *Signer, row core.Signer) {
	dest.Weight = row.Weight
	dest.Key = row.Publickey
	dest.Type = MustKeyTypeFromAddress(dest.Key)
}

// PopulateMaster fills out the fields of the signer, using a blocksafe account to
// provide the data.
func PopulateMasterSigner(dest *Signer, row core.Account) {
	dest.Weight = int32(row.Thresholds[0])
	dest.Key = row.Accountid
	dest.Type = MustKeyTypeFromAddress(dest.Key)
}
