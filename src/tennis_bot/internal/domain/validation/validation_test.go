package validation

import (
	"errors"
	"testing"
	"time"
)

func TestValidateTimeBounds(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	nowNextYear := time.Date(2026, time.December, 31, 12, 0, 0, 0, time.UTC)
	test := []struct {
		name        string
		now         time.Time
		start, end  time.Time
		expectedErr error
	}{
		{
			name:        "valid time",
			now:         now,
			start:       now.Add(3 * time.Hour),
			end:         now.Add(4 * time.Hour),
			expectedErr: nil,
		},
		{
			name:        "minimal lead time reached",
			now:         now,
			start:       now,
			end:         now.Add(2 * time.Hour),
			expectedErr: ErrCantBook,
		},
		{
			name:        "end before start",
			now:         now,
			start:       now.Add(3 * time.Hour),
			end:         now.Add(2 * time.Hour),
			expectedErr: ErrInvalidTimeFormat,
		},
		{
			name:        "book after court closes",
			now:         now,
			start:       now.Add(12 * time.Hour),
			end:         now.Add(13 * time.Hour),
			expectedErr: ErrInvalidTimeRange,
		},
		{
			name:        "max hours to book reached",
			now:         now,
			start:       now.Add(2 * time.Hour),
			end:         now.Add(6 * time.Hour),
			expectedErr: ErrLimitReached,
		},
		{
			name:        "happy path, days limit not exceed",
			now:         now,
			start:       now.AddDate(0, 0, 3),
			end:         now.AddDate(0, 0, 3).Add(1 * time.Hour),
			expectedErr: nil,
		},
		{
			name:        "happy path, days limit  exceeded",
			now:         now,
			start:       now.AddDate(0, 0, 7),
			end:         now.AddDate(0, 0, 7).Add(1 * time.Hour),
			expectedErr: ErrDayDiffReached,
		},
		{
			name:        "corner case: book in next year",
			now:         nowNextYear,
			start:       nowNextYear.AddDate(0, 0, 5),
			end:         nowNextYear.AddDate(0, 0, 5).Add(1 * time.Hour),
			expectedErr: nil,
		},
		{
			name:        "corner case: failed book in next year",
			now:         nowNextYear,
			start:       nowNextYear.AddDate(0, 0, 7),
			end:         nowNextYear.AddDate(0, 0, 7).Add(1 * time.Hour),
			expectedErr: ErrDayDiffReached,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateTimeBounds(tc.now, tc.start, tc.end); !errors.Is(err, tc.expectedErr) {
				t.Log(tc.now, tc.start, tc.end)
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
