package apperr

import (
	"log"
	"net/http"
)

const (
	Unexpected AppErrType = iota
	NotFound
	BadRequest
)

var (
	TypeToStatusCode = map[AppErrType]int{
		Unexpected: http.StatusInternalServerError,
		NotFound:   http.StatusNotFound,
		BadRequest: http.StatusBadRequest,
	}

	TypeToMessage = map[AppErrType]string{
		Unexpected: "Something went wrong",
		NotFound:   "Not found",
		BadRequest: "Bad Request",
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
		log.Println(e.Err)
	}

	return e.Message
}
