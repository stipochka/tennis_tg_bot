package reservation

import "errors"

var (
	ErrReservationOverlap = errors.New("reservation time is crossing others")
)
