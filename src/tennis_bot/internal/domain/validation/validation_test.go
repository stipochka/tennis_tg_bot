package validation

import (
	"errors"
	"testing"
	"time"
)

func TestValidateTimeBounds(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

	test := []struct {
		name        string
		start, end  time.Time
		expectedErr error
	}{
		{
			name:        "valid time",
			start:       now.Add(3 * time.Hour),
			end:         now.Add(4 * time.Hour),
			expectedErr: nil,
		},
		{
			name:        "minimal lead time reached",
			start:       now,
			end:         now.Add(2 * time.Hour),
			expectedErr: ErrCantBook,
		},
		{
			name:        "end before start",
			start:       now.Add(3 * time.Hour),
			end:         now.Add(2 * time.Hour),
			expectedErr: ErrInvalidTimeFormat,
		},
		{
			name:        "start after court closes",
			start:       now.Add(12 * time.Hour),
			end:         now.Add(13 * time.Hour),
			expectedErr: ErrInvalidTimeRange,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateTimeBounds(now, tc.start, tc.end); !errors.Is(err, tc.expectedErr) {
				t.Log(now, tc.start, tc.end)
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
