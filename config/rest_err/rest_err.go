package rest_err

import (
	"net/http"

	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type RestErr struct {
	Message string  `json:"message"`
	Code    int     `json:"code"`
	Err     string  `json:"err"`
	Causes  []Cause `json:"causes"`
}

type Cause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewBadRequestError(message string, causes ...Cause) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     "bad_request",
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     "internal_server_error",
		Causes:  nil,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusNotFound,
		Err:     "not_found",
		Causes:  nil,
	}
}

func ConvertError(internalError *internal_errors.InternalError) *RestErr {
	switch internalError.Err {
	case "bad_request":
		return NewBadRequestError(internalError.Message)
	case "not_found":
		return NewNotFoundError(internalError.Message)
	default:
		return NewInternalServerError(internalError.Message)
	}
}
