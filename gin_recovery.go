package notify

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GinRecovery() gin.HandlerFunc {
	return recoverError()
}

func recoverError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e, isError := err.(error)
				if isError {
					ErrorRequest(e, c.Request)

					panic(e)
				}

				s, isString := err.(string)
				if isString {
					e = errors.New(s)
					ErrorRequest(e, c.Request)

					panic(e)
				}
			}
		}()

		c.Next()
	}
}
