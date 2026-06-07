package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/service"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrRemoteService  = errors.New("remote service error")
	ErrDuplicateValue = errors.New("duplicate value")
)

type APIError struct {
	Status  int
	Message string
}

func (e APIError) Error() string {
	return e.Message
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}

func HandleError(w http.ResponseWriter, err error) {
	var apiErr APIError
	if errors.As(err, &apiErr) {
		WriteText(w, apiErr.Status, apiErr.Message)
		return
	}

	switch {
	case errors.Is(err, ErrUnauthorized), errors.Is(err, service.ErrInvalidCredentials):
		WriteText(w, http.StatusUnauthorized, "Invalid credentials.")
	case errors.Is(err, ErrBadRequest), errors.Is(err, service.ErrInvalidRequestData), errors.Is(err, service.ErrRemoteRejected):
		WriteText(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrDuplicateValue), errors.Is(err, service.ErrConflict):
		WriteText(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrNotFound), errors.Is(err, service.ErrNotFound):
		WriteText(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrRemoteService), errors.Is(err, service.ErrRemoteService), errors.Is(err, service.ErrInvalidRemoteResult):
		WriteText(w, http.StatusBadGateway, "Risk service is unavailable.")
	default:
		WriteText(w, http.StatusInternalServerError, "An unexpected error occurred.")
	}
}

func HandleGinError(c *gin.Context, err error) {
	var apiErr APIError
	if errors.As(err, &apiErr) {
		c.String(apiErr.Status, apiErr.Message)
		return
	}

	switch {
	case errors.Is(err, ErrUnauthorized), errors.Is(err, service.ErrInvalidCredentials):
		c.String(http.StatusUnauthorized, "Invalid credentials.")
	case errors.Is(err, ErrBadRequest), errors.Is(err, service.ErrInvalidRequestData), errors.Is(err, service.ErrRemoteRejected):
		c.String(http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrDuplicateValue), errors.Is(err, service.ErrConflict):
		c.String(http.StatusConflict, err.Error())
	case errors.Is(err, ErrNotFound), errors.Is(err, service.ErrNotFound):
		c.String(http.StatusNotFound, err.Error())
	case errors.Is(err, ErrRemoteService), errors.Is(err, service.ErrRemoteService), errors.Is(err, service.ErrInvalidRemoteResult):
		c.String(http.StatusBadGateway, "Risk service is unavailable.")
	default:
		c.String(http.StatusInternalServerError, "An unexpected error occurred.")
	}
}
