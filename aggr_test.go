package aggr

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type item struct {
	Date   time.Time
	Value  float64
	Merged bool
}

type aggTestCase struct {
	data        []item
	exp         []item
	granularity int
	expN        int
	loc         *time.Location
}

func TestByDate(t *testing.T) {
	testCases := []aggTestCase{
		{
			loc:         time.UTC,
			granularity: Day,
			expN:        3,
			data: []item{
				{
					Value: 1,
					Date:  time.Date(2020, time.September, 4, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 3,
					Date:  time.Date(2020, time.September, 6, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 2,
					Date:  time.Date(2020, time.September, 5, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 1,
					Date:  time.Date(2020, time.September, 4, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 2,
					Date:  time.Date(2020, time.September, 5, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 3,
					Date:  time.Date(2020, time.September, 6, 11, 12, 13, 14, time.UTC),
				},
			},
			exp: []item{
				{
					Value:  2,
					Date:   time.Date(2020, time.September, 4, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
				{
					Value:  4,
					Date:   time.Date(2020, time.September, 5, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
				{
					Value:  6,
					Date:   time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
			},
		},
		{
			loc:         time.UTC,
			granularity: Month,
			expN:        2,
			data: []item{
				{
					Value: 3,
					Date:  time.Date(2019, time.September, 4, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 3,
					Date:  time.Date(2019, time.September, 6, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 3,
					Date:  time.Date(2019, time.September, 5, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 6,
					Date:  time.Date(2019, time.October, 5, 11, 12, 13, 14, time.UTC),
				},
			},
			exp: []item{
				{
					Value:  9,
					Date:   time.Date(2019, time.September, 1, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
				{
					Value:  6,
					Date:   time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
			},
		},
		{
			loc:         time.UTC,
			granularity: Week,
			expN:        2,
			data: []item{
				{
					Value: 3,
					Date:  time.Date(2019, time.September, 10, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 4,
					Date:  time.Date(2019, time.September, 9, 11, 12, 13, 14, time.UTC),
				},
				{
					Value: 5,
					Date:  time.Date(2019, time.September, 8, 11, 12, 13, 14, time.UTC),
				},
			},
			exp: []item{
				{
					Value:  5,
					Date:   time.Date(2019, time.September, 2, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
				{
					Value:  7,
					Date:   time.Date(2019, time.September, 9, 0, 0, 0, 0, time.UTC),
					Merged: true,
				},
			},
		},
		{
			loc:         time.UTC,
			granularity: Week,
			expN:        0,
			data:        []item{},
			exp:         []item{},
		},
	}

	for i, tc := range testCases {
		data := tc.data
		ln := ByDate(
			data,
			tc.granularity,
			func(i int) time.Time {
				// getting date as time.Time from item
				return data[i].Date
			},
			func(toMerge []int, i int, d time.Time) {
				merged := item{
					Merged: true,
					Date:   d,
				}

				// combine all values of items being merged
				for _, j := range toMerge {
					merged.Value += data[j].Value
				}

				// writing merged item to ith position
				data[i] = merged
			},
			tc.loc,
		)
		data = data[:ln]

		assert.Equal(t, tc.expN, ln, fmt.Sprintf("test case %d", i+1))
		assert.Equal(t, tc.exp, data, fmt.Sprintf("test case %d", i+1))
	}
}

func TestWrongGranularity(t *testing.T) {
	data := []struct{}{
		{},
		{},
	}
	assert.Panics(t, func() {
		ByDate(data, 6666, func(int) time.Time { return time.Now() }, func([]int, int, time.Time) {}, time.UTC)
	})
}

func TestWrongInputData(t *testing.T) {
	assert.Panics(t, func() {
		ByDate(1, Day, func(int) time.Time { return time.Now() }, func([]int, int, time.Time) {}, time.UTC)
	})
}

func TestSortedKeys(t *testing.T) {
	// as map order is random we repeat test 100 times
	// to make sure random order doesn't affect anything
	for at := 0; at < 100; at++ {
		type testCase struct {
			keys [][2]int
			exp  [][2]int
		}

		testCases := []testCase{
			{
				keys: [][2]int{
					{2014, 12},
					{2014, 7},
					{2016, 10},
					{2014, 5},
				},
				exp: [][2]int{
					{2014, 5},
					{2014, 7},
					{2014, 12},
					{2016, 10},
				},
			},
			{
				keys: [][2]int{
					{2019, 365},
					{2019, 11},
					{2019, 256},
					{2019, 124},
				},
				exp: [][2]int{
					{2019, 11},
					{2019, 124},
					{2019, 256},
					{2019, 365},
				},
			},
			{
				keys: [][2]int{
					{2011, 11},
					{2022, 11},
					{2019, 11},
					{2010, 11},
				},
				exp: [][2]int{
					{2010, 11},
					{2011, 11},
					{2019, 11},
					{2022, 11},
				},
			},
		}

		for tcN, tc := range testCases {
			m := make(map[[2]int][]int)
			for i, k := range tc.keys {
				m[k] = []int{i}
			}
			res := sortedKeys(m)
			assert.Equal(t, tc.exp, res, fmt.Sprintf("at %d, test case %d", at, tcN))
		}
	}
}
