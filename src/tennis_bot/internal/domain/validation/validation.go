package validation

import (
	"time"
)

const (
	minLeadTime     = 1
	maxHourToBook   = 3
	courtOpenedHour = 7
	courtClosedHour = 23
)

func ValidateTimeBounds(now, start, end time.Time) error {
	type condErr struct {
		condition bool
		err       error
	}

	checks := []condErr{
		{start.Sub(start).Hours() < minLeadTime, ErrCantBook},
		{start.Year() == end.Year(), ErrInvalidDate},
		{start.Month() == end.Month(), ErrInvalidDate},
		{start.Day() == end.Day(), ErrInvalidDate},
		{end.After(start), ErrInvalidTimeFormat},
		{start.Hour() >= courtOpenedHour, ErrInvalidTimeRange},
		{end.Hour() <= courtClosedHour, ErrInvalidTimeRange},
		{start.Minute() == 0 && end.Minute() == 0, ErrInvalidTimeFormat},
		{start.Second() == 0 && end.Second() == 0, ErrInvalidTimeFormat},
		{end.Sub(start).Hours() <= maxHourToBook, ErrLimitReached},
	}

	for _, check := range checks {
		if err := errIfFalse(check.condition, check.err); err != nil {
			return err
		}
	}

	return nil
}

func errIfFalse(cond bool, err error) error {
	if !cond {
		return err
	}

	return nil
}
