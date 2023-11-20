package runtimex

import (
	"fmt"
	"github.com/lafrinte/nops/psutil"
	"github.com/lafrinte/nops/str"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

// SetGoMemoryLimitPercent
/*
 use percent to reset debug.SetMemoryLimit. when in container system, SetGoMemoryLimitPercent will skip all action
 @argument:
   p: range 1 100
*/
func SetGoMemoryLimitPercent(p int) error {
	if !psutil.IsContainer() {
		total := psutil.GetTotalMem()
		if total > 0 && (p > 0 && p <= 100) {
			debug.SetMemoryLimit(int64(math.Floor(float64(total) * (float64(p) / 100))))
			return nil
		}

		return fmt.Errorf("arguments may invalid: (totalMem: %d, percent %d)", total, p)
	}

	log.Info().Msg("container env, use default MemoryLimit in cgroup")

	return nil
}

// SetGoMemoryLimit
/*
 use percent to reset debug.SetMemoryLimit. when in container system, SetGoMemoryLimit will skip all action
 @argument:
   memString: <int><unit>
     :int: valid int
     :unit: B,KB,MB,GB,TB,PB. case insensitive. unit configuration empty means 'B'
*/
func SetGoMemoryLimit(memString string) error {
	if !psutil.IsContainer() {
		v, err := parsingMemSetting(memString)
		if err != nil {
			return fmt.Errorf("invalid memory configuration: %s", err)
		}

		debug.SetMemoryLimit(int64(v))
		return nil
	}

	log.Info().Msg("container env, use default MemoryLimit in cgroup")

	return nil
}

func parsingMemSetting(memString string) (uint64, error) {
	firstCharIndex := 0

	memString = strings.ToUpper(strings.Trim(memString, " "))

	if str.IsNumeric(memString) {
		v, _ := strconv.Atoi(memString)
		return uint64(v), nil
	}

	for index, s := range memString {
		// to string
		v := string(s)

		// is numeric
		_, err := strconv.Atoi(v)
		if err != nil {
			if firstCharIndex == 0 {
				firstCharIndex = index
				break
			}
		}
	}

	unit := memString[firstCharIndex:]
	value, _ := strconv.Atoi(memString[:firstCharIndex])

	switch unit {
	case "B":
		return uint64(value), nil
	case "KB":
		return uint64(value) * KB, nil
	case "MB":
		return uint64(value) * MB, nil
	case "GB":
		return uint64(value) * GB, nil
	case "TB":
		return uint64(value) * TB, nil
	case "PB":
		return uint64(value) * PB, nil
	default:
		return 0, fmt.Errorf("format error. configuration proxy: [int][unit] and unit only enables B,KB,MB,GB,TB,PB")
	}
}
