package common

import "time"

type ActivityLogType string

const (
	ActivityCustomerRegister        ActivityLogType = "CUSTOMER_REGISTER"
	ActivityCustomerBooking         ActivityLogType = "CUSTOMER_BOOKING"
	ActivityCustomerBrowse          ActivityLogType = "CUSTOMER_BROWSE"
	ActivityCustomerBookingUpdate   ActivityLogType = "CUSTOMER_BOOKING_UPDATE"
	ActivityCustomerBookingCancel   ActivityLogType = "CUSTOMER_BOOKING_CANCEL"
	ActivityAdminBookingCreate      ActivityLogType = "ADMIN_BOOKING_CREATE"
	ActivityAdminBookingUpdate      ActivityLogType = "ADMIN_BOOKING_UPDATE"
	ActivityAdminBookingCancel      ActivityLogType = "ADMIN_BOOKING_CANCEL"
	ActivityAdminBookingCompleted   ActivityLogType = "ADMIN_BOOKING_COMPLETED"
)

type ActivityLogEntry struct {
	ID        string          `json:"id"`
	Type      ActivityLogType `json:"type"`
	Message   string          `json:"message"`
	Timestamp time.Time       `json:"timestamp"`
}

type ActivityLogResponse struct {
	Activities []ActivityLogEntry `json:"activities"`
	Total      int                `json:"total"`
}