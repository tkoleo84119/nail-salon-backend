package adminSchedule

// -------------------------------------------------------------------------------------

type GetResponse struct {
	ID        string            `json:"id"`
	WorkDate  string            `json:"workDate"`
	Note      string            `json:"note"`
	TimeSlots []GetTimeSlotInfo `json:"timeSlots"`
}

type GetTimeSlotInfo struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}
