// Package controller provides shared HTTP helpers for the API.
package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
)

const maxJSONBodyBytes int64 = 1 << 20

// BindJSON parses an application/json request body into the provided destination.
func BindJSON(ec *appcontext.GinContext, destination any) bool {
	contentType := ec.GetHeader("Content-Type")
	if contentType == "" || !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		HandleError(ec, http.StatusUnsupportedMediaType, ErrInvalidContentType, "Content-Type must be application/json")
		return true
	}

	ec.Request.Body = http.MaxBytesReader(ec.Writer, ec.Request.Body, maxJSONBodyBytes)
	decoder := json.NewDecoder(ec.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		HandleError(ec, http.StatusBadRequest, errors.Join(ErrInvalidRequestBody, err), "Invalid request body")
		return true
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			err = ErrInvalidRequestBody
		}
		HandleError(ec, http.StatusBadRequest, err, "Invalid request body")
		return true
	}

	return false
}

// HandleError logs the request failure and writes a safe API error response.
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

// errorCode maps HTTP statuses to stable API error codes.
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

// ParseID validates a path or query identifier as a positive integer.
func ParseID(value string) (int64, error) {
	return parseID(value)
}

// ParsePair validates the asset and vulnerability identifiers from a request context.
func ParsePair(ec *appcontext.GinContext) (int64, int64, bool) {
	return parsePair(ec)
}

// parseID parses a positive integer identifier.
func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidIdentifier
	}
	return id, nil
}

// parsePair parses the asset and vulnerability identifiers from the request.
func parsePair(ec *appcontext.GinContext) (int64, int64, bool) {
	assetID, err := parseID(ec.Param("id"))
	if err != nil {
		return 0, 0, false
	}
	vulnerabilityID, err := parseID(ec.Param("vulnerabilityId"))
	if err != nil {
		return 0, 0, false
	}
	return assetID, vulnerabilityID, true
}
