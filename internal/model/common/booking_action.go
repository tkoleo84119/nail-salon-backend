package common

type BookingAction string

const (
	BookingActionCreated   BookingAction = "created"
	BookingActionUpdated   BookingAction = "updated"
	BookingActionCancelled BookingAction = "cancelled"
)
