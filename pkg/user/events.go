package user

import "time"

type UserCreated struct {
	Occurred time.Time `bson:"occurred_on" json:"occurred_on"`
	UserID   string    `bson:"user_id" json:"user_id"`
	Username string    `bson:"username" json:"username"`
}

func NewUserCreated(userID string, username string) *UserCreated {
	return &UserCreated{time.Now(), userID, username}
}

func (e *UserCreated) OccurredOn() time.Time {
	return e.Occurred
}
