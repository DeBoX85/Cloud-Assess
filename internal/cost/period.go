package cost

import "time"

// PreviousCompletedMonth returns the exact previous completed UTC calendar month.
func PreviousCompletedMonth(now time.Time) (time.Time, time.Time) {
	now = now.UTC()
	start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
	return start, end
}
