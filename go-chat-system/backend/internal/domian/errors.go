package domain

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrNotFound         = &Error{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrUnauthorized     = &Error{Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrInvalidInput     = &Error{Code: "INVALID_INPUT", Message: "Invalid input"}
	ErrAlreadyExists    = &Error{Code: "ALREADY_EXISTS", Message: "Resource already exists"}
	ErrPermissionDenied = &Error{Code: "PERMISSION_DENIED", Message: "Permission denied"}
	ErrInternalError    = &Error{Code: "INTERNAL_ERROR", Message: "Internal server error"}
)
