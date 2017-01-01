package status

import (
	"time"
	"github.com/hellofresh/goengine"
)

var (
	OutOfOffice Reason = "out_of_office"
	Remote      Reason = "working_remote"
	Sick        Reason = "sick"
	Vacation    Reason = "vacation"
	WorkTrip    Reason = "work trip"

	ThisMorning   TimePeriod = "this_morning"
	ThisAfternoon TimePeriod = "this_afternoon"
	Today         TimePeriod = "today"
	Tomorrow      TimePeriod = "tomorrow"
)

type TimePeriod string
type Reason string

type Status struct {
	*goengine.AggregateRootBased
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	IsActive   bool `bson:"is_active"`
	UserID     string `bson:"id"`
	Reason     Reason `bson:"reason"`
	IsAllDay   bool `bson:"is_all_day"`
	StartsAt   time.Time `bson:"starts_at"`
	EndsAt     time.Time `bson:"ends_at"`
	Message    string `bson:"message"`
	TimePeriod TimePeriod `bson:"time_period"`
}

func NewStatus(id string, userID string) *Status {
	status := new(Status)
	status.AggregateRootBased = goengine.NewAggregateRootBased(status)
	status.RecordThat(NewStatusCreated(id, userID))

	return status
}

func NewStatusFromHistory(id string, repo WriteRepository) (*Status, error) {
	status := new(Status)
	status.AggregateRootBased = goengine.NewEventSourceBasedWithID(status, id)
	err := repo.Reconstitute(id, status)

	return status, err
}

func (s *Status) Remove() {
	s.RecordThat(NewStatusRemoved(s.ID))
}

func (s *Status) AllDay() {
	s.RecordThat(NewAllDay(s.ID))
}

func (s *Status) AtThisTime(startsAt time.Time, endsAt time.Time) error {
	if startsAt.Before(time.Now()) {
		return ErrStartTimeInvalid
	}

	if endsAt.Before(startsAt) {
		return ErrEndTimeBeforeStartTime
	}

	s.RecordThat(NewStatusTime(s.ID, startsAt, endsAt))
	return nil
}

func (s *Status) ChangeTimePeriod(timePeriod TimePeriod) {
	var startsAt time.Time
	var endsAt time.Time
	var now time.Time = time.Now().UTC()

	switch timePeriod {
	case ThisMorning:
		startsAt = time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
		endsAt = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		s.AtThisTime(startsAt, endsAt)
	case ThisAfternoon:
		startsAt = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		endsAt = time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
		s.AtThisTime(startsAt, endsAt)
	case Today:
		s.AllDay()
	case Tomorrow:
		s.AllDay()
	}

	s.RecordThat(NewStatusTimePeriodChanged(s.ID, timePeriod))
}

func (s *Status) Because(reason Reason, message string) error {
	if reason == "" {
		return ErrReasonInvalid
	}

	if reason != s.Reason {
		return nil
	}

	s.RecordThat(NewStatusReasonChanged(s.ID, reason, message))
	return nil
}

func (s *Status) WhenStatusCreated(event *StatusCreated) {
	s.ID = event.ID
	s.CreatedAt = time.Now()
	s.UserID = event.UserID
	s.IsActive = true
	s.updated()
}

func (s *Status) WhenStatusRemoved(event *StatusRemoved) {
	s.IsActive = false
}

func (s *Status) WhenAllDay(event *AllDay) {
	s.IsAllDay = true
}

func (s *Status) WhenReasonChanged(event *StatusReasonChanged) {
	s.Reason = event.Reason
	s.Message = event.Message
}

func (s *Status) WhenStatusTime(event *StatusTime) {
	s.StartsAt = event.StartsAt.UTC()
	s.EndsAt = event.EndsAt.UTC()
}

func (s *Status) WhenStatusTimePeriodChanged(event *StatusTimePeriodChanged) {
	s.TimePeriod = event.TimePeriod
}

func (s *Status) updated() {
	s.UpdatedAt = time.Now()
}
