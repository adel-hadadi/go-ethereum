package respond

import (
	"encoding/json"
	"ethereum/util/apperr"
	"log"
	"net/http"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(w http.ResponseWriter, r *http.Request, data any, code ...int) {
	statusCode := http.StatusOK

	if len(code) > 0 {
		statusCode = code[0]
	}

	res := Response{
		Data: data,
	}

	bytes, _ := json.Marshal(&res)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

const DefaultFaileMessage = "Something went wrong"

func Faile(w http.ResponseWriter, message string, err error, code ...int) {
	statusCode := http.StatusInternalServerError

	if len(code) > 0 {
		statusCode = code[0]
	}

	log.Println("Error => ", err)

	if message == "" {
		message = DefaultFaileMessage
	}

	res := Response{
		Error:   true,
		Message: message,
	}

	bytes, _ := json.Marshal(&res)

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func WithErr(w http.ResponseWriter, err error) {
	appErr := err.(apperr.AppErr)

	res := Response{
		Error:   true,
		Message: appErr.Error(),
	}

	bytes, _ := json.Marshal(&res)

	w.WriteHeader(appErr.StatusCode)
	w.Write(bytes)
}
