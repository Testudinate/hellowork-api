package user

import (
	"github.com/hellofresh/goengine"
	"time"
)

type User struct {
	*goengine.AggregateRootBased
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"created_at" json:"created_at"`
	Username  string    `bson:"username" json:"username"`
}

func NewUser(id string, username string) *User {
	user := new(User)
	user.AggregateRootBased = goengine.NewAggregateRootBased(user)
	user.RecordThat(NewUserCreated(id, username))

	return user
}

func NewUserFromHistory(id string, streamName goengine.StreamName, repo goengine.AggregateRepository) (*User, error) {
	user := new(User)
	user.AggregateRootBased = goengine.NewEventSourceBasedWithID(user, id)
	err := repo.Reconstitute(id, user, streamName)

	return user, err
}

func (u *User) WhenUserCreated(event *UserCreated) {
	u.ID = event.UserID
	u.Username = event.Username
	u.CreatedAt = event.OccurredOn()
}
