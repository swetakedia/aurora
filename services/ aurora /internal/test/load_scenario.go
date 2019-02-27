package test

import (
	"github.com/blocksafe/go/services/aurora/internal/test/scenarios"
)

func loadScenario(scenarioName string, includeAurora bool) {
	blocksafeCorePath := scenarioName + "-core.sql"
	auroraPath := scenarioName + "-aurora.sql"

	if !includeAurora {
		auroraPath = "blank-aurora.sql"
	}

	scenarios.Load(BlocksafeCoreDatabaseURL(), blocksafeCorePath)
	scenarios.Load(DatabaseURL(), auroraPath)
}
