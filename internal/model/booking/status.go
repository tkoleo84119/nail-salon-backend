package booking

const (
	BookingStatusScheduled = "SCHEDULED"
	BookingStatusCancelled = "CANCELLED"
	BookingStatusCompleted = "COMPLETED"
	BookingStatusNoShow    = "NO_SHOW"
)

func IsValidBookingStatus(status string) bool {
	switch status {
	case BookingStatusScheduled, BookingStatusCancelled, BookingStatusCompleted, BookingStatusNoShow:
		return true
	default:
		return false
	}
}
