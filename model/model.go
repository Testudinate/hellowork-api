package model

import (
	"fmt"
	"time"
)

var (
	OutOfOffice Reason = "out of office"
	Remote      Reason = "working remote"
	Sick        Reason = "sick"
	Vacation    Reason = "vacation"
	WorkTrip    Reason = "work trip"
)

type UserID string
type Reason string

type User struct {
	ID       UserID
	Username string
	Statuses []*Status
}

func NewUser(id UserID) *User {
	return &User{ID: id, Statuses: make([]*Status, 0)}
}

func (u *User) AddStatus(status *Status) {
	u.Statuses = append(u.Statuses, status)
}

func (u *User) GetStatus() *Status {
	return u.Statuses[len(u.Statuses) - 1]
}

func (u *User) IsAvailable(date time.Time) bool {
	var available bool
	for _, status := range u.Statuses {
		available = !status.isValid(date)
	}

	return available
}

type Status struct {
	Description string
	From        time.Time
	To          time.Time
	Reason      Reason
}

func NewStatus(description string, from time.Time, to time.Time, reason Reason) *Status {
	return &Status{description, from, to, reason}
}

func (s *Status) isValid(date time.Time) bool {
	before := s.From.Before(date)
	after := s.To.After(date)

	return !before && !after
}
