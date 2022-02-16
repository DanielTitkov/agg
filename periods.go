package aggr

import (
	"time"
)

func weekIdx(t time.Time) [2]int {
	y, w := t.ISOWeek()
	return [2]int{y, w}
}

func monthIdx(t time.Time) [2]int {
	return [2]int{t.Year(), int(t.Month())}
}

func dayIdx(t time.Time) [2]int {
	return [2]int{t.Year(), t.YearDay()}
}

func idxToDayStart(i [2]int, loc *time.Location) time.Time {
	t := time.Date(i[0], 1, 1, 0, 0, 0, 0, loc)
	return t.AddDate(0, 0, int(i[1]-1))
}

func idxToWeekStart(i [2]int, loc *time.Location) time.Time {
	year, week := i[0], i[1]
	// start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, loc)

	// roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// difference in weeks:
	_, w := t.ISOWeek()
	return t.AddDate(0, 0, (week-w)*7)
}

func idxToMonthStart(i [2]int, loc *time.Location) time.Time {
	return time.Date(i[0], time.Month(i[1]), 1, 0, 0, 0, 0, loc)
}
