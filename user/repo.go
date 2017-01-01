package user

import (
	"time"
)

type Repository interface {
	Find(id string) *User
	FindAll() []*User
	FindAllOut(date time.Time) []*User
	FindAllByID(usernames []string, date time.Time) []*User
	Add(user *User)
	Remove(username string)
}
