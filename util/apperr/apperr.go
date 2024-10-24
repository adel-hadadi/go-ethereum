package apperr

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	Unexpected AppErrType = iota
	NotFound
	BadRequest
	Confilict
)

var (
	TypeToStatusCode = map[AppErrType]int{
		Unexpected: http.StatusInternalServerError,
		NotFound:   http.StatusNotFound,
		BadRequest: http.StatusBadRequest,
		Confilict:  http.StatusConflict,
	}

	TypeToMessage = map[AppErrType]string{
		Unexpected: "Something went wrong",
		NotFound:   "Not found",
		BadRequest: "Bad Request",
		Confilict:  "Data Confilict",
	}
)

type AppErrType int

type AppErr struct {
	Type       AppErrType
	StatusCode int
	Err        error
	Message    string
}

func New(errType AppErrType) AppErr {
	return AppErr{
		Type:       errType,
		StatusCode: TypeToStatusCode[errType],
		Message:    TypeToMessage[errType],
	}
}

func (e AppErr) WithStatusCode(statusCdoe int) AppErr {
	e.StatusCode = statusCdoe

	return e
}

func (e AppErr) WithMessage(msg string) AppErr {
	e.Message = msg

	return e
}

func (e AppErr) WithErr(err error) AppErr {
	e.Err = err

	return e
}

func (e AppErr) Error() string {
	if e.Err != nil {
		logrus.Warning(e.Err)
	}

	return e.Message
}
