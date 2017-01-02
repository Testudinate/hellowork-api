package status

import (
	"net/http"

	"github.com/italolelis/hellowork-api/pkg/errors"
)

var (
	ErrEndTimeBeforeStartTime = errors.New(http.StatusBadRequest, "end time can't be before start time")
	ErrStartTimeInvalid       = errors.New(http.StatusBadRequest, "star time needs to be a valid")
	ErrReasonInvalid          = errors.New(http.StatusBadRequest, "the status reason can't be blank")

	ErrStatusNotFound = errors.New(http.StatusNotFound, "the status was not found")
)
