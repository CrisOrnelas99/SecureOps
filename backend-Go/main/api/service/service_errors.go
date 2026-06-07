package service

type ServiceError struct {
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

var (
	ErrInvalidRequestData  = &ServiceError{Message: "invalid request data"}
	ErrConflict            = &ServiceError{Message: "conflict"}
	ErrNotFound            = &ServiceError{Message: "not found"}
	ErrInvalidCredentials  = &ServiceError{Message: "invalid credentials"}
	ErrRemoteService       = &ServiceError{Message: "remote service error"}
	ErrRemoteRejected      = &ServiceError{Message: "remote service rejected request"}
	ErrInvalidRemoteResult = &ServiceError{Message: "invalid remote service response"}
)
