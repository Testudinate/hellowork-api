package repo

import (
	"time"

	"github.com/italolelis/hellowork/model"
)

type InMemory struct {
	users map[model.UserID]*model.User
}

func NewInMemory() *InMemory {
	return &InMemory{make(map[model.UserID]*model.User)}
}

func (r *InMemory) Find(id string) *model.User {
	user, exists := r.users[model.UserID(id)]

	if !exists {
		return nil
	}

	return user
}

func (r *InMemory) FindAll() []*model.User {
	var users []*model.User

	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

func (r *InMemory) FindAllOut(date time.Time) []*model.User {
	var users []*model.User

	for _, user := range r.users {
		if user.IsAvailable(date) {
			users = append(users, user)
		}
	}

	return users
}

func (r *InMemory) FindAllByID(ids []string, date time.Time) []*model.User {
	var users []*model.User

	for _, id := range ids {
		user, exists := r.users[model.UserID(id)]
		if exists {
			users = append(users, user)
		}
	}

	return users
}

func (r *InMemory) Add(user *model.User) {
	r.users[user.ID] = user
}

func (r *InMemory) Remove(id string) {
	delete(r.users, model.UserID(id))
}
