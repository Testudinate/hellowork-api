package status

import "time"

type Command struct {
	UserID     string    `json:"user_id"`
	Reason     string    `json:"reason"`
	Message    string    `json:"message"`
	IsAllDay   string    `json:"is_all_day"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
	TimePeriod string    `json:"time_period"`
}

func (c Command) HasTimePeriod() bool {
	return len(c.TimePeriod) > 0
}

type CreateStatus struct {
	Command
}

type UpdateStatus struct {
	Command
	ID string `json:"id"`
}

type RemoveStatus struct {
	ID string `json:"id"`
}
