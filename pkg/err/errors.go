package err

import (
	"google.golang.org/grpc/codes"
)

// // ApplicationError ...
// type ApplicationError struct{}

// // IsLDAPManagerError ...
// func (e *ApplicationError) IsLDAPManagerError() bool {
// 	return true
// }

// // Code ...
// func (e *ApplicationError) Code() codes.Code {
// 	return codes.Internal //e.code
// }

// Error ...
type Error interface {
	error
	IsLDAPManagerError() bool
	Code() codes.Code
}

// // ValidationError ...
// type ValidationError struct {
// 	ApplicationError
// 	Message string
// }

// // Error ...
// func (e *ValidationError) Error() string {
// 	return e.Message
// }

// // Code ...
// func (e *ValidationError) Code() codes.Code {
// 	return codes.InvalidArgument
// }
