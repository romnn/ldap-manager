package err

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ApplicationError is an application error that is visible to the end user.
//
// GRPC and HTTP check if an error is an ApplicationError and transparently
// pass them to the user.
type ApplicationError interface {
	error
	IsLDAPManagerError() bool
	StatusError() error
}

// ValidationError ...
type ValidationError struct {
	error
	Message string
}

// Error ...
func (e *ValidationError) Error() string {
	return e.Message
}

// StatusError ...
func (e *ValidationError) StatusError() error {
	return status.Errorf(codes.InvalidArgument, e.Error())
}
