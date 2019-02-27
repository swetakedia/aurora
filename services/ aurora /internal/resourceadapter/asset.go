package resourceadapter

import (
	"context"

	. "github.com/blocksafe/go/protocols/aurora"
	"github.com/blocksafe/go/xdr"
)

func PopulateAsset(ctx context.Context, dest *Asset, asset xdr.Asset) error {
	return asset.Extract(&dest.Type, &dest.Code, &dest.Issuer)
}
