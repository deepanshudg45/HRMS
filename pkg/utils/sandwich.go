package utils

import "time"

func ApplySandwichPolicy(from, to time.Time, holidays []time.Time) int {
	if to.Before(from) {
		return 0
	}

	total := 0

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		total++
	}

	return total // includes weekends + holidays
}