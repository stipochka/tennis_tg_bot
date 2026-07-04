package domain

import "time"

type ReservationKind string

var (
	ReservationKindBooking ReservationKind = "booking"
	ReservationKindBlocked ReservationKind = "blocked"
)

type ReservationStatus string

var (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusRejected  ReservationStatus = "rejected"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	ID         int64
	CourtID    int64
	UserID     int64
	Kind       ReservationKind
	Start      time.Time
	End        time.Time
	Status     ReservationStatus
	CreatedAt  time.Time
	ReviewedAt time.Time
}
