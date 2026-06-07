package controller

import "errors"

var ErrInvalidRequestBody = errors.New("invalid request body")

type ControllerError struct {
	Kind    error
	Message string
}

func (e ControllerError) Error() string {
	return e.Message
}

func (e ControllerError) Unwrap() error {
	return e.Kind
}

func invalidRequestBody() error {
	return ControllerError{Kind: ErrInvalidRequestBody, Message: "Invalid request body"}
}
