package errors

import "net/http"

var (
	ErrUserNotFound = New(http.StatusNotFound, "the user was not found")
)
