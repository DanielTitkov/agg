# agg
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Little Go package for making aggregation over slice on any objects by date with user-provided functions.

## Usage

**agg.ByDate** encapsulated logic for aggregation by date. In order for it to be able to do in user needs to provide functions to get date from aggregated items and also to merge any number of items to aggregated result. Function returns the number **N** of aggregated result items. After applying ByDate function slice will be sorted by date. First **N** positions will be occupied with aggregated items. As soon as number of aggregated results is always less or equal to lenght of initial slice some posisions may be occupied with unaggregated items. In order to get slice on only aggregated results you need to cut it up to **N**: `data = data[:N]`.

The date which ByDate sends back to merge function represent the start of the time period (day, week or month) to which items to be aggregated are assinged. 

See more in example.

## Example

```go
type item struct {
	Date   time.Time
	Value  float64
	Merged bool
}

func main() {
	data := []item{
		{
			Value: 1,
			Date:  time.Now().AddDate(0, 0, -2), // 2 days ago
		},
		{
			Value: 1,
			Date:  time.Now().AddDate(0, 0, -2), // also 2 days ago
		},
		{
			Value: 2,
			Date:  time.Now().AddDate(0, 0, -1), // yesterday
		},
		{
			Value: 2,
			Date:  time.Now().AddDate(0, 0, -1), // also yesterday
		},
		{
			Value: 3,
			Date:  time.Now(), // today
		},
		{
			Value: 3,
			Date:  time.Now(), // also today
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
```

This example with print something like following:
```
Result: [{Date:2022-02-15 00:00:00 +0000 UTC Value:2 Merged:true} {Date:2022-02-16 00:00:00 +0000 UTC Value:4 Merged:true} {Date:2022-02-17 00:00:00 +0000 UTC Value:6 Merged:true}]
```