package err

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error ...
type Error interface {
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
