package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errorMessages []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "url":
			errorMessages = append(errorMessages, fmt.Sprintf("invalid url : %s", err.Field()))
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}
