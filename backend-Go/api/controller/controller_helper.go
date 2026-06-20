package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
)

const maxJSONBodyBytes int64 = 1 << 20

func BindJSON(ec *appcontext.GinContext, destination any) bool {
	contentType := ec.GetHeader("Content-Type")
	if contentType == "" || !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		HandleError(ec, http.StatusUnsupportedMediaType, errors.New("unsupported content type"), "Content-Type must be application/json")
		return true
	}

	ec.Request.Body = http.MaxBytesReader(ec.Writer, ec.Request.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(ec.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		HandleError(ec, http.StatusBadRequest, err, "Invalid request body")
		return true
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("request body must contain a single JSON object")
		}
		HandleError(ec, http.StatusBadRequest, err, "Invalid request body")
		return true
	}

	return false
}

func HandleError(ec *appcontext.GinContext, status int, err error, message string) bool {
	if err == nil {
		return false
	}

	ec.Logger().Printf("request error status=%d error=%v message=%q", status, err, message)
	ec.JSON(status, dto.ErrorResponse{
		Code:      errorCode(status),
		Message:   message,
		RequestID: ec.TransactionID(),
	})
	return true
}

func errorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnsupportedMediaType:
		return "UNSUPPORTED_MEDIA_TYPE"
	default:
		return "INTERNAL_ERROR"
	}
}
