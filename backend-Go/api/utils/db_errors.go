package utils

type DatabaseError struct {
	Message string
}

func (e DatabaseError) Error() string {
	return e.Message
}

var (
	ErrForeignKeyViolation      = &DatabaseError{Message: "foreign key violation"}
	ErrCheckConstraintViolation = &DatabaseError{Message: "check constraint violation"}
	ErrUniqueViolation          = &DatabaseError{Message: "unique violation"}
)
