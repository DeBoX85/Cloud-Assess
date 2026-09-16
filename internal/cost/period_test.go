package cost

import (
	"testing"
	"time"
)

func TestPreviousCompletedMonth(t *testing.T) {
	tests := []struct {
		name      string
		now       time.Time
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "regular month",
			now:       time.Date(2026, time.September, 16, 12, 0, 0, 0, time.FixedZone("CEST", 2*60*60)),
			wantStart: time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
		{
			name:      "year boundary",
			now:       time.Date(2026, time.January, 10, 5, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := PreviousCompletedMonth(tt.now)
			if !start.Equal(tt.wantStart) || !end.Equal(tt.wantEnd) {
				t.Fatalf("got %s - %s, want %s - %s", start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}
