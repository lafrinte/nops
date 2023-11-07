package sock

import (
	logger "github.com/lafrinte/nops/log"
)

var log = logger.DefaultLogger().GetLogger().With().Str("module", "sock").Logger()
