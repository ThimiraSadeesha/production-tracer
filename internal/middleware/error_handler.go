package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/thimira/production-tracer/internal/exceptions"

	"github.com/gin-gonic/gin"
	"github.com/rashintha/logger"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorln("🔥 PANIC recovered.")
				logger.Errorln(string(debug.Stack()))

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": exceptions.InternalServerError.Error(),
				})
			}
		}()
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleAppError(c, err)
		}
	}
}

func handleAppError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var statusCode int
	var message string

	switch {
	case errors.Is(err, exceptions.InvalidUsernameOrPassword):
		statusCode, message = http.StatusUnauthorized, err.Error()

	case errors.Is(err, exceptions.UserAccountInactive):
		statusCode, message = http.StatusForbidden, err.Error()

	case errors.Is(err, exceptions.UserNotFoundInOurDB):
		statusCode, message = http.StatusNotFound, err.Error()

	case errors.Is(err, exceptions.UserAlreadyExists):
		statusCode, message = http.StatusConflict, err.Error()

	case errors.Is(err, exceptions.TokenExpired),
		errors.Is(err, exceptions.InvalidatedToken):
		statusCode, message = http.StatusUnauthorized, err.Error()

	case errors.Is(err, exceptions.ResourceNotFound):
		statusCode, message = http.StatusNotFound, err.Error()

	case errors.Is(err, exceptions.RequestBodyValidationFailed),
		errors.Is(err, exceptions.RequestQueryValidationFailed),
		errors.Is(err, exceptions.RequestParamsValidationFailed):
		statusCode, message = http.StatusBadRequest, err.Error()

	default:
		statusCode, message = http.StatusInternalServerError, exceptions.InternalServerError.Error()
		logger.Errorln("💥 Unhandled error:" + err.Error())
	}

	c.AbortWithStatusJSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}
