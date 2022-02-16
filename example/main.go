package main

import (
	"log"
	"time"

	"github.com/DanielTitkov/agg"
)

type item struct {
	Date   time.Time
	Value  float64
	Merged bool
}

func main() {
	data := []item{
		{
			Value: 1,
			Date:  time.Now().AddDate(0, 0, -2),
		},
		{
			Value: 1,
			Date:  time.Now().AddDate(0, 0, -2),
		},
		{
			Value: 2,
			Date:  time.Now().AddDate(0, 0, -1),
		},
		{
			Value: 2,
			Date:  time.Now().AddDate(0, 0, -1),
		},
		{
			Value: 3,
			Date:  time.Now(),
		},
		{
			Value: 3,
			Date:  time.Now(),
		},
	}

	i := agg.ByDate(
		data,
		agg.Day,
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
		time.UTC,
	)

	// Now we need to slice the result to keep only aggregated values.
	// This is important because ByDate is not able
	// to remove not-aggregated items from slice.
	data = data[:i]

	log.Printf("Result: %+v", data)
}
