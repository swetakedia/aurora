package compliance

import (
	"github.com/asaskevich/govalidator"
	"github.com/blocksafe/go/address"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.CustomTypeTagMap.Set("blocksafe_address", govalidator.CustomTypeValidator(isBlocksafeAddress))
}

func isBlocksafeAddress(i interface{}, context interface{}) bool {
	addr, ok := i.(string)

	if !ok {
		return false
	}

	_, _, err := address.Split(addr)

	if err == nil {
		return true
	}

	return false
}
