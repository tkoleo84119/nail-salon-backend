package schedule

// DeleteTimeSlotResponse represents the response after deleting a time slot
type DeleteTimeSlotResponse struct {
	Deleted []string `json:"deleted"`
}