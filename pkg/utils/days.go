package utils

import "time"

func isHoliday(date time.Time, holidays []time.Time) bool {
	for _, h := range holidays {
		if date.Year() == h.Year() && date.YearDay() == h.YearDay() {
			return true
		}
	}
	return false
}

func WorkingDays(from, to time.Time, holidays []time.Time) int {
	if to.Before(from) {
		return 0
	}

	count := 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()

		if weekday == time.Saturday || weekday == time.Sunday {
			continue
		}
		if isHoliday(d, holidays) {
			continue
		}
		count++
	}
	return count
}