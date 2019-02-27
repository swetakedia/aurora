package ethereum

import (
	"math/big"

	"github.com/blocksafe/go/services/bifrost/common"
)

func (t Transaction) ValueToBlocksafe() string {
	valueEth := new(big.Rat)
	valueEth.Quo(new(big.Rat).SetInt(t.ValueWei), weiInEth)
	return valueEth.FloatString(common.BlocksafeAmountPrecision)
}
