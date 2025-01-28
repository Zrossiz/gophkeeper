package apperrors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exist")
	ErrHashPassword      = errors.New("error hash password")
	ErrDBQuery           = errors.New("db error")
	ErrJWTGeneration     = errors.New("jwt generation error")
	ErrInvalidPassword   = errors.New("invalid login or password")

	ErrInternalServer     = "internal server error"
	ErrInvalidRequestBody = "invalid request body"
)
