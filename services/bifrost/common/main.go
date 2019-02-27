package common

import (
	"github.com/blocksafe/go/support/log"
)

const BlocksafeAmountPrecision = 7

func CreateLogger(serviceName string) *log.Entry {
	return log.DefaultLogger.WithField("service", serviceName)
}
