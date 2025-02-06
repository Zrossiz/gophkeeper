// Package apperrors defines common application errors that can be used across the application.
// These errors are intended to provide consistent error messages and handling for various scenarios,
// such as user authentication, database operations, and server-side issues.
package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a requested user is not found in the system.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when attempting to create a user that already exists.
	ErrUserAlreadyExists = errors.New("user already exist")

	// ErrHashPassword is returned when there is an error hashing a password.
	ErrHashPassword = errors.New("error hash password")

	// ErrDBQuery is returned when there is an error executing a database query.
	ErrDBQuery = errors.New("db error")

	// ErrJWTGeneration is returned when there is an error generating a JWT token.
	ErrJWTGeneration = errors.New("jwt generation error")

	// ErrInvalidPassword is returned when the provided login credentials are invalid.
	ErrInvalidPassword = errors.New("invalid login or password")

	// ErrInternalServer is a string error message for internal server errors.
	// This is not an error type but a message that can be used in responses.
	ErrInternalServer = "internal server error"

	// ErrInvalidRequestBody is a string error message for invalid request bodies.
	// This is not an error type but a message that can be used in responses.
	ErrInvalidRequestBody = "invalid request body"
)
