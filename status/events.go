package status

import "time"

type StatusCreated struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
	UserID     string           `bson:"user_id" json:"user_id"`
}

func NewStatusCreated(id string, userID string) *StatusCreated {
	return &StatusCreated{time.Now(), id, userID}
}

func (e *StatusCreated) OccurredOn() time.Time {
	return e.occurredOn
}

type StatusRemoved struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
}

func NewStatusRemoved(id string) *StatusRemoved {
	return &StatusRemoved{time.Now(), id}
}

func (e *StatusRemoved) OccurredOn() time.Time {
	return e.occurredOn
}

type AllDay struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
}

func NewAllDay(id string) *AllDay {
	return &AllDay{time.Now(), id}
}

func (e *AllDay) OccurredOn() time.Time {
	return e.occurredOn
}

type StatusTime struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
	StartsAt   time.Time        `bson:"starts_at" json:"starts_at"`
	EndsAt     time.Time        `bson:"ends_at" json:"ends_at"`
}

func NewStatusTime(id string, startsAt time.Time, endsAt time.Time) *StatusTime {
	return &StatusTime{time.Now(), id, startsAt, endsAt}
}

func (e *StatusTime) OccurredOn() time.Time {
	return e.occurredOn
}

type StatusReasonChanged struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
	Reason     Reason           `bson:"reason" json:"reason"`
	Message    string           `bson:"message" json:"message"`
}

func NewStatusReasonChanged(id string, reason Reason, message string) *StatusReasonChanged {
	return &StatusReasonChanged{time.Now(), id, reason, message}
}

func (e *StatusReasonChanged) OccurredOn() time.Time {
	return e.occurredOn
}

type StatusTimePeriodChanged struct {
	occurredOn time.Time        `bson:"occurred_on" json:"occurred_on"`
	ID         string           `bson:"id" json:"id"`
	TimePeriod TimePeriod   `bson:"time_period" json:"time_period"`
}

func NewStatusTimePeriodChanged(id string, timePeriod TimePeriod) *StatusTimePeriodChanged {
	return &StatusTimePeriodChanged{time.Now(), id, timePeriod}
}

func (e *StatusTimePeriodChanged) OccurredOn() time.Time {
	return e.occurredOn
}
