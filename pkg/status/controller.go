package status

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/italolelis/hellowork-api/pkg/errors"
	"github.com/italolelis/hellowork-api/pkg/request"
	"github.com/italolelis/hellowork-api/pkg/response"
	"github.com/italolelis/hellowork-api/pkg/router"
)

// Controller is the api rest controller
type Controller struct {
	writeRepo WriteRepository
	readRepo  ReadStatusRepository
}

// NewUsersController creates a new instance of Controller
func NewController(writeRepo WriteRepository, readRepo ReadStatusRepository) *Controller {
	return &Controller{writeRepo, readRepo}
}

func (u *Controller) Get(w http.ResponseWriter, r *http.Request) {
	userID := router.FromContext(r.Context()).ByName("user_id")
	data, err := u.readRepo.FindAllByUserID(userID)
	if err != nil {
		panic(err.Error())
	}

	response.JSON(w, http.StatusOK, data)
}

func (u *Controller) GetBy(w http.ResponseWriter, r *http.Request) {
	id := router.FromContext(r.Context()).ByName("id")
	status, err := u.readRepo.Find(id)
	if err != nil {
		log.Panic(err)
	}

	if status.ID == "" {
		panic(ErrStatusNotFound)
	}

	response.JSON(w, http.StatusOK, status)
}

func (c *Controller) PutBy(w http.ResponseWriter, r *http.Request) {
	var err error
	var cmd UpdateStatus

	id := router.FromContext(r.Context()).ByName("id")
	status, err := NewStatusFromHistory(id, c.writeRepo)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	if status.ID == "" {
		panic(errors.ErrUserNotFound)
	}

	err = request.BindJSON(r, cmd)
	if nil != err {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	err = status.Because(c.parseReason(cmd.Reason), cmd.Message)
	if nil != err {
		panic(err)
	}

	if cmd.HasTimePeriod() {
		status.ChangeTimePeriod(c.parseTimePeriod(cmd.TimePeriod))
	} else {
		err = status.AtThisTime(cmd.StartsAt, cmd.EndsAt)
		if nil != err {
			panic(err)
		}
	}

	err = c.writeRepo.Add(status)
	if nil != err {
		panic(errors.New(http.StatusBadRequest, err.Error()))
	}

	response.JSON(w, http.StatusOK, nil)
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	var cmd CreateStatus

	err := request.BindJSON(r, &cmd)
	if nil != err {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	status := NewStatus(c.writeRepo.NextIdentity(), cmd.UserID)
	err = status.Because(c.parseReason(cmd.Reason), cmd.Message)
	if nil != err {
		panic(err)
	}

	if cmd.HasTimePeriod() {
		status.ChangeTimePeriod(c.parseTimePeriod(cmd.TimePeriod))
	} else {
		err = status.AtThisTime(cmd.StartsAt, cmd.EndsAt)
		if nil != err {
			panic(err)
		}
	}

	err = c.writeRepo.Add(status)
	if nil != err {
		panic(errors.New(http.StatusBadRequest, err.Error()))
	}

	response.JSON(w, http.StatusCreated, nil)
}

func (c *Controller) DeleteBy(w http.ResponseWriter, r *http.Request) {
	var cmd RemoveStatus
	cmd.ID = router.FromContext(r.Context()).ByName("id")

	status, err := NewStatusFromHistory(cmd.ID, c.writeRepo)
	if err != nil {
		panic(err)
	}
	status.Remove()
	c.writeRepo.Add(status)

	response.JSON(w, http.StatusOK, nil)
}

func (c *Controller) parseTimePeriod(timePeriod string) TimePeriod {
	switch strings.ToLower(timePeriod) {
	case "this_morning":
		return ThisMorning
	case "this_afternoon":
		return ThisAfternoon
	case "today":
		return Today
	case "tomorrow":
		return Tomorrow
	default:
		return ThisMorning
	}
}

func (c *Controller) parseReason(reason string) Reason {
	switch strings.ToLower(reason) {
	case "out of office":
		return OutOfOffice
	case "remote":
		return Remote
	case "sick":
		return Sick
	case "vacation":
		return Vacation
	case "work trip":
		return WorkTrip
	default:
		return OutOfOffice
	}
}
