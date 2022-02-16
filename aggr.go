package aggr

import (
	"reflect"
	"sort"
	"time"
)

const (
	Day = iota
	Week
	Month
	wrongGranularityPanicMsg = "granularity value is not known"
)

// ByDate aggregates slice by date of objects with required granularity.
// User needs to provide getDateFn which tells Slice how to get date for the object with index i.
// Also user needs to provide mergeFn which merges objects with given indexes
// and writes the result to given index.
// Sorts slice by date from earlier to later.
// Panics if data is not slice.
// Panics if granularity is not valid.
func ByDate(
	data interface{},
	granularity int,
	getDateFn func(i int) time.Time,
	mergeFn func(toMerge []int, i int, date time.Time),
	loc *time.Location,
) int {
	rv := reflect.ValueOf(data)
	rl := rv.Len()

	// make sure slice is sorted
	sort.Slice(data, func(i, j int) bool {
		return getDateFn(i).Before(getDateFn(j))
	})

	// function which presents date unit (day, week, etc.)
	// as a pair of ints unique for that date unit
	getIdx := getTimeToIdxFn(granularity)

	// This is map of date unit (day, week, etc.) index
	// to list of indexes of objects in input data.
	// TODO: maybe switch to using ordered map e.g. https://github.com/wk8/go-ordered-map
	// it may be more efficient though more testing required to make sure.
	valueMap := make(map[[2]int][]int)
	for i := 0; i < rl; i++ {
		di := getIdx(getDateFn(i))
		valueMap[di] = append(valueMap[di], i)
	}

	// get function which turns date unit index to time.Date
	// in provided timezone
	getPeriodStartFn := getIdxToTimeFn(granularity)

	keys := sortedKeys(valueMap)
	for i, k := range keys {
		// this should iterate over map in sorted order
		// and put merged items to ith place in input slice.
		mergeFn(valueMap[k], i, getPeriodStartFn(k, loc))
	}

	return len(keys)
}

func getTimeToIdxFn(granularity int) func(time.Time) [2]int {
	switch granularity {
	case Day:
		return dayIdx
	case Week:
		return weekIdx
	case Month:
		return monthIdx
	default:
		panic(wrongGranularityPanicMsg)
	}
}

func getIdxToTimeFn(granularity int) func([2]int, *time.Location) time.Time {
	switch granularity {
	case Day:
		return idxToDayStart
	case Week:
		return idxToWeekStart
	case Month:
		return idxToMonthStart
	default:
		panic(wrongGranularityPanicMsg)
	}
}

func sortedKeys(m map[[2]int][]int) [][2]int {
	keys := make([][2]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i int, j int) bool {
		yearI, otherI, yearJ, otherJ := keys[i][0], keys[i][1], keys[j][0], keys[j][1]
		if yearI < yearJ {
			return true
		} else if yearI > yearJ {
			return false
		}
		if otherI < otherJ {
			return true
		}
		return false
	})

	return keys
}
