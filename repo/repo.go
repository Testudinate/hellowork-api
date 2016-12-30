package repo

import (
	"time"
	"github.com/italolelis/hellowork-api/model"
)

type Repository interface {
	Find(id string) *model.User
	FindAll() []*model.User
	FindAllOut(date time.Time) []*model.User
	FindAllByID(usernames []string, date time.Time) []*model.User
	Add(user *model.User)
	Remove(username string)
}
