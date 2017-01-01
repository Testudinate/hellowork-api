package user

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/goengine"
	"github.com/hellofresh/janus/request"
	"github.com/italolelis/hellowork-api/errors"
	"github.com/italolelis/hellowork-api/response"
	"github.com/italolelis/hellowork-api/router"
)

// UsersController is the api rest controller
type UsersController struct {
	writeRepo  goengine.AggregateRepository
	streamName goengine.StreamName
}

// NewUsersController creates a new instance of Controller
func NewUsersController(repo goengine.AggregateRepository) *UsersController {
	return &UsersController{repo, "users"}
}

func (u *UsersController) Get(w http.ResponseWriter, r *http.Request) {
	//data, err := u.repo.FindAll()
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//response.JSON(w, http.StatusOK, data)
}

func (u *UsersController) GetBy(w http.ResponseWriter, r *http.Request) {
	id := router.FromContext(r.Context()).ByName("id")
	user, err := NewUserFromHistory(id, u.streamName, u.writeRepo)
	if err != nil {
		log.Panic(err)
	}

	if user.ID == "" {
		panic(errors.ErrUserNotFound)
	}

	response.JSON(w, http.StatusOK, user)
}

func (u *UsersController) PutBy(w http.ResponseWriter, r *http.Request) {
	var err error

	id := router.FromContext(r.Context()).ByName("id")
	user, err := NewUserFromHistory(id, u.streamName, u.writeRepo)
	if err != nil {
		log.Panic(err)
	}

	if user.ID == "" {
		panic(errors.ErrUserNotFound)
	}

	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	err = request.BindJSON(r, user)
	if nil != err {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	err = u.writeRepo.Save(user, u.streamName)
	if nil != err {
		panic(errors.New(http.StatusBadRequest, err.Error()))
	}

	response.JSON(w, http.StatusOK, nil)
}

func (u *UsersController) Post(w http.ResponseWriter, r *http.Request) {
	var cmd CreateUser

	err := request.BindJSON(r, &cmd)
	if nil != err {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	user := NewUser(cmd.ID, cmd.Username)
	err = u.writeRepo.Save(user, u.streamName)
	if nil != err {
		panic(errors.New(http.StatusBadRequest, err.Error()))
	}

	response.JSON(w, http.StatusOK, nil)
}
