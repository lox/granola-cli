package output

import "time"

type Meeting struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
}

type MeetingDetail struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	Summary   string    `json:"summary,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	Attendees []string  `json:"attendees,omitempty"`
}
