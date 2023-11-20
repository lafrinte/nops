package runtimex

import (
	"fmt"
	"github.com/lafrinte/nops/psutil"
	"math"
	"runtime"
)

// SetGoMaxProcPercent
/*
 use percent to reset runtime.GOMAXPROCS. when in container system, SetGoMaxProcPercent will skip all action
 @argument:
   p: range 1 100

 @example:
   1:
     vcpu := 8
     p := 10     # means 10% of cpu resource
     GOMAXPROCS equal to 1
   2:
     vcpu := 10
     p := 20     # means 20% of cpu resource
     GOMAXPROCS equal to 2
*/
func SetGoMaxProcPercent(p int) error {
	if !psutil.IsContainer() {
		cores := psutil.VCPUCount()
		if cores > 0 && (p > 0 && p <= 100) {
			rs := float64(cores) * (float64(p) / 100)
			if rs < 1 {
				runtime.GOMAXPROCS(int(math.Ceil(rs)))
				return nil
			}

			runtime.GOMAXPROCS(int(math.Floor(rs)))
			return nil
		}

		log.Warn().Int("cores", cores).Int("percent", p).Msg("arguments may invalid: (cores, percent)")

		return fmt.Errorf("arguments may invalid: (cores: %d, percent %d)", cores, p)
	}

	log.Info().Msg("container env, use default GOMAXPROC in cgroup")
	return nil
}
