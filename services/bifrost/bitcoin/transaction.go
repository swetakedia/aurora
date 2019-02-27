package bitcoin

import (
	"math/big"

	"github.com/blocksafe/go/services/bifrost/common"
)

func (t Transaction) ValueToBlocksafe() string {
	valueSat := new(big.Int).SetInt64(t.ValueSat)
	valueBtc := new(big.Rat).Quo(new(big.Rat).SetInt(valueSat), satInBtc)
	return valueBtc.FloatString(common.BlocksafeAmountPrecision)
}
