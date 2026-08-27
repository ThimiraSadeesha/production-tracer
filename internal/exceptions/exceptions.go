package exceptions

import "errors"

var (
	InvalidUsernameOrPassword     = errors.New("username or password invalid")
	UserAccountInactive           = errors.New("user account is not active")
	UserNotFoundInOurDB           = errors.New("user does not exists in our database")
	UserAlreadyExists             = errors.New("user already exists in our database")
	InternalServerError           = errors.New("internal server error")
	ResourceNotFound              = errors.New("resource not found")
	UnsupportedTypeCasting        = errors.New("unsupported models type")
	RequestBodyValidationFailed   = errors.New("request body validation failed")
	RequestQueryValidationFailed  = errors.New("request query validation failed")
	RequestParamsValidationFailed = errors.New("request params validation failed")
	InvalidatedToken              = errors.New("invalid user token")
	TokenExpired                  = errors.New("token has been expired")
	InvalidUserId                 = errors.New("you are trying to change that is not yours")
)
