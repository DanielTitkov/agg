package aggr

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type idxTestCase struct {
	date   string
	layout string
	exp    [2]int
}

type startTestCase struct {
	idx [2]int
	exp time.Time
	loc *time.Location
}

func TestIdxToDayStart(t *testing.T) {
	locCEST, _ := time.LoadLocation("Europe/Berlin")
	for i, tc := range []startTestCase{
		{
			idx: [2]int{2019, 365},
			exp: time.Date(2019, time.December, 31, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			// for leap years day indexes should be already ajusted for that
			idx: [2]int{2020, 366},
			exp: time.Date(2020, time.December, 31, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			idx: [2]int{2021, 256},
			exp: time.Date(2021, time.September, 13, 0, 0, 0, 0, locCEST),
			loc: locCEST,
		},
	} {
		assert.Equal(t, tc.exp, idxToDayStart(tc.idx, tc.loc), fmt.Sprintf("test case %d", i+1))
	}
}

func TestIdxToMonthStart(t *testing.T) {
	locCEST, _ := time.LoadLocation("Europe/Berlin")
	for i, tc := range []startTestCase{
		{
			idx: [2]int{2019, 12},
			exp: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			idx: [2]int{2021, 9},
			exp: time.Date(2021, time.September, 1, 0, 0, 0, 0, locCEST),
			loc: locCEST,
		},
	} {
		assert.Equal(t, tc.exp, idxToMonthStart(tc.idx, tc.loc), fmt.Sprintf("test case %d", i+1))
	}
}

func TestIdxToWeekStart(t *testing.T) {
	locCEST, _ := time.LoadLocation("Europe/Berlin")
	for i, tc := range []startTestCase{
		{
			idx: [2]int{2022, 7},
			exp: time.Date(2022, time.February, 14, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			idx: [2]int{2020, 31},
			exp: time.Date(2020, time.July, 27, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			idx: [2]int{2018, 31},
			exp: time.Date(2018, time.July, 30, 0, 0, 0, 0, time.UTC),
			loc: time.UTC,
		},
		{
			idx: [2]int{2020, 49},
			exp: time.Date(2020, time.November, 30, 0, 0, 0, 0, locCEST),
			loc: locCEST,
		},
	} {
		assert.Equal(t, tc.exp, idxToWeekStart(tc.idx, tc.loc), fmt.Sprintf("test case %d", i+1))
	}
}

func TestDayIdx(t *testing.T) {
	for _, tc := range []idxTestCase{
		{
			date:   "2014-11-12T11:45:26.371Z",
			layout: "2006-01-02T15:04:05.000Z",
			exp:    [2]int{2014, 316},
		},
		{
			date:   "2020-12-31T11:45:26.371Z",
			layout: "2006-01-02T15:04:05.000Z",
			exp:    [2]int{2020, 366},
		},
	} {
		date, err := time.Parse(tc.layout, tc.date)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, tc.exp, dayIdx(date))
	}
}

func TestMonthIdx(t *testing.T) {
	for _, tc := range []idxTestCase{
		{
			date:   "2014-11-12T11:45:26.371Z",
			layout: "2006-01-02T15:04:05.000Z",
			exp:    [2]int{2014, 11},
		},
	} {
		date, err := time.Parse(tc.layout, tc.date)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, tc.exp, monthIdx(date))
	}
}

func TestWeekIdx(t *testing.T) {
	for _, tc := range []idxTestCase{
		{
			date:   "2014-11-12T11:45:26.371Z",
			layout: "2006-01-02T15:04:05.000Z",
			exp:    [2]int{2014, 46},
		},
	} {
		date, err := time.Parse(tc.layout, tc.date)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, tc.exp, weekIdx(date))
	}
}
