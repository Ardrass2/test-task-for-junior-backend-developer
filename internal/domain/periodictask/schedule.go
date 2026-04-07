package periodictask

import "time"

func (pt *PeriodicTask) MatchesDate(date time.Time) bool {
	d := truncateToDate(date)
	start := truncateToDate(pt.StartDate)

	if d.Before(start) {
		return false
	}
	if pt.EndDate != nil {
		end := truncateToDate(*pt.EndDate)
		if d.After(end) {
			return false
		}
	}

	switch pt.RecurrenceType {
	case RecurrenceDaily:
		return pt.matchesDaily(d, start)
	case RecurrenceMonthly:
		return pt.matchesMonthly(d)
	case RecurrenceSpecificDates:
		return pt.matchesSpecificDates(d)
	case RecurrenceEvenOdd:
		return pt.matchesEvenOdd(d)
	default:
		return false
	}
}

func (pt *PeriodicTask) matchesDaily(date, start time.Time) bool {
	interval := 1
	if pt.RecurrenceParams.Interval != nil && *pt.RecurrenceParams.Interval > 0 {
		interval = *pt.RecurrenceParams.Interval
	}
	days := int(date.Sub(start).Hours() / 24)
	return days%interval == 0
}

func (pt *PeriodicTask) matchesMonthly(date time.Time) bool {
	if pt.RecurrenceParams.Day == nil {
		return false
	}
	return date.Day() == *pt.RecurrenceParams.Day
}

func (pt *PeriodicTask) matchesSpecificDates(date time.Time) bool {
	formatted := date.Format("2006-01-02")
	for _, d := range pt.RecurrenceParams.Dates {
		if d == formatted {
			return true
		}
	}
	return false
}

func (pt *PeriodicTask) matchesEvenOdd(date time.Time) bool {
	if pt.RecurrenceParams.Parity == nil {
		return false
	}
	dayOfMonth := date.Day()
	if *pt.RecurrenceParams.Parity == ParityEven {
		return dayOfMonth%2 == 0
	}
	return dayOfMonth%2 != 0
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
